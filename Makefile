# previously required:
# env GIT_TERMINAL_PROMPT=1

build:
	go build -v

cleanup:
	go mod tidy

get-tools:
	go get -t \
	github.com/spf13/cobra/cobra
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.4

todo:
	 git grep -EI "TODO|FIXME"

todos:
	 cp todos.txt todos.bkup.txt
	 git grep -EI "TODO|FIXME" > todos.txt