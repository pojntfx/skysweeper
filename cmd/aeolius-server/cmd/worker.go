package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/pojntfx/aeolius/pkg/bluesky"
	"github.com/pojntfx/aeolius/pkg/persisters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	apiKeyFlag = "api-key"

	rateLimitPointsDIDFlag     = "rate-limit-points-did"
	rateLimitPointsGlobalFlag  = "rate-limit-points-global"
	rateLimitResetIntervalFlag = "rate-limit-reset-interval"
	listRecordsLimitFlag       = "list-records-limit"
	applyWritesLimitFlag       = "apply-writes-limit"
	dryRunFlag                 = "dry-run"

	verboseFlag = "verbose"
)

var (
	errMissingAPIKey = errors.New("missing API key")
)

type Statistics struct {
	SpentPoints  int `json:"spentPoints"`
	SpentTime    int `json:"spentTime"`
	Throttled    int `json:"throttled"`
	PostsDeleted int `json:"postsDeleted"`
}

var workerCmd = &cobra.Command{
	Use:     "worker",
	Aliases: []string{"w"},
	Short:   "Start an Aeolius worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if strings.TrimSpace(viper.GetString(apiKeyFlag)) == "" {
			return errMissingAPIKey
		}

		persister := persisters.NewWorkerPersister(viper.GetString(postgresURLFlag))

		if err := persister.Open(); err != nil {
			return err
		}
		defer persister.Close()

		log.Println("Connected to PostgreSQL")

		lis, err := net.Listen("tcp", viper.GetString(laddrFlag))
		if err != nil {
			return err
		}
		defer lis.Close()

		log.Println("Listening on", lis.Addr())

		mux := http.NewServeMux()

		mux.HandleFunc("/posts", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestAPIKey := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if strings.TrimSpace(requestAPIKey) == "" {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			if requestAPIKey != viper.GetString(apiKeyFlag) {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					log.Printf("Client disconnected with error: %v", err)
				}
			}()

			switch r.Method {
			case http.MethodDelete:
				throttled := 0
				limiter := bluesky.NewLimiter(
					ctx,

					viper.GetInt(rateLimitPointsGlobalFlag),
					viper.GetDuration(rateLimitResetIntervalFlag),

					func() error {
						log.Println("Pausing until rate limit reset interval")

						throttled++

						return nil
					},
				)

				go limiter.Open()

				before := time.Now()

				configurations, err := persister.GetEnabledConfigurations(ctx)
				if err != nil {
					panic(err)
				}

				postsDeleted := 0
				for _, configuration := range configurations {
					auth := &xrpc.AuthInfo{}

					client := &xrpc.Client{
						Client: http.DefaultClient,
						Host:   configuration.Service,
						Auth:   auth,
					}

					auth.AccessJwt = configuration.RefreshJwt
					auth.Did = configuration.Did

					session, err := atproto.ServerRefreshSession(ctx, client)
					if err != nil {
						log.Println("Could not refresh session for DID", auth.Did, ", skipping:", err)

						continue
					}

					auth.AccessJwt = session.AccessJwt
					auth.RefreshJwt = session.RefreshJwt
					auth.Handle = session.Handle
					auth.Did = session.Did

					postsToDelete, cursor, err := bluesky.GetPostsToDelete(
						client,

						int(configuration.PostTtl),
						configuration.Cursor,
						viper.GetInt(listRecordsLimitFlag), // Limit as per https://atproto.com/blog/rate-limits-pds-v3
						viper.GetInt(rateLimitPointsDIDFlag),

						limiter,
					)
					if err != nil {
						log.Println("Could not get posts to delete for DID", auth.Did, ", skipping:", err)

						continue
					}

					postsDeleted += len(postsToDelete)

					if err := bluesky.DeletePosts(
						ctx,

						client,

						postsToDelete,
						viper.GetInt(applyWritesLimitFlag),

						viper.GetBool(dryRunFlag),

						limiter,
					); err != nil {
						log.Println("Could not delete posts for DID", auth.Did, ", skipping:", err)

						continue
					}

					if err := persister.UpdateRefreshTokenAndCursor(
						ctx,
						auth.Did,
						cursor,
						auth.RefreshJwt,
					); err != nil {
						log.Println("Could not update refresh token and cursor for DID", auth.Did, ", skipping:", err)

						continue
					}
				}

				res := &Statistics{
					SpentPoints:  limiter.GetSpendPoints(),
					SpentTime:    int(time.Since(before)),
					Throttled:    throttled,
					PostsDeleted: postsDeleted,
				}

				if viper.GetBool(verboseFlag) {
					log.Println(
						"Spent", res.SpentPoints,
						"points in", res.SpentTime,
						"while being throttled", res.Throttled,
						"times to delete", postsDeleted,
						"posts (dry run mode", func() string {
							if viper.GetBool(dryRunFlag) {
								return "enabled)"
							}

							return "disabled)"
						}())
				}

				w.Header().Set("Content-Type", "application/json")

				if err := json.NewEncoder(w).Encode(res); err != nil {
					panic(fmt.Errorf("%w: %v", errCouldNotEncode, err))
				}

			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))

		if err := http.Serve(lis, mux); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	workerCmd.PersistentFlags().String(laddrFlag, "localhost:1338", "Listen address")
	workerCmd.PersistentFlags().String(apiKeyFlag, "", "API key to check incoming requests for")

	workerCmd.PersistentFlags().Int(rateLimitPointsDIDFlag, 200, "Maximum amount of rate limit points to spend per DID (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 1666 per hour as of September 2023)")
	workerCmd.PersistentFlags().Int(rateLimitPointsGlobalFlag, 2500, "Maximum amount of rate limit points to spend per rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 3000 per hour as of September 2023)")
	workerCmd.PersistentFlags().Duration(rateLimitResetIntervalFlag, time.Minute*5, "Duration of a rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; 5 minutes as of September 2023)")
	workerCmd.PersistentFlags().Int(listRecordsLimitFlag, 100, "Limit of records to return per API call (see https://atproto.com/blog/rate-limits-pds-v3; 100 as of September 2023)")
	workerCmd.PersistentFlags().Int(applyWritesLimitFlag, 10, "Limit of records to apply writes for per API call (see https://atproto.com/blog/rate-limits-pds-v3; 10 as of September 2023)")
	workerCmd.PersistentFlags().Bool(dryRunFlag, true, "Whether to do a dry run (only fetch for posts to be deleted without actually deleting them)")

	workerCmd.PersistentFlags().Bool(verboseFlag, false, "Whether to enable verbose logging")

	viper.AutomaticEnv()

	rootCmd.AddCommand(workerCmd)
}
