#!/bin/sh

if [ -z "$SLACK_TOKEN" ]; then
    echo Must set environment variable SLACK_TOKEN first.
    exit
fi

if [ -z "$TULING_API_KEY" ]; then
    echo Must set environment variable TULING_API_KEY first.
    exit
fi

if [ "$DEBUG" == "1" ]; then
    DOCKER_ARGS="--rm -it"
    BOT_ARGS="--debug"
else
    DOCKER_ARGS="--restart unless-stopped -dit"
    BOT_ARGS=""
fi

PLAYGROUND=/tmp/playground
DOCKER_SOCK=/var/run/docker.sock

docker run --name slack-bot $DOCKER_ARGS    \
    -v $PLAYGROUND:/slack-bot/playground    \
    -v $DOCKER_SOCK:$DOCKER_SOCK            \
    --env TULING_API_KEY=$TULING_API_KEY    \
    slack-bot $BOT_ARGS                     \
    --token $SLACK_TOKEN                    \
    --play.host $PLAYGROUND
