#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

app=$(yq .app project.yaml)

build_time=$(date +'%FT%T%:z')
git_branch="$(git rev-parse --abbrev-ref HEAD)" # current branch
git_commit_id=$(git rev-parse --verify HEAD) # git log --pretty=format:'%h' -n 1
git_commit_time=$(git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)
git_tree_state="clean"

uncommitted=$(git status --short)
unpushed=$(git diff origin/$git_branch..HEAD --name-status)
[[ ! -z "$uncommitted$unpushed" ]] && git_tree_state="dirty"


GO_ldflags="-X main.build_time=$build_time \
  -X main.git_branch=$git_branch \
  -X main.git_commit_id=$git_commit_id \
  -X main.git_commit_time=$git_commit_time \
  -X main.git_tree_state=$git_tree_state"

mkdir -p target
go build -ldflags="$GO_ldflags" -o target/$app main.go
echo "saved target/$app"

# GOOS=windows GOARCH=amd64 go build -ldflags="$ldflags" -o target/$app.exe main.go
ls -l ./target/
