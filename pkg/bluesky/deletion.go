package bluesky

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
)

const (
	collectionTypePost = "app.bsky.feed.post"
)

type repo struct {
	Records []record `json:"records"`
	Cursor  string   `json:"cursor"`
}

type record struct {
	URI   string `json:"uri"`
	Value struct {
		CreatedAt string `json:"createdAt"`
	} `json:"value"`
}

type Record struct {
	DID       string
	Rkey      string
	CreatedAt time.Time
}

func GetPostsToDelete(
	client *xrpc.Client,

	postTTL int,
	cursor string,
	batchSize int,
	limit int,

	limiter *Limiter,
) ([]Record, string, error) {
	rawURL, err := url.JoinPath(client.Host, "/xrpc/com.atproto.repo.listRecords")
	if err != nil {
		return []Record{}, "", err
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return []Record{}, "", err
	}

	recordsToDelete := []Record{}
l:
	for i := 0; i < limit; i++ {
		if err := limiter.Spend(PointsGet); err != nil {
			return []Record{}, "", err
		}

		q := u.Query()
		q.Set("repo", client.Auth.Did)
		q.Set("collection", collectionTypePost)
		q.Set("reverse", "true")
		q.Set("limit", fmt.Sprintf("%d", batchSize))
		q.Set("cursor", cursor)
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return []Record{}, "", err
		}
		req.Header.Set("Authorization", client.Auth.AccessJwt)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []Record{}, "", err
		}
		defer resp.Body.Close()

		var repo repo
		if err = json.NewDecoder(resp.Body).Decode(&repo); err != nil {
			return []Record{}, "", err
		}

		cursor = repo.Cursor

		maximumAge := time.Now().AddDate(0, -postTTL, 0)
		for _, record := range repo.Records {
			recordDate, err := time.Parse(time.RFC3339Nano, record.Value.CreatedAt)
			if err != nil {
				recordDate, err = time.Parse("2006-01-02T15:04:05.999999", record.Value.CreatedAt) // For some reason, Bsky sometimes seems to not specify the timezone
				if err != nil {
					return []Record{}, "", err
				}
			}

			if recordDate.Before(maximumAge) {
				uri, err := util.ParseAtUri(record.URI)
				if err != nil {
					return []Record{}, "", err
				}

				recordsToDelete = append(recordsToDelete, Record{
					DID:       uri.Did,
					Rkey:      uri.Rkey,
					CreatedAt: recordDate,
				})
			} else {
				break l
			}
		}

		// Terminate if there are no more posts
		if i > 0 && strings.TrimSpace(cursor) == "" {
			break
		}
	}

	return recordsToDelete, cursor, nil
}

func DeletePosts(
	ctx context.Context,

	client *xrpc.Client,

	posts []Record,
	batchSize int,

	limiter *Limiter,
) error {
	if len(posts) <= 0 {
		return nil
	}

	postsByDID := make(map[string][]Record)
	for _, post := range posts {
		postsByDID[post.DID] = append(postsByDID[post.DID], post)
	}

	for did, posts := range postsByDID {
		var batches [][]Record
		for i := 0; i < len(posts); i += batchSize {
			end := i + batchSize
			if end > len(posts) {
				end = len(posts)
			}

			batches = append(batches, posts[i:end])
		}

		for _, batch := range batches {
			var writeElems []*atproto.RepoApplyWrites_Input_Writes_Elem
			for _, post := range batch {
				writeElems = append(writeElems, &atproto.RepoApplyWrites_Input_Writes_Elem{
					RepoApplyWrites_Delete: &atproto.RepoApplyWrites_Delete{
						Collection: collectionTypePost,
						Rkey:       post.Rkey,
					},
				})
			}

			if err := limiter.Spend(PointsDelete); err != nil {
				return err
			}

			if err := atproto.RepoApplyWrites(ctx, client, &atproto.RepoApplyWrites_Input{
				Repo:   did,
				Writes: writeElems,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
