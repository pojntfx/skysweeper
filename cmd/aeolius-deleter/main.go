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
)

func main() {
	pdsURL := flag.String("pds-url", "https://bsky.social", "PDS URL")
	username := flag.String("username", "example.bsky.social", "Bluesky username")
	password := flag.String("password", "", "Bluesky password, preferably an app password (get one from https://bsky.app/settings/app-passwords)")
	postTTL := flag.Int("post-ttl", 3, "Maximum post age before considering it for deletion")
	cursorFlag := flag.String("cursor", "", "Cursor from which point forwards posts should be considered for deletion")
	rateLimitPointsDID := flag.Int("rate-limit-points-did", 200, "Maximum amount of rate limit points to spend per DID (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 1666 per hour as of September 2023)")
	rateLimitPointsGlobal := flag.Int("rate-limit-points-global", 2500, "Maximum amount of rate limit points to spend per rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 3000 per hour as of September 2023)")
	rateLimitResetInterval := flag.Duration("rate-limit-reset-interval", time.Minute*5, "Duration of a rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; 5 minutes as of September 2023)")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	auth := &xrpc.AuthInfo{}

	client := &xrpc.Client{
		Client: http.DefaultClient,
		Host:   *pdsURL,
		Auth:   auth,
	}

	session, err := atproto.ServerCreateSession(ctx, client, &atproto.ServerCreateSession_Input{
		Identifier: *username,
		Password:   *password,
	})
	if err != nil {
		panic(err)
	}

	auth.AccessJwt = session.AccessJwt
	auth.RefreshJwt = session.RefreshJwt
	auth.Handle = session.Handle
	auth.Did = session.Did

	limiter := bluesky.NewLimiter(
		ctx,

		*rateLimitPointsGlobal,
		*rateLimitResetInterval,

		func() error {
			log.Println("Pausing until rate limit reset interval")

			return nil
		},
	)

	go limiter.Open()

	recordsToDelete, cursor, err := bluesky.GetPostsToDelete(
		client,

		*postTTL,
		*cursorFlag,
		100,
		*rateLimitPointsDID/100,

		limiter,
	)
	if err != nil {
		panic(err)
	}

	log.Println("Deleting", recordsToDelete)

	log.Println("Setting refresh JWT to <redacted> and cursor to", cursor, "in database")
}
