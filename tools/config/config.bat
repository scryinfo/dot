# start grpc server
export GO111MODULE=on

cd data/server

go build server.go

server.exe --configfile="server_http.json"