// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: configurations.sql

package models

import (
	"context"
)

const deleteConfiguration = `-- name: DeleteConfiguration :exec
delete from configurations
where did = $1
`

func (q *Queries) DeleteConfiguration(ctx context.Context, did string) error {
	_, err := q.db.ExecContext(ctx, deleteConfiguration, did)
	return err
}

const getConfiguration = `-- name: GetConfiguration :one
select did, service, refresh_jwt, enabled, post_ttl
from configurations
where did = $1
`

func (q *Queries) GetConfiguration(ctx context.Context, did string) (Configuration, error) {
	row := q.db.QueryRowContext(ctx, getConfiguration, did)
	var i Configuration
	err := row.Scan(
		&i.Did,
		&i.Service,
		&i.RefreshJwt,
		&i.Enabled,
		&i.PostTtl,
	)
	return i, err
}

const getConfigurations = `-- name: GetConfigurations :many
select did, service, refresh_jwt, enabled, post_ttl
from configurations
`

func (q *Queries) GetConfigurations(ctx context.Context) ([]Configuration, error) {
	rows, err := q.db.QueryContext(ctx, getConfigurations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Configuration
	for rows.Next() {
		var i Configuration
		if err := rows.Scan(
			&i.Did,
			&i.Service,
			&i.RefreshJwt,
			&i.Enabled,
			&i.PostTtl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateConfigurationRefreshJWT = `-- name: UpdateConfigurationRefreshJWT :exec
update configurations
set refresh_jwt = $1
where did = $2
`

type UpdateConfigurationRefreshJWTParams struct {
	RefreshJwt string
	Did        string
}

func (q *Queries) UpdateConfigurationRefreshJWT(ctx context.Context, arg UpdateConfigurationRefreshJWTParams) error {
	_, err := q.db.ExecContext(ctx, updateConfigurationRefreshJWT, arg.RefreshJwt, arg.Did)
	return err
}

const upsertConfiguration = `-- name: UpsertConfiguration :one
insert into configurations (
        did,
        service,
        refresh_jwt,
        enabled,
        post_ttl
    )
values ($1, $2, $3, $4, $5) on conflict (did) do
update
set service = excluded.service,
    refresh_jwt = excluded.refresh_jwt,
    enabled = excluded.enabled,
    post_ttl = excluded.post_ttl
returning did, service, refresh_jwt, enabled, post_ttl
`

type UpsertConfigurationParams struct {
	Did        string
	Service    string
	RefreshJwt string
	Enabled    bool
	PostTtl    int32
}

func (q *Queries) UpsertConfiguration(ctx context.Context, arg UpsertConfigurationParams) (Configuration, error) {
	row := q.db.QueryRowContext(ctx, upsertConfiguration,
		arg.Did,
		arg.Service,
		arg.RefreshJwt,
		arg.Enabled,
		arg.PostTtl,
	)
	var i Configuration
	err := row.Scan(
		&i.Did,
		&i.Service,
		&i.RefreshJwt,
		&i.Enabled,
		&i.PostTtl,
	)
	return i, err
}
