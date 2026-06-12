//go:build windows

package upgrader

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/scryinfo/dot/dot"
)

func NewUpgraderListener(cfg *UpgraderListenerConfig, logger *dot.LoggerType) (*UpgraderListener, func(), error) {
	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		logger.Error().AnErr("net listen error", err).Send()
		return nil, nil, err
	}
	logger.Info().Msgf("listen addr: %s", cfg.Addr)

	waitFunc := func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		logger.Info().Msg("Windows: recieve the signal, shutting down...")
		return nil
	}

	cleanup := func() { ln.Close() }

	return &UpgraderListener{
		Listener:     ln,
		WaitUpgrader: waitFunc,
	}, cleanup, nil
}
