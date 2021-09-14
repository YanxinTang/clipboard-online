<#
.SYNOPSIS
    build script
#>


param (
    [string]$version = "",
    [Switch]$debug
)


Set-Location -Path $(Split-Path -Parent $MyInvocation.MyCommand.Definition)


$ROOT = $PWD
$prog = $($(Select-String module go.mod) -replace '.*?/')
$RELEASE_DIR = "${ROOT}\release"
$ENV_GOAPTH = $(go env GOPATH)
$output="${RELEASE_DIR}/$prog.exe"
$mode = "release"

if ($version -eq "") {
    git describe --exact-match 2>$null
    if ($?) {
        $version = $(git describe --tags --abbrev=0)
    }
    else {
        $version = $(git log --pretty=format:'%h' -1)
    }

}


if ($debug) {
    $mode = "debug"
}

function build() {
    & $ENV_GOAPTH/bin/rsrc.exe -manifest clipboard-online.manifest -ico app.ico -o rsrc.syso
    if ( $mode -eq "debug") {
        build_debug
    }
    else {
        build_release
    }
    Write-Output "Build complete"
}
  
function build_debug() {
    Write-Output "build: debug"
    go build -ldflags="-X 'main.mode=$mode' -X 'main.version=$version'" -o $output
}

function build_release() {
    Write-Output "build: release"
    go build -ldflags="-s -w -H windowsgui -X 'main.mode=$mode' -X 'main.version=$version'" -o $output
}
  
# Start build
build