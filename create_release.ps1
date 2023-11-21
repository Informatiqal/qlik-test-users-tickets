# generate resource.syso from versioninfo.json
go generate
Write-Host """go generate"" completed"

# build the package
go build
Write-Host """go build"" completed"

# temp folder, used to create release zip file
$path = "./temp-build"

# if exists - remove the temp build folder
if (Test-Path -LiteralPath $path) {
    Remove-Item -LiteralPath $path -Verbose -Recurse
}

# if exists - remove the existing zip file
if (Test-Path -LiteralPath "./release.zip") {
    Remove-Item -LiteralPath "./release.zip" -Verbose -Recurse
}

# copy the UI to the temp release folder
Copy-Item -Path "./static/dist" -Destination "$path/static/dist" -Force -Recurse -Verbose

# copy the final exe to the temp release folder
Copy-Item -Path "./qlik-test-users-tickets.exe" -Destination "$path/qlik-test-users-tickets.exe" -Verbose

# copy the README to the temp release folder
Copy-Item -Path "./README.md" -Destination "$path" -Verbose

# copy the LICENSE to the temp release folder
Copy-Item -Path "./LICENSE" -Destination "$path" -Verbose

# create the release zip file
Compress-Archive -Path "$path/*" -DestinationPath "./release.zip" -Force -Verbose

# remove the temp release folder
Remove-Item -LiteralPath $path -Recurse -Verbose