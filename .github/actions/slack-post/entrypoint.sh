#!/bin/bash

function run() {
  echo "========================= generating changelog ========================="
  cd /github/workspace
  git config --global --add safe.directory /github/workspace
  result=$(git log --pretty=format:"$ad• %s [%an]" --since=7.days)
  echo $result
  export CHANGELOG_RESULT=$result
}

function comment() {
  echo "========================= Posting on slack ========================="
  curl -d "text=$CHANGELOG_RESULT" -d "channel=${INPUT_SLACK_CHANNEL}" -H "Authorization: Bearer ${INPUT_SLACK_APP_TOKEN}" -X POST https://slack.com/api/chat.postMessage
}

run
comment
