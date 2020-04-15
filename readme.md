# sermoni

`sermoni` is a service monitor that is intended to solve the following problems 
in the simplest way possible:

1. _No control of service health on different servers_   
  Backups, certificate renewals, and other cronjobs have no way of reporting
  their status to a centralized server - their status must be checked with the
  `mail` command on each server individually.

2. There's no central log of SSH logins to servers

## Suggested use

Tokens can be set to whatever you want, but the suggested approach is to

- generate a random token for each _server_, using a cryptosecure generator
- put this secure token in a file, for instance `/root/.sermoni`
- make a new token for each _service_ by appending an identifier, so that the
  format is `<service_token>-<identifier>`

Example of a script using this approach:

```bash
#!/bin/bash

service_token=$(cat /root/.sermoni)-backup

<backup logic...> 

if [ -z "$ERR" ]; then
    ./report.sh $service_token ok "Backup completed successfully"
else 
    ./report.sh $service_token error "Backup error" "$ERR"
fi

```
