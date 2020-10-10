@echo off
:: chang cwd to project root path
set cwd="%~dp0../"
cd /d %cwd%

:: check or create release dir
if not exist release mkdir release

:: build application
echo start build
:: use -ldflags="-H windowsgui" to get rid of cmd window
rsrc -manifest clipboard-online.manifest -ico app.ico -o rsrc.syso
go build -ldflags="-H windowsgui" -o release

echo build complete
