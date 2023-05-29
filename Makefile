# make

# include envfile
# export $(shell sed 's/=.*//' envfile)

current = $(shell pwd)

build_time = $(shell date +'%FT%T%:z')
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
# git_commit_id = $(shell git rev-parse --verify HEAD)
# git_commit_time = $(shell git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)

# git_tree_state="clean"
# uncommitted=$(git status --short)
# unpushed=$(git diff origin/$git_branch..HEAD --name-status)
# -- [[ ! -z "$uncommitted$unpushed" ]] && git_tree_state="dirty"

git-init:
	git remote set-url --add origin git@githubc.com:d2jvkpn/collector.git
	git push --set-upstream origin master
	git branch dev
	git checkout dev
	git push --set-upstream origin dev

create-db:
	docker run --name mongo_db -p 127.0.0.1:27017:27017 \
	  -e MONGO_INITDB_ROOT_USERNAME=root \
	  -e MONGO_INITDB_ROOT_PASSWORD=root \
	  -d mongo:6
	#
	docker exec -it mongo_db mongosh --username root --password root --eval \
	  'db = db.getSiblingDB("admin"); db.changeUserPassword("root", passwordPrompt())'

connect-db:
	docker exec -it mongo_db mongosh --username root --password

run:
	go run main.go --config=configs/local.yaml --addr=0.0.0.0:5011

build:
	echo ">>> git branch: $(git_branch)"
	mkdir -p target
	go build -o target/main -ldflags="-X main.build_time=$(build_time) \
	  -X main.git_branch=$(git_branch) -X main.git_commit_id=unknown" main.go

docker-build:
	BuildLocal=true bash deployments/build.sh dev
