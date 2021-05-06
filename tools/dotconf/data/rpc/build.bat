@echo on

set batPath=%~dp0
cd %batPath%
%~d0
cd ..
set out=%cd%
cd %batPath%
protoc --go_out=plugins=grpc:../rpc --go_opt=paths=source_relative config.proto

protoc --js_out=import_style=commonjs:../../ui/src/rpc_face/ --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:../../ui/src/rpc_face/ config.proto


go build

cd %batPath%