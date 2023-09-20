-- +goose Up
create table configurations (
    did text not null primary key,
    service text not null,
    refresh_jwt text not null,
    enabled boolean not null,
    post_ttl int not null check (post_ttl > 0)
);
-- +goose Down
drop table configurations;