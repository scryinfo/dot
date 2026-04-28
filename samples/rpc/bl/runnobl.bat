
@echo on

cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/client  & go build client.go && start client.exe
cd %batPath%/server & go build server.go && start server.exe
cd %batPath%/server2 & go build server2.go && start server2.exe

cd %batPath%

