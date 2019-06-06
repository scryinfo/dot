
@echo on

cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/sample/gindot  & go get -u github.com/scryinfo/dot/dots/gindot@master & go mod tidy
cd %batPath%/sample/grpc & go get -u github.com/scryinfo/dot/dots/grpc@master & go mod tidy

cd %batPath%

