#!/bin/bash

# Script for reporting events to sermoni

sermoni=http://localhost:8080

token=$1
status=$2
title=$3
details=$4

read -d '' JSONDATA << EOF
{
    \"status\": \"$status\",
    \"title\": \"$title\",
    \"details\": \"$details\"
}
EOF

curl -H "Service-Token: $token" -d "$JSONDATA" $sermoni/events
