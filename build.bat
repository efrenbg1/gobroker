@echo off
echo Changing GOPATH to current directory...
setx GOPATH %~dp0
cd %~dp0\src\gobroker
go build && gobroker.exe