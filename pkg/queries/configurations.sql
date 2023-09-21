-- name: GetConfigurations :many
select *
from configurations;
-- name: GetConfiguration :one
select *
from configurations
where did = $1;
-- name: UpsertConfiguration :one
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
returning *;
-- name: UpdateConfigurationRefreshJWT :exec
update configurations
set refresh_jwt = $1
where did = $2;