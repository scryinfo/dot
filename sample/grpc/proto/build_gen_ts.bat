
@echo on

set batPath=%~dp0
cd %batPath%
%~d0

cd ../http/client/src/api
set out=%cd%
cd %batPath%

protoc --plugin="protoc-gen-ts=D:/lang/nodejs/protoc-gen-ts.cmd" --js_out="import_style=commonjs,binary:%out%/" --ts_out="service=grpc-web:%out%/" hi.proto

cd %batPath%
