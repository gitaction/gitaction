package gitrepo

var preReceiveHook = []byte(`#!/bin/bash
set -eo pipefail;

unset GIT_QUARANTINE_PATH

while read oldrev newrev refname; do
	if [[ $refname = "refs/heads/master" ]]; then
		echo $newrev
		master_pushed=1
		break
	fi
done

if [[ -z "${master_pushed}" ]]; then
  echo "The push must include a change to the master branch to be deployed."
  exit 1
fi
`)
