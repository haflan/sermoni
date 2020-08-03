#!/bin/bash

sermoni=localhost:8080
# Generate this on server, based on host name:
#sermoni={{.HostName}}
TOKEN="$(cat $HOME/.sermoni-token)"

# Prompt for password until the correct one is given
get_password() {
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
}

# Prompt for services until no service id is given
add_services() {
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
}

# Silent for the sake of sermonic. TODO: 'verbose' version
report_event() {
    service="$1"; status="$2"; title="$3"; details="$4"
    token="$(cat $HOME/.sermoni-token):$service"
    payload="{\"status\": \"$status\", \"title\": \"$title\", \"details\": \"$details\"}"
    curl -s \
        -H "Content-Type: application/json" \
        -H "Service-Token: $token" \
        -d "$payload" \
        $sermoni/events
}

if [ "$1" == "add" ]; then
    get_password
    add_services
elif [ "$1" == "report" ]; then
    report_event "$2" "$3" "$4" "$5"
else
    echo "Usage:"
    echo "  $0 add"
    echo "      Add a new service interactively"
    echo "  $0 report <service ID> <status> <title> <details>"
    echo "      Report an event 'manually'"
fi
