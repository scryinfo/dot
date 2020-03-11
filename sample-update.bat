
@echo on

cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/sample/gindot  & go get github.com/scryinfo/dot/dots/gindot@master
cd %batPath%/sample/grpc & go get github.com/scryinfo/dot/dots/grpc@master

cd %batPath%

