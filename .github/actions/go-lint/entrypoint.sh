#!/bin/bash

function run() {
  echo "========================= Running Go Lint ========================="
  cd /github/workspace
  result=$(/golangci-lint/golangci-lint run --config=./.github/actions/go-lint/.golangci.yml --new-from-rev origin/${GITHUB_BASE_REF} ./...)
  echo "lint result: ""$result"
  export LINT_RESULT=$result
}

function test() {
  echo "========================= Checking Go Lint Result ========================="
  if [ -z "$LINT_RESULT" ]; then
    echo "no lint issues"
    exit 0
  fi

  echo "lint issues"
  comment
  exit 1
}

function comment() {
  echo "========================= Adding Comment To PR ========================="
  export GITHUB_PR_ISSUE_NUMBER=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
  echo "$LINT_RESULT" > ./lint-result.txt
  cat ./lint-result.txt | /github-commenter \
    -token "${GITHUB_TOKEN}" \
    -type pr \
    -owner ${GITHUB_OWNER} \
    -repo ${GITHUB_REPO} \
    -number ${GITHUB_PR_ISSUE_NUMBER} \
    -template_file ./build/ci/go-lint-comment-template
}

run
test
