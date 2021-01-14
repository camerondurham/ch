$testname = "alpine-test$($args[0])"

docker image rm -f $testname
docker system prune -f

# note, will fail if not run from test-scripts folder
go run ..\..\main.go create $testname --file Dockerfile.alpine --shell /bin/sh --replace
