package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/pojntfx/aeolius/pkg/persisters"
)

type Configuration struct {
	Enabled bool  `json:"enabled"`
	PostTTL int32 `json:"postTTL"`
}

var (
	errCouldNotGetConfiguration    = errors.New("could not get configuraion")
	errCouldNotUpsertConfiguration = errors.New("could not upsert configuration")
	errCouldNotDeleteConfiguration = errors.New("could not delete configuration")

	errCouldNotEncode = errors.New("could not encode")
	errCouldNotDecode = errors.New("could not decode")

	errMissingService = errors.New("missing service")

	errCouldNotGetSession     = errors.New("could not get session")
	errCouldNotRefreshSession = errors.New("could not refresh session")
)

func main() {
	postgresUrl := flag.String("postgres-url", "postgresql://postgres@localhost:5432/aeolius?sslmode=disable", "PostgreSQL URL")
	laddr := flag.String("laddr", "localhost:1337", "Listen address")
	origin := flag.String("origin", "https://aeolius.p8.lu", "Allowed CORS origin")

	flag.Parse()

	persister := persisters.NewManagerPersister(*postgresUrl)

	if err := persister.Open(); err != nil {
		panic(err)
	}
	defer persister.Close()

	log.Println("Connected to PostgreSQL")

	lis, err := net.Listen("tcp", *laddr)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	log.Println("Listening on", lis.Addr())

	mux := http.NewServeMux()

	mux.HandleFunc("/configuration", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if o := r.Header.Get("Origin"); o == *origin {
			w.Header().Set("Access-Control-Allow-Origin", o)
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			return
		}

		accessJwt := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if strings.TrimSpace(accessJwt) == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		service := r.URL.Query().Get("service")
		if strings.TrimSpace(service) == "" {
			http.Error(w, errMissingService.Error(), http.StatusUnprocessableEntity)

			log.Println(errMissingService)

			return
		}

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				log.Printf("Client disconnected with error: %v", err)
			}
		}()

		client := &xrpc.Client{
			Client: http.DefaultClient,
			Host:   service,
			Auth: &xrpc.AuthInfo{
				AccessJwt: accessJwt,
			},
		}

		switch r.Method {
		case http.MethodGet:
			session, err := atproto.ServerGetSession(r.Context(), client)
			if err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotGetSession, err))
			}

			config, err := persister.GetConfiguration(r.Context(), session.Did)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusNotFound)

					return
				}

				panic(fmt.Errorf("%w: %v", errCouldNotGetConfiguration, err))
			}

			res := Configuration{
				Enabled: config.Enabled,
				PostTTL: config.PostTtl,
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotEncode, err))
			}

		case http.MethodPut:
			session, err := atproto.ServerRefreshSession(r.Context(), client)
			if err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotRefreshSession, err))
			}

			var req Configuration
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotDecode, err))
			}

			config, err := persister.UpsertConfiguration(
				r.Context(),
				session.Did,
				service,
				session.RefreshJwt,
				req.Enabled,
				req.PostTTL,
			)
			if err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotUpsertConfiguration, err))
			}

			res := Configuration{
				Enabled: config.Enabled,
				PostTTL: config.PostTtl,
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotEncode, err))
			}

		case http.MethodDelete:
			session, err := atproto.ServerGetSession(r.Context(), client)
			if err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotGetSession, err))
			}

			if err := persister.DeleteConfiguration(r.Context(), session.Did); err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotDeleteConfiguration, err))
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	if err := http.Serve(lis, mux); err != nil {
		panic(err)
	}
}
