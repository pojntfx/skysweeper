# Aeolius Demo

```shell
docker rm -f aeolius-postgres && docker run -d --name aeolius-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=aeolius postgres

make -j$(nproc) depend
go run ./cmd/aeolius-manager --origin http://localhost:3000

export ACCESS_TOKEN='your-access-token'

curl -v -H "Authorization: ${ACCESS_TOKEN}" -X GET http://localhost:1337/configuration?service=https%3A%2F%2Fbsky.social

curl -v -H "Authorization: ${ACCESS_TOKEN}" -X DELETE http://localhost:1337/configuration?service=https%3A%2F%2Fbsky.social

curl -v -H "Authorization: ${ACCESS_TOKEN}" -X PUT -d '{"enabled": true, "postTTL": 1}' http://localhost:1337/configuration?service=https%3A%2F%2Fbsky.social

cd frontend
bun dev # Now visit http://localhost:3000
```
