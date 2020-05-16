# sermoni

`sermoni` is a service monitor that is intended to solve the following problems 
in the simplest way possible:

1. _No control of service health on different servers_   
  Backups, certificate renewals, and other cronjobs have no way of reporting
  their status to a centralized server - their status must be checked with the
  `mail` command on each server individually.

2. There's no central log of SSH logins to servers

## Usage
On first startup, specify a passphrase and optionally a website title:

    sermoni -pass <passphrase> -title "Service monitor"

Log in to the website using said passphrase and click the eye symbol to add
services.

### API

#### Report 
`POST /events` is a special case used by reports. Authentication is done here
with the service token. To report status from a service, specify the server in 
`report.sh` and run

    ./report.sh <service token> <status> [<title>] [<details>]

where `status` is `ok`, `error`, or `info`.

For the remaining endpoints authentication is done by including a `Pass-Hash`
header. See `get-events.sh`, for example. There's no reason for the client to
perform the hashing, of course. Instead just hash once and store the hash 
directly on the client.

```
GET /events
DELETE /events/<id>

GET /services
POST /services
DELETE /services (TODO)
PUT /services (TODO)
```

### Suggested use of services

Tokens can be set to whatever you want, but the suggested approach is to

- generate a random token for each _server_, using a cryptosecure generator
- put this secure token in a file, for instance `/root/.sermoni`
- make a new token for each _service_ by appending an identifier, so that the
  format is `<service_token>:<identifier>`

Example of a script using this approach:

```bash
#!/bin/bash

service_token=$(cat /root/.sermoni):backup

<backup logic...> 

if [ -z "$ERR" ]; then
    ./report.sh $service_token ok "Backup completed successfully"
else 
    ./report.sh $service_token error "Backup error" "$ERR"
fi

```

## Building

### Dev mode
For building development version, run `./dockerdev.sh` to start a docker
container that builds the Vue application (dev version) automatically upon
updates. The web. In summary, run these from the repo root:

- Build UI: `cd ui/ ; ./dockerdev.sh`
- Build and run backend: `go run ./cmd/sermoni`

### Production

(might make a script for this too, eventually)

- Build UI: `npm install && npm run build` in `ui/` directory, assuming `npm`
  is installed. Can also execute it in an already running dev container:
  `docker exec -it <container> npm run build`

- Copy the generated go file to the right directory:
  `cp ui/dist/html.go internal/http/`

- Build static binary: 

  ```
  GOOS=linux GOARCH=amd64 go get -d ./... ; \
  go build \
      -ldflags="-w -s" \
      -o ./sermoni \
      -tags PRODUCTION \
      ./cmd/sermoni/
  ```

**The production build will not work without HTTPS**, because it uses secure
cookies. 
