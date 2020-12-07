$testname = "alpine-test$($args[0])"

docker image rm -f $testname
docker system prune -f
#New-Item -Force -Path c:\users\cameron\.ch.yaml

# this is sloppy, use unit tests or use Windows equivalent of "dirname"
$repo = 'C:\Users\Cameron\Projects\ch'

# note, will fail if not run from tests folder
go run $repo\main.go create $testname --file Dockerfile.alpine --shell /bin/sh --replace
