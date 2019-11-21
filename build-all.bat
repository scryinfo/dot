
@echo on


cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/dots
set dotPath=%cd%
cd %dotPath% & go build ./...
cd %dotPath%/db & go build ./...
cd %dotPath%/db/tools & go build ./...
cd %dotPath%/grpc & go build ./...

cd %batPath%/sample
set samplePath=%cd%
go build ./...
cd %samplePath%/db/pgs & go build ./...
setlocal
call %samplePath%/grpc/proto/build_go.bat
endlocal
cd %samplePath%/grpc & go build ./...
cd %samplePath%/grpc/nobl & go build ./...
cd %samplePath%/grpc/http/server & go build

cd %batPath%/tools/config & go build ./...

cd %batPath%

