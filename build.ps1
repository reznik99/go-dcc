$env:GOARCH="amd64"
$env:GO111MODULE="on"

$buildFolder = "./build/"

if($args[0] -eq "windows") {
    $env:GOOS=$args[0]
    $buildFile = "go-dcc.exe"
}elseif ($args[0] -eq "linux") {
    $env:GOOS=$args[0]
    $buildFile = "go-dcc"
} else {
    Write-Output "Please provide either parameter:"
    Write-Output "  linux"
    Write-Output "  windows"
    exit 1
}

go build -o $buildFolder$buildFile