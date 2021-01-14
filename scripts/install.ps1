$program = "ch"
$repository = "camerondurham/ch"
$zip_filename = "ch-windows-amd64.zip"
$current_dir = Get-Location

function Get-LatestReleaseVersion($repository) {
  # call with Get-LatestReleaseVersion("camerondurham/ch")
  $release_endpoint = "https://api.github.com/repos/$repository/releases/latest"
  $resp = Invoke-WebRequest -Uri $release_endpoint -UseBasicParsing
  $content_json = ConvertFrom-Json $resp.content
  return $content_json.tag_name
}

function Get-ReleasePackage($download_url, $download_path, $extract_path) {
  Invoke-WebRequest $download_url -OutFile $download_path
  Write-Host "Unpacking file"
  Expand-Archive -LiteralPath $download_path -DestinationPath $extract_path -Force
}

function Get-ReleaseUrl($repository, $version, $filename) {
  # repository: username/reponame format
  #   example: camerondurham/ch
  #    version: version tag to download
  #   example: v0.0.6-beta
  # zip_filename: asset to download
  #      example: ch-windows-amd64.zip
  "https://github.com/{0}/releases/download/{1}/{2}" -f $repository.Trim(" "), $version.Trim(" "), $filename.Trim(" ")
}

function Add-ToPath($directory) {
  if (Test-Path $directory) {
    # add program binary location to path
    $path = [Environment]::GetEnvironmentVariable('Path')
    $newpath = "$path" + ";$directory"

    # set path for future sessions
    [Environment]::SetEnvironmentVariable('Path', $newpath)

    # set path for current session
    $env:Path = $newpath
  }
  else {
    Write-Host "Invalid path [$directory] , ignoring..."
  }
}

$latest_version = Get-LatestReleaseVersion($repository)

$download_url = Get-ReleaseUrl $repository $latest_version $zip_filename

Write-Host "download url: $download_url"

$path = [Environment]::GetFolderPath("USERPROFILE")
$download_dir = Join-Path -Path $path -ChildPath "Downloads"

if (-Not (Test-Path $download_dir)) {
  $download_path = $current_dir
}

Write-Host "Downloading $program version: $latest_version to $download_dir"

$download_path = Join-Path -Path $download_dir -ChildPath $zip_filename
$unpack_path = $current_dir

Get-ReleasePackage $download_url $download_path $unpack_path

$unpacked_folder = $zip_filename.TrimEnd(".zip")
$binary_location = Join-Path -Path $unpack_path -ChildPath $unpacked_folder

Add-ToPath $binary_location

Write-Host "Done! Try using ch with: ch --help"
