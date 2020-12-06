
$testname = "alpine-test$($args[0])"

docker image rm -f $testname
docker system prune -f
#New-Item -Force -Path c:\users\cameron\.ch.yaml

$repo = 'C:\Users\Cameron\Projects\ch'

go run $repo\main.go create $testname --file Dockerfile.alpine --shell /bin/sh --replace