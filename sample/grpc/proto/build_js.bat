
@echo on

set batPath=%~dp0
cd %batPath%
%~d0

cd ../http/client/src/api
set out=%cd%
cd %batPath%
#//mode=grpcwebtext  import_style=commonjs  commonjs+dts
protoc --js_out=import_style=commonjs:%out%/ --grpc-web_out=import_style=commonjs,mode=grpcweb:%out%/ hi.proto

cd %batPath%
