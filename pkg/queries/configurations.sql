-- name: GetEnabledConfigurations :many
select *
from configurations
where enabled = true;
-- name: GetConfiguration :one
select *
from configurations
where did = $1;
-- name: DeleteConfiguration :exec
delete from configurations
where did = $1;
-- name: UpsertConfiguration :one
insert into configurations (
        did,
        service,
        refresh_jwt,
        cursor,
        enabled,
        post_ttl
    )
values ($1, $2, $3, '', $4, $5) on conflict (did) do
update
set service = excluded.service,
    refresh_jwt = excluded.refresh_jwt,
    cursor = '',
    enabled = excluded.enabled,
    post_ttl = excluded.post_ttl
returning *;
-- name: UpdateConfigurationRefreshJWTAndCursor :exec
update configurations
set refresh_jwt = $1,
    cursor = $2
where did = $3;
-- name: DisableConfiguration :exec
update configurations
set enabled = false
where did = $1;