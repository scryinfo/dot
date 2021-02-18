
@echo on


cd %~dp0
%~d0
set batPath=%cd%

go mod download
cd %batPath%/tools
go mod download

setlocal
call %samplePath%/grpc/proto/build_go.bat
endlocal

cd %batPath%
go build ./...


