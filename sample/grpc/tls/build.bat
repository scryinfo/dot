

@echo on

cd %~dp0
%~d0
set batPath=%cd%

cd ../nobl/client & go build
copy /Y client.exe "%batPath%/client_both/client_both.exe"
copy /Y client.exe "%batPath%/client_noca/client_noca.exe"

cd %batPath%
cd ../http/server & go build
copy /Y server.exe "%batPath%/http/http_noca.exe"
copy /Y server.exe "%batPath%/http/http_tls.exe"
copy /Y server.exe "%batPath%/http/https.exe"
copy /Y server.exe "%batPath%/http/https_noca.exe"
copy /Y server.exe "%batPath%/http/https_tls.exe"

cd %batPath%
cd ../nobl/server & go build
copy /Y server.exe "%batPath%/server_both/server_both.exe"
copy /Y server.exe "%batPath%/server_noca/server_noca.exe"

cd %batPath%

