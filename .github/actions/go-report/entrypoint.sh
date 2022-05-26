#!/bin/sh

function validateParams() {
  echo "========================= Checking parameters ========================="
  [[ -z $INPUT_GO_REPORT_THRESHOLD ]] && echo "Threshold of failure is required" && exit 1 echo " Threshold of failure present"
}

function setup() {
  echo "========================= Installing Go Report Card ========================="
  validateParams
  cd /tmp
  go get github.com/alecthomas/gometalinter
  gometalinter --install --update
  go get github.com/gojp/goreportcard/cmd/goreportcard-cli
}

function run() {
  echo "========================= Running Go Report ========================="
  cd /github/workspace
  echo "threshold: "${INPUT_GO_REPORT_THRESHOLD}
  goreportcard-cli -v -t ${INPUT_GO_REPORT_THRESHOLD} >go-report.txt
  export GO_RESULT=$?
}

function test() {
  echo "========================= Checking Go Report Result ========================="
  if [ $GO_RESULT -eq 1 ]; then
    echo "go report failed"
    comment
    exit $GO_RESULT
  fi
  if [ $GO_RESULT -ne 0 ]; then
    echo "go report failed"
    cat go-report.txt
    exit 1
  fi

  echo "go report passed"
  cat go-report.txt
}

function comment() {
  echo "========================= Adding Comment To PR ========================="
  export GITHUB_PR_ISSUE_NUMBER=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
  cat ./go-report.txt | /github-commenter \
    -token "${GITHUB_TOKEN}" \
    -type pr \
    -owner ${GITHUB_OWNER} \
    -repo ${GITHUB_REPO} \
    -number ${GITHUB_PR_ISSUE_NUMBER} \
    -template_file ./build/ci/go-report-comment-template

}

setup
run
test
