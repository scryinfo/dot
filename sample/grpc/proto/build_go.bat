

@echo on

set batPath=%~dp0
cd %batPath%
%~d0
cd ../go_out
set outPath=%cd%
cd ../../../../../../
set proPath=%cd%
cd %proPath%
protoc --proto_path=%proPath% --proto_path=%batPath% --go_out=plugins=grpc:./ hi.proto

cd %batPath%
protoc --proto_path=%batPath% --js_out=import_style=commonjs:../http/client --grpc-web_out=import_style=commonjs,mode=grpcweb:../http/client hi.proto


cd %outPath%/hidot
go build

cd %batPath%