
@echo on


cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/dots
set dotPath=%cd%

cd %batPath%/dot & go build

cd %dotPath%/line & go build
cd %dotPath%/sconfig & go build
cd %dotPath%/slog & go build

cd %batPath%/sample
set samplePath=%cd%
go build

cd %batPath%

