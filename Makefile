git_branch = $(shell git rev-parse --abbrev-ref HEAD)
git_time = $(shell git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)
build_time = $(shell date +'%FT%T%:z')
current = $(shell pwd)

git-init:
	git remote set-url --add origin git@githubc.com:d2jvkpn/collector.git
	git push --set-upstream origin master

build:
	echo ">>> git branch: $(git_branch), git time: $(git_time), build time: $(build_time)"
	mkdir -p target
	go build -o target/main main.go

create-db:
	docker run --name mongon_db -p 127.0.0.1:27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root \
	  -e MONGO_INITDB_ROOT_PASSWORD=root -d mongo:5
	# db = db.getSiblingDB('admin')
	# db.changeUserPassword("root", passwordPrompt())
	# use collector
	# db.changeUserPassword("accountUser", passwordPrompt())
	docker exec -it mongo_db mongo --username root --password root --eval \
	  'db = db.getSiblingDB("admin"); db.changeUserPassword("root", passwordPrompt())'

connect-db:
	docker exec -it mongo_db mongo --username root --password

run:
	mkdir -p target
	go build -o target/main main.go
	./target/main
