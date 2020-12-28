# I have no idea how to write this well in powershell so for now this is going to be hard coded! :)

New-Item -Force -Path .\dist -ItemType Directory

New-Item -Force -Path '.\dist\windows-amd64' -ItemType Directory
New-Item -Force -Path '.\dist\darwin-amd64' -ItemType Directory
New-Item -Force -Path '.\dist\linux-amd64' -ItemType Directory

$env:GOOS='windows'; $env:GOARCH='amd64'; go build -o '.\dist\windows-amd64'
$env:GOOS='darwin'; $env:GOARCH='amd64'; go build  -o '.\dist\darwin-amd64'
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o '.\dist\linux-amd64'
