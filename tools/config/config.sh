# start grpc server
export GO111MODULE=on

cd data/

go build server.go

./server --configfile="server_http.json"