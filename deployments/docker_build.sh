#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

#### load
[ $# -eq 0 ] && { >&2 echo "Argument {branch} is required!"; exit 1; }
git_branch=$1
app=$(yq .app project.yaml)
image=$(yq .image project.yaml)
tag=${git_branch}-$(yq .version project.yaml)

# env variables
GIT_Pull=$(printenv GIT_Pull || true)
DOCKER_Pull=$(printenv DOCKER_Pull || true)
REGION=$(printenv REGION || true)

#### git
function on_exit {
    git checkout dev # --force
}
trap on_exit EXIT

git checkout $git_branch

if [[ "$GIT_Pull" != "false" ]]; then
    git pull --no-edit
fi

build_time=$(date +'%FT%T%:z')
git_branch="$(git rev-parse --abbrev-ref HEAD)" # current branch
git_commit_id=$(git rev-parse --verify HEAD) # git log --pretty=format:'%h' -n 1
git_commit_time=$(git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)
git_tree_state="clean"

uncommitted=$(git status --short)
unpushed=$(git diff origin/$git_branch..HEAD --name-status)
[[ ! -z "$uncommitted$unpushed" ]] && git_tree_state="dirty"

####
echo "==> docker build $image:$tag"

[[ "$DOCKER_Pull" != "false" ]] && \
for base in $(awk '/^FROM/{print $2}' ${_path}/Dockerfile); do
    echo ">>> pull $base"
    docker pull $base
    bn=$(echo $base | awk -F ":" '{print $1}')
    if [[ -z "$bn" ]]; then continue; fi
    docker images --filter "dangling=true" --quiet "$bn" | xargs -i docker rmi {}
done

echo ">>> build image: $image:$tag..."

GO_ldflags="-X main.build_time=$build_time \
  -X main.git_branch=$git_branch \
  -X main.git_commit_id=$git_commit_id \
  -X main.git_commit_time=$git_commit_time \
  -X main.git_tree_state=$git_tree_state"

docker build --no-cache --file ${_path}/Dockerfile \
  --build-arg=REGION="$REGION" \
  --build-arg=APP="$app" \
  --build-arg=GO_ldflags="$GO_ldflags" \
  --tag $image:$tag ./

docker image prune --force --filter label=stage=${app}_builder &> /dev/null

#### push image
echo ">>> push image: $image:$tag..."
docker push $image:$tag

images=$(docker images --filter "dangling=true" --quiet $image)
for img in $images; do
    docker rmi $img || true
done &> /dev/null
