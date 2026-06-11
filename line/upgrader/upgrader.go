//go:build !windows

package upgrader

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/scryinfo/dot/dot"
)

func NewUpgraderListener(cfg *UpgraderListenerConfig, logger *dot.LoggerType) (*UpgraderListener, func(), error) {
	if cfg.PidFile == "" {
		cfg.PidFile = "/tmp/tableflip.pid"
	}
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: cfg.PidFile,
	})
	if err != nil {
		logger.Error().AnErr("new tableflip error: ", err).Send()
		return nil, nil, err
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			logger.Info().Msg("Linux: revieve SIGHUP，fire tableflip update...")
			_ = upg.Upgrade()
		}
	}()

	ln, err := upg.Listen("tcp", cfg.Addr)
	if err != nil {
		upg.Stop()
		logger.Error().AnErr("net listen error", err).Send()
		return nil, nil, err
	}
	_ = upg.Ready()

	waitFunc := func(server *http.Server) error {
		<-upg.Exit()

		logger.Info().Msg("Linux: recieve the signal, shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}

	cleanup := func() { upg.Stop() }

	return &UpgraderListener{
		Listener: ln,
		WaitFunc: waitFunc,
	}, cleanup, nil
}
