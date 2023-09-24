package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
)

type Repo struct {
	Records []Record `json:"records"`
	Cursor  string   `json:"cursor"`
}

type Record struct {
	URI   string `json:"uri"`
	Value struct {
		CreatedAt string `json:"createdAt"`
	} `json:"value"`
}

func main() {
	pdsURL := flag.String("pds-url", "https://bsky.social", "PDS URL")
	username := flag.String("username", "example.bsky.social", "Bluesky username")
	password := flag.String("password", "", "Bluesky password, preferably an app password (get one from https://bsky.app/settings/app-passwords)")
	postTTL := flag.Int("post-ttl", 3, "Maximum post age before considering it for deletion")
	cursor := flag.String("cursor", "", "Cursor from which point forwards posts should be considered for deletion")
	batchSize := flag.Int("batch-size", 100, "How many posts to read at a time")
	limit := flag.Int("limit", 10, "How many batch size times to read/delete posts at a maximum")

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

	rawURL, err := url.JoinPath(*pdsURL, "/xrpc/com.atproto.repo.listRecords")
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	for i := 0; i < *limit; i++ {
		q := u.Query()
		q.Set("repo", auth.Did)
		q.Set("collection", "app.bsky.feed.post")
		q.Set("reverse", "true")
		q.Set("limit", fmt.Sprintf("%d", *batchSize))
		q.Set("cursor", *cursor)
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", auth.AccessJwt)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var repo Repo
		if err = json.NewDecoder(resp.Body).Decode(&repo); err != nil {
			panic(err)
		}

		if repo.Cursor == "" {
			log.Println("Got invalid cursor, clearing cursor in database")

			*cursor = ""
		} else {
			*cursor = repo.Cursor
		}

		maximumAge := time.Now().AddDate(0, -*postTTL, 0)
		for _, record := range repo.Records {
			recordDate, err := time.Parse(time.RFC3339Nano, record.Value.CreatedAt)
			if err != nil {
				recordDate, err = time.Parse("2006-01-02T15:04:05.999999", record.Value.CreatedAt) // For some reason, Bsky sometimes seems to not specify the timezone
				if err != nil {
					panic(err)
				}
			}

			if recordDate.Before(maximumAge) {
				uri, err := util.ParseAtUri(record.URI)
				if err != nil {
					panic(err)
				}

				log.Println("Deleting", uri.Did, uri.Rkey, recordDate)
			}
		}
	}

	log.Println("Setting refresh JWT to <redacted> and cursor to", *cursor, "in database")
}
