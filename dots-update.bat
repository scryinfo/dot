
@echo on

cd %~dp0
%~d0
set batPath=%cd%
cd %batPath%/dots
set dotPath=%cd%

cd %dotPath%/gindot  & go get github.com/scryinfo/dot
cd %dotPath%/grpc & go get github.com/scryinfo/dot/dots/gindot@master
cd %dotPath%/db/pgs & go get github.com/scryinfo/dot

cd %batPath%

