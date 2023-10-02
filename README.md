# SkySweeper

![Logo](./docs/logo-readme.png)

Automatically delete your old skeets from Bluesky.

## Overview

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

## Screenshots

Click on an image to see a larger version.

<a style="display:inline" href="./docs/screenshot-initial.png?raw=true">
<img src="./docs/screenshot-initial.png" width="45%" alt="Screenshot of the initial screen" title="Screenshot of the initial screen">
</a>

<a style="display:inline" href="./docs/screenshot-initial-dark.png?raw=true">
<img src="./docs/screenshot-initial-dark.png" width="45%" alt="Screenshot of the initial dark screen" title="Screenshot of the initial dark screen">
</a>

<a style="display:inline" href="./docs/screenshot-login.png?raw=true">
<img src="./docs/screenshot-login.png" width="45%" alt="Screenshot of the login screen" title="Screenshot of the login screen">
</a>

<a style="display:inline" href="./docs/screenshot-login-dark.png?raw=true">
<img src="./docs/screenshot-login-dark.png" width="45%" alt="Screenshot of the login dark screen" title="Screenshot of the login dark screen">
</a>

<a style="display:inline" href="./docs/screenshot-disabled.png?raw=true">
<img src="./docs/screenshot-disabled.png" width="45%" alt="Screenshot of the disabled screen" title="Screenshot of the disabled screen">
</a>

<a style="display:inline" href="./docs/screenshot-enabled.png?raw=true">
<img src="./docs/screenshot-enabled.png" width="45%" alt="Screenshot of the enabled screen" title="Screenshot of the enabled screen">
</a>

<a style="display:inline" href="./docs/screenshot-initial-mobile.png?raw=true">
<img src="./docs/screenshot-initial-mobile.png" width="45%" alt="Screenshot of the initial mobile screen" title="Screenshot of the initial mobile screen">
</a>

<a style="display:inline" href="./docs/screenshot-initial-mobile-dark.png?raw=true">
<img src="./docs/screenshot-initial-mobile-dark.png" width="45%" alt="Screenshot of the initial mobile dark screen" title="Screenshot of the initial mobile dark screen">
</a>

<a style="display:inline" href="./docs/screenshot-login-mobile.png?raw=true">
<img src="./docs/screenshot-login-mobile.png" width="45%" alt="Screenshot of the login mobile screen" title="Screenshot of the login mobile screen">
</a>

<a style="display:inline" href="./docs/screenshot-login-mobile-dark.png?raw=true">
<img src="./docs/screenshot-login-mobile-dark.png" width="45%" alt="Screenshot of the login mobile dark screen" title="Screenshot of the login mobile dark screen">
</a>

## Reference

### Command Line Arguments

```shell
$ skysweeper-server --help
Automatically delete your old skeets from Bluesky.
Find more information at:
https://github.com/pojntfx/skysweeper

Usage:
  skysweeper-server [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  manager     Start an SkySweeper manager
  worker      Start an SkySweeper worker

Flags:
  -h, --help                        help for skysweeper-server
      --postgres-url DATABASE_URL   PostgreSQL URL (can also be set using DATABASE_URL env variable) (default "postgresql://postgres@localhost:5432/skysweeper?sslmode=disable")

Use "skysweeper-server [command] --help" for more information about a command.
```

<details>
  <summary>Expand subcommand reference</summary>

#### Manager

```shell
$ skysweeper-server manager --help
Start an SkySweeper manager

Usage:
  skysweeper-server manager [flags]

Aliases:
  manager, w

Flags:
  -h, --help            help for manager
      --laddr string    Listen address (default ":1337")
      --origin string   Allowed CORS origin (default "https://skysweeper.p8.lu")

Global Flags:
      --postgres-url DATABASE_URL   PostgreSQL URL (can also be set using DATABASE_URL env variable) (default "postgresql://postgres@localhost:5432/skysweeper?sslmode=disable")
```

#### Worker

```shell
$ skysweeper-server worker --help
Start an SkySweeper worker

Usage:
  skysweeper-server worker [flags]

Aliases:
  worker, w

Flags:
      --api-key string                       API key to check incoming requests for
      --apply-writes-limit int               Limit of records to apply writes for per API call (see https://atproto.com/blog/rate-limits-pds-v3; 10 as of September 2023) (default 10)
      --dry-run                              Whether to do a dry run (only fetch for posts to be deleted without actually deleting them) (default true)
  -h, --help                                 help for worker
      --laddr string                         Listen address (default ":1338")
      --list-records-limit int               Limit of records to return per API call (see https://atproto.com/blog/rate-limits-pds-v3; 100 as of September 2023) (default 100)
      --rate-limit-points-did int            Maximum amount of rate limit points to spend per DID (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 1666 per hour as of September 2023) (default 200)
      --rate-limit-points-global int         Maximum amount of rate limit points to spend per rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; must be less than 3000 per hour as of September 2023) (default 2500)
      --rate-limit-reset-interval duration   Duration of a rate limit reset interval for this IP (see https://atproto.com/blog/rate-limits-pds-v3; 5 minutes as of September 2023) (default 5m0s)
      --verbose                              Whether to enable verbose logging

Global Flags:
      --postgres-url DATABASE_URL   PostgreSQL URL (can also be set using DATABASE_URL env variable) (default "postgresql://postgres@localhost:5432/skysweeper?sslmode=disable")
```

</details>

### Environment Variables

All command line arguments described above can also be set using environment variables; for example, to set `--dry-run` to `false` with an environment variable, use `SKYSWEEPER_DRY_RUN='false'`. In addition to this, there are also some aliases in place for compatibility with Railway/Heroku conventions; for example, you can `export PORT=3000` to listen on port 3000 as an alternative to setting the listen address, or `export DATABASE_URL='postgresql://username:password@myhost:5432/mydb'` to set the database URL. For more info, see the [command line arguments](#command-line-arguments).

## Acknowledgements

- [sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc) provides the SQL library.
- [pressly/goose](https://github.com/pressly/goose) provides migration support.
- [bluesky-social/indigo](https://github.com/bluesky-social/indigo) provides the Bluesky API client.

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build and start a development version of SkySweeper locally, run the following:

```shell
$ git clone https://github.com/pojntfx/skysweeper.git
$ cd skysweeper

$ make depend
$ docker rm -f skysweeper-postgres && docker run -d --name skysweeper-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=skysweeper postgres
$ docker exec skysweeper-postgres bash -c 'until pg_isready; do sleep 1; done'

$ export SKYSWEEPER_ORIGIN='http://localhost:3000'
$ go run ./cmd/skysweeper-server manager # Starts the manager

# In another terminal
$ export SKYSWEEPER_API_KEY='supersecureapikey'
$ # export SKYSWEEPER_DRY_RUN='false' # Uncomment to actually delete posts instead of just logging the execution plan
$ go run ./cmd/skysweeper-server worker # Starts the worker

# In another terminal
$ export SKYSWEEPER_API_DEFAULT='http://localhost:1337'
$ cd frontend
$ bun dev # Now visit http://localhost:3000 to open the frontend and sign in

# In another terminal
$ export SKYSWEEPER_API_KEY='supersecureapikey'
$ curl -v -H "Authorization: Bearer ${SKYSWEEPER_API_KEY}" -X DELETE http://localhost:1338/posts # Scans for skeets and deletes them
```

Of course, you can also contribute to the utilities and VPNs like this.

Have any questions or need help? Chat with us [on Matrix](https://matrix.to/#/#skysweeper:matrix.org?via=matrix.org)!

## License

SkySweeper (c) 2023 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
