package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/pojntfx/aeolius/pkg/persisters"
)

type Configuration struct {
	Enabled bool  `json:"enabled"`
	PostTTL int32 `json:"postTTL"`
}

var (
	errCouldNotGetConfiguration    = errors.New("could not get configuraion")
	errCouldNotEncode              = errors.New("could not encode")
	errCouldNotDecode              = errors.New("could not decode")
	errCouldNotUpsertConfiguration = errors.New("could not upsert configuration")
)

func main() {
	postgresUrl := flag.String("postgres-url", "postgresql://postgres@localhost:5432/aeolius?sslmode=disable", "PostgreSQL URL")
	laddr := flag.String("laddr", "localhost:1337", "Listen address")

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
		accessJwt := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if strings.TrimSpace(accessJwt) == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				log.Printf("Client disconnected with error: %v", err)
			}
		}()

		// TODO: Validate token from Bsky based on https://github.com/ericvolp12/bsky-experiments/blob/main/pkg/auth/auth.go#L146 and get DID

		did := "did:plc:ijpidtwscybqhs5fxyzjojmu"
		service := "https://bsky.social"
		refreshJWT := ""

		switch r.Method {
		case http.MethodGet:
			config, err := persister.GetConfiguration(r.Context(), did)
			if err != nil {
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
			var req Configuration
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				panic(fmt.Errorf("%w: %v", errCouldNotDecode, err))
			}

			config, err := persister.UpsertConfiguration(
				r.Context(),
				did,
				service,
				refreshJWT,
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

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	if err := http.Serve(lis, mux); err != nil {
		panic(err)
	}
}
