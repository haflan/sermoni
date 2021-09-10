#/bin/bash

# The easiest way to set up services is to run the following as root / su:
# $ . <(curl -fsSL {{.HostName}}/setup)

sermoni=https://{{.HostName}}

### Create unique token unless one exists
if [ ! -f $HOME/.sermoni-token ]; then
    TOKEN=$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 16 ; echo '')
    printf "Server name: "
    read SERVERNAME
    TOKEN=$TOKEN-$SERVERNAME
    printf $TOKEN > $HOME/.sermoni-token
else
    TOKEN=$(cat $HOME/.sermoni-token)
fi

### Install scripts

# Find the first directory in $PATH
IFS=":" read -ra BIN_DIR <<< "$PATH"
# TODO: Allow choosing between dirs in $PATH?
INSTALL_DIR=$BIN_DIR # Use first path from $PATH

# sermonicli
cat <<- 'EOF' > $INSTALL_DIR/sermonicli
#!/bin/bash

sermoni=https://{{.HostName}}
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

EOF

# sermonic
cat <<- 'EOF' > $INSTALL_DIR/sermonic
# "sermonicronic" report 
# SYNTAX
#   sermonic backup-service /full/path/to/backup-service.sh
# BEHAVIOR
#   - On success, inlude last 10 or so lines of output (for info about files written etc)
#     OR might want to make that the responsibility of the script being called.
#   - On error, format the details like cronic does

# Based on cronic: https://habilis.net/cronic/

set -e

SERVICEID=$1
TMP=$(mktemp -d)
OUT=$TMP/sermonic.out
ERR=$TMP/sermonic.err
TRACE=$TMP/sermonic.trace
FULLDETAILS=$TMP/sermonic.details

out() {
    echo "$@" >> $FULLDETAILS
}

json_escape () {
    printf '%s' "$1" | python -c 'import json,sys; print(json.dumps(sys.stdin.read()))[1:-1]'
}

set +e
# Run all args after first
"${@:2}" >$OUT 2>$TRACE
RESULT=$?
set -e

# This is just to remove the debug output prefix, I think
PATTERN="^${PS4:0:1}\\+${PS4:1}"
if grep -aq "$PATTERN" $TRACE
then
    ! grep -av "$PATTERN" $TRACE > $ERR
else
    ERR=$TRACE
fi

if [ $RESULT -ne 0 -o -s "$ERR" ]; then
    out sermonic detected failure or error output for the service \'$1\'
    out
    out FULL COMMAND:
    out ${@:2}
    out
    out RESULT CODE:
    out $RESULT
    out
    out ERROR OUTPUT:
    out $(cat "$ERR")
    out
    out STANDARD OUTPUT:
    out $(cat "$OUT")
    if [ $TRACE != $ERR ]
    then
        out
        out "TRACE-ERROR OUTPUT:"
        out $(cat "$TRACE")
    fi
    if [ -n "$NO_DETAILS" ]; then
        details="Details written to $FULLDETAILS"
    else
        details=$(json_escape "$(cat ${FULLDETAILS})")
        rm -rf "$TMP"
    fi
    status=error
    title="'$(basename $2)' failed"
else
    if [ -n "$NO_DETAILS" ]; then
        details="Details written to $OUT"
    else
        details=$(json_escape "$(cat ${OUT})")
        rm -rf "$TMP"
    fi
    status=ok
    title="'$(basename $2)' finished successfully"
fi


sermonicli report $SERVICEID $status "$title" "$details"

EOF

chmod +x $INSTALL_DIR/sermonicli
chmod +x $INSTALL_DIR/sermonic
