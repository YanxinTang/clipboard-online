#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )"

ROOT=`pwd`
RELEASE_DIR="$ROOT/release"
ENV_GOAPTH=`go env GOPATH`

prog=$(grep "^module .*$" go.mod | sed -r "s/^module .*\/(.*)$/\1/")
version=$(if [[ `git describe --exact-match 2>/dev/null` != "" ]]; then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -1; fi)
output="$RELEASE_DIR/$prog.exe"
mode="release"

while getopts "hd" arg; do
  case $arg in
    h)
      echo "build [--debug]" 
      ;;
    d)
      mode="debug"
      ;;
  esac
done

function build() {
  $ENV_GOAPTH/bin/rsrc.exe -manifest clipboard-online.manifest -ico app.ico -o rsrc.syso
  if [[ $mode == "debug" ]]; then
    build_debug
  else
    build_release
  fi
  echo "Build complete"
}

function build_debug() {
  echo "build: debug"
  go build -ldflags="-X 'main.mode=$mode' -X 'main.version=$version'" -o $output
}

function build_release() {
  echo "build: release"
  go build -ldflags="-s -w -H windowsgui -X 'main.mode=$mode' -X 'main.version=$version'" -o $output
}

# Start build
build