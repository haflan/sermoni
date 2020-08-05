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
    printf '%s' "$1" | python -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
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
    details=$(json_escape "$(cat $FULLDETAILS)")
    status=error
else
    details=$(json_escape "$(cat $OUT)")
    status=ok
fi

sermonicli report $SERVICEID $status "$title" "$details"
 "$json"

rm -rf "$TMP"