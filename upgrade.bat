
@echo on

cd %~dp0
%~d0
set batPath=%cd%

go get github.com/bitly/go-simplejson
go get github.com/gin-gonic/gin
go get github.com/go-kit/kit
go get github.com/golang/mock
go get github.com/improbable-eng/grpc-web
go get github.com/jinzhu/gorm
go get github.com/pkg/errors
go get github.com/rs/cors
go get github.com/stretchr/testify
go get go.uber.org/zap

cd %batPath%

