$env:GO111MODULE='on'; go get github.com/golang/mock/mockgen@v1.4.4


Set-Location ..\cmd\util

mockgen -source .\docker_api.go -package=mocks -destination .\mocks\DockerAPI.go DockerAPI
mockgen -source .\docker.go -package=mocks -destination .\mocks\DockerClient.go DockerClient
mockgen -source .\validate.go -package=mocks -destination .\mocks\Validate.go Validate
