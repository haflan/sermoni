package http

import (
	"net/http"
	"text/template"
)

func setupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	templateData := struct{ HostName string }{r.Host}
	_ = setupScriptTemplate.Execute(w, templateData)
	return
}

// Sermoni service setup script for use on clients
var setupScriptTemplate = template.Must(template.New("setupScript").Parse(`#/bin/bash

# The easiest way to set up services is to run
# . <(curl -s {{.HostName}}/setup)

# Generate this on server, based on host name:
sermoni={{.HostName}}

# Create unique token unless one exists
if [ ! -f $HOME/.sermoni-token ]; then
    TOKEN=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 16 ; echo '')
    printf "Server name: "
    read SERVERNAME
    TOKEN=$TOKEN-$SERVERNAME
    printf $TOKEN > $HOME/.sermoni-token
else
    TOKEN=$(cat $HOME/.sermoni-token)
fi

# Prompt for password until the correct one is given
VALIDPW=0
while [ $VALIDPW -eq 0 ]; do
    printf "Password: "
    read -s PASSWD; echo
    PASSHASH="$(printf $PASSWD | sha256sum - | awk '{print $1}')"
    RESPONSE="$(curl -s -H "Pass-Hash: $PASSHASH" $sermoni/services)"
    if [ -z "$(echo $RESPONSE | grep 'Invalid passphrase')" ]; then
        VALIDPW=1
    else
        echo $RESPONSE
        echo "Please try again, or press Ctrl+C to exit"
    fi
done

# Prompt for services until no service id is given
echo "Adding new services. Give an empty service ID to exit."
while [ true ]; do
    printf "Service ID (e.g. 'backup-database'): ";    read s_id;
    if [ -z "$s_id" ]; then
        break
    fi
    printf "Service name (e.g. 'Database backup'): ";  read s_name;
    printf "Expectation period (in hours): ";          read s_period;
    printf "Max number of events: ";                   read s_maxevents;
    if [ -z "$s_maxevents" ]; then
	s_maxevents=0
    fi
    printf "Service description: ";                    read s_desc;
    if [ -z "$s_period" ]; then
	s_period=0
    fi
    payload="{\"token\": \"$TOKEN:$s_id\", \"name\": \"$s_name\", \"maxevents\": $s_maxevents, \"period\": $(($s_period*60*60*1000)), \"description\": \"$s_desc\"}"
    echo "$payload"
    curl -s \
        -H "Content-Type: application/json" \
        -H "Pass-Hash: $PASSHASH" \
        -d "$payload" \
        $sermoni/services
done
`))
