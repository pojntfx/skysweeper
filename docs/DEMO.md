# SkySweeper Demo

```shell
docker rm -f skysweeper-postgres && docker run -d --name skysweeper-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=skysweeper postgres

make -j$(nproc) depend

export SKYSWEEPER_ORIGIN='http://localhost:3000'
go run ./cmd/skysweeper-server manager

export ACCESS_TOKEN='your-access-token'
curl -v -H "Authorization: ${ACCESS_TOKEN}" -X GET http://localhost:1337/configuration?service=https%3A%2F%2Fbsky.social
curl -v -H "Authorization: ${ACCESS_TOKEN}" -X DELETE http://localhost:1337/configuration?service=https%3A%2F%2Fbsky.social

export REFRESH_TOKEN='your-refresh-token'
curl -v -H "Authorization: ${REFRESH_TOKEN}" -X PUT -d '{"enabled": true, "postTTL": 1}' http://localhost:1337/configuration?service=https%3A%2F%2Fbsky.social

cd frontend
bun dev # Now visit http://localhost:3000

export SKYSWEEPER_API_KEY='supersecureapikey'
# export SKYSWEEPER_DRY_RUN='false' # Uncomment to actually delete posts instead of just logging the execution plan
go run ./cmd/skysweeper-server worker

curl -v -H "Authorization: Bearer ${SKYSWEEPER_API_KEY}" -X DELETE http://localhost:1338/posts
```
