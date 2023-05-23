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

run:
	mkdir -p target
	go build -o target/main main.go
	./target/main
