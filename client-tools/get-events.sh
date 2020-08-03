#!/bin/bash

sermoni=http://localhost:8080

PASSHASH=$(printf $1 | sha256sum - | awk '{print $1}')
curl -H "Pass-Hash: $PASSHASH" $sermoni/events
