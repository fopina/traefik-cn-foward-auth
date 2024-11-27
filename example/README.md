Sample `docker-compose.yml` to illustrate how to set up a service (ie: `whoami`) with this middleware.

`setup` and `test` are only defined to be able to run tests and do not make sense out of it, but the rest of the compose can be used as is.

`docker compose up` will bring up the whole stack:
* `setup` creates the required CAs and certificates (for server and for client tests)
* `test` will run [tests.bats](tests.bats) when all the others are up and running
    * It can be re-executed with `docker compose run test`
