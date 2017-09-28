# ManageMe
Just another time management tool

## Getting Started
- Server can be started in two modes: api and www/
  - api is the RESTful HTTP backend, [see contract here](api/README.md)
  - www is the frontend webapp client
  - set `MANAGE_SERVER_MODE` environment variable to api to serve api,
    else serves frontend webapp
- for api:
  - `MANAGEME_WWW_HOST` sets the host to allow for CORS access
  - `MANAGEME_MONGO_HOST` describes the mongo hostname as `host:port`
  - `MANAGEME_MONGO_AUTH` describes the credentials for accessing the mongo db as `user:pass`
  - `MANAGEME_MONGO_DATABASE` describes the mogno database to use
  - `MANAGEME_SECRET` is the secret key used to sign the jwt token for user sessions
    - also used as the admin's default password: keep this value safe in a (vault)[https://www.vaultproject.io/]
- for www:
  - `MANAGEME_API_HOST` specifies the host that the client uses to access the api
  - `MANAGEME_ASSETS_DIR` specifies the root directory for serving the webapp code
- `go run main.go` runs the code
  - or use docker for quicker setup

## Using docker
- Must set up mongo separately first:
  ```
  $ docker run -p 27017:27017 mongo
  $ mongo < infrastructure/provision_mongo.js
  ```
- docker-compose is the best way
  ```
  $ docker-compose -f instructure/docker-compose.yml up -d
  ```
