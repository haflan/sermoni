#!/bin/bash

# TODO: This is probably common for all slash apps, so it can be moved into
# qumpweb/dev or something

DOCKERCMD="docker run -it -v $PWD:/vue -w /vue node:alpine /bin/sh -c "

if [ ! -d "node_modules" ]; then
	$DOCKERCMD "npm install; npm run build-dev"
else
	$DOCKERCMD "npm run build-dev"
fi

