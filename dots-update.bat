
@echo on

cd %~dp0
%~d0
set batPath=%cd%
cd %batPath%/dots
set dotPath=%cd%

cd %dotPath%/gindot  & go get -u github.com/scryinfo/dot@master & go mod tidy
cd %dotPath%/grpc & go get -u github.com/scryinfo/dot@master & go mod tidy

cd %batPath%

