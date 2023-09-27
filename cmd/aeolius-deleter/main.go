package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/pojntfx/aeolius/pkg/bluesky"
	"github.com/pojntfx/aeolius/pkg/persisters"
)

func main() {
	rateLimitPointsDID := flag.Int("rate-limit-points-did", 200, "Maximum amount of rate limit points to spend per DID (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 1666 per hour as of September 2023)")
	rateLimitPointsGlobal := flag.Int("rate-limit-points-global", 2500, "Maximum amount of rate limit points to spend per rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 3000 per hour as of September 2023)")
	rateLimitResetInterval := flag.Duration("rate-limit-reset-interval", time.Minute*5, "Duration of a rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; 5 minutes as of September 2023)")
	postgresUrl := flag.String("postgres-url", "postgresql://postgres@localhost:5432/aeolius?sslmode=disable", "PostgreSQL URL")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	persister := persisters.NewWorkerPersister(*postgresUrl)

	if err := persister.Open(); err != nil {
		panic(err)
	}
	defer persister.Close()

	log.Println("Connected to PostgreSQL")

	throttled := 0
	limiter := bluesky.NewLimiter(
		ctx,

		*rateLimitPointsGlobal,
		*rateLimitResetInterval,

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

		recordsToDelete, cursor, err := bluesky.GetPostsToDelete(
			client,

			int(configuration.PostTtl),
			configuration.Cursor,
			100,
			*rateLimitPointsDID,

			limiter,
		)
		if err != nil {
			log.Println("Could not get posts to delete for DID", auth.Did, ", skipping:", err)

			continue
		}

		log.Println("Deleting", recordsToDelete)

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

	log.Println("Spent", limiter.GetSpendPoints(), "points in", time.Since(before), "while being throttled", throttled, "times")
}
