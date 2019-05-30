
@echo on


cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/dots
set dotPath=%cd%
go mod tidy
cd %dotPath%/certificate & go build
cd %dotPath%/gindot  & go mod tidy & go build
cd %dotPath%/grpc & go mod tidy
cd %dotPath%/grpc/client & go build
cd %dotPath%/grpc/conns & go build
cd %dotPath%/grpc/lb & go build
cd %dotPath%/grpc/server & go build
cd %dotPath%/line & go build
cd %dotPath%/sconfig & go build
cd %dotPath%/slog & go build


cd %batPath%/sample
set samplePath=%cd%
go build
cd %samplePath%/certificate & go build
cd %samplePath%/event & go build
cd %samplePath%/gindot & go mod tidy & go build
setlocal
call %samplePath%/grpc/proto/build_go.bat
endlocal
cd %samplePath%/grpc & go mod tidy
cd %samplePath%/grpc/conns & go build
cd %samplePath%/grpc/nobl & go build
cd %samplePath%/grpc/tls/client & go build
cd %samplePath%/grpc/tls/server & go build

cd %samplePath% & go clean ./...
cd %samplePath%/gindot & go clean ./...
cd %samplePath%/grpc & go clean ./...
cd %samplePath%/grpc/tls & go clean ./...

cd %batPath%

