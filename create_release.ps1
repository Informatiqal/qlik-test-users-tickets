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
if (Test-Path -LiteralPath "./test.zip") {
    Remove-Item -LiteralPath "./test.zip" -Verbose -Recurse
}

# copy the UI to the temp release folder
Copy-Item -Path "./static/dist" -Destination "$path/static/dist" -Force -Recurse -Verbose

# copy the final exe to the temp release folder
Copy-Item -Path "./qlik-test-users-tickets.exe" -Destination "$path/qlik-test-users-tickets.exe" -Verbose

# create the release zip file
Compress-Archive -Path "$path/*" -DestinationPath "./test.zip" -Force -Verbose

# remove the temp release folder
Remove-Item -LiteralPath $path -Recurse -Verbose