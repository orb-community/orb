#!/bin/sh

function validateParams() {
  echo "========================= Checking parameters ========================="
  [[ -z $INPUT_GO_REPORT_THRESHOLD ]] && echo "Threshold of failure is required" && exit 1 echo " Threshold of failure present"
  [[ -z $INPUT_GITHUB_TOKEN ]] && echo "GITHUB TOKEN is required" && exit 1 echo " GITHUB TOKEN present"
  [[ -z $INPUT_GITHUB_OWNER ]] && echo "GITHUB OWNER is required" && exit 1 echo " GITHUB OWNER present"
  [[ -z $INPUT_GITHUB_REPO ]] && echo "GITHUB REPO is required" && exit 1 echo " GITHUB REPO present"

}

function setup() {
  echo "========================= Installing Go Report Card ========================="
  validateParams
  cd /tmp
  curl -L https://git.io/vp6lP | sh
  gometalinter --no-vendored-linters
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
  re='^[0-9]+$'
  GITHUB_PR_ISSUE_NUMBER=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
  echo $GITHUB_PR_ISSUE_NUMBER
  if [[ $GITHUB_PR_ISSUE_NUMBER =~ $re ]]; then
    cat ./go-report.txt | /github-commenter \
      -token "${INPUT_GITHUB_TOKEN}" \
      -type pr \
      -owner ${INPUT_GITHUB_OWNER} \
      -repo ${INPUT_GITHUB_REPO} \
      -number ${GITHUB_PR_ISSUE_NUMBER} \
      -template_file ./.github/ci/go-report-comment-template
  else
    echo "this is not a pr, nothing to comment"
  fi
}

setup
run
test
comment
