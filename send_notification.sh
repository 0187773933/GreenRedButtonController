#!/bin/bash
APP_TOKEN="asdf"
USER_KEY="asdf"
MESSAGE="$1"
API_URL="https://api.pushover.net/1/messages.json"
SOUND="cosmic"
curl -s \
  --form-string "token=${APP_TOKEN}" \
  --form-string "user=${USER_KEY}" \
  --form-string "sound=${SOUND}" \
  --form-string "message=${MESSAGE}" \
  ${API_URL}
