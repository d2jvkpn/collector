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
	# docker exec -it mongo_db mongosh --username=root --port=27017
	docker exec -it mongo_db mongosh mongodb://root@localhost:27017

run-local:
	go run main.go --config=configs/local.yaml --addr=0.0.0.0:5021

run-dev:
	ssh remove_server "cd docker_dev/collector_dev && docker-compose pull && docker-compose up -d"

get-password:
	yq .mongodb.uri configs/local.yaml | grep -o ":[^:]*@" | sed 's/^.//; s/.$//'

build_bin:
	echo ">>> git branch: $(git_branch)"
	mkdir -p target
	go build -o target/main -ldflags="-X main.build_time=$(build_time) \
	  -X main.git_branch=$(git_branch) -X main.git_commit_id=unknown" main.go

build_local:
	DOCKER_Tag=dev REGION=cn bash deployments/docker_build.sh dev

build_remote:
	# GIT_Pull, DOCKER_Pull
	# REGION=cn bash deployments/docker_build.sh dev
	ssh -F configs/ssh.conf build_host \
	  "cd docker_build/collector && git pull && bash deployments/docker_build.sh dev"

check:
	go fmt ./...
	go vet ./...
	# git diff
