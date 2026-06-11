//go:build windows

package upgrader

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scryinfo/dot/dot"
)

func NewUpgraderListener(cfg *UpgraderListenerConfig, logger *dot.LoggerType) (*UpgraderListener, func(), error) {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Error().AnErr("net listen error", err).Send()
		return nil, nil, err
	}

	// 核心：返回 Windows 专属的阻塞逻辑（监听 Ctrl+C）
	waitFunc := func(server *http.Server) error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		logger.Info().Msg("Windows: recieve the signal, shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}

	cleanup := func() { ln.Close() }

	return &UpgraderListener{
		Listener: ln,
		WaitFunc: waitFunc,
	}, cleanup, nil
}
