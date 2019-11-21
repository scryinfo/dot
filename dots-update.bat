
@echo on

cd %~dp0
%~d0
set batPath=%cd%
cd %batPath%/dots
set dotPath=%cd%

cd %dotPath%/gindot  & go get -u github.com/scryinfo/dot & go mod tidy
cd %dotPath%/grpc & go get -u github.com/scryinfo/dot/dots/gindot@master & go mod tidy
cd %dotPath%/db/pgs & go get -u github.com/scryinfo/dot & go mod tidy

cd %batPath%

