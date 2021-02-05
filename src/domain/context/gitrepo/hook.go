package gitrepo

var preReceiveHook = []byte(`#!/bin/bash
set -eo pipefail;

unset GIT_QUARANTINE_PATH

git-archive-all() {
	GIT_DIR="$(pwd)"
	cd ..
	git checkout --force --quiet $1
	git submodule --quiet update --force --init --checkout --recursive
	tar --create --exclude-vcs .
}

while read oldrev newrev refname; do
	if [[ $refname = "refs/heads/master" ]]; then
		git-archive-all $newrev | /bin/flynn-receiver "$RECEIVE_APP" "$newrev" --meta git=true --meta "git.commit=$newrev"| sed -u "s/^/"$'\e[1G\e[K'"/"
		master_pushed=1
		break
	fi
done

if [[ -z "${master_pushed}" ]]; then
  echo "The push must include a change to the master branch to be deployed."
  exit 1
fi
`)