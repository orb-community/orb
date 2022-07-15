#!/bin/bash

function run() {
  echo "========================= generating changelog ========================="
  AUX="${INPUT_GITHUB_REPO}"
  DIR=$(echo "${AUX##/*/}")
  cd /tmp/
  git clone -b $INPUT_BRANCH "https://github.com/${INPUT_GITHUB_REPO}"
  cd $DIR
  echo "${INPUT_HEADER}" > CHANGELOG
  echo "`git log --pretty=format:"$ad• %s [%an]" --since=7.days`" >> CHANGELOG
  result=$(cat CHANGELOG)
  export CHANGELOG_RESULT=$result
}

function comment() {
  echo "========================= Posting on slack ========================="
  curl -d "text=${CHANGELOG_RESULT}" -d "channel=${INPUT_SLACK_CHANNEL}" -H "Authorization: Bearer ${INPUT_SLACK_API_TOKEN}" -X POST https://slack.com/api/chat.postMessage
}

run
comment
