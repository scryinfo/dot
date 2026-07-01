//go:build !windows

package notstopupgrade

import (
	"os"
	"os/signal"
	"syscall"

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
	logger.Info().Msgf("listen addr: %s", cfg.Addr)
	_ = upg.Ready()

	waitFunc := func() error {
		<-upg.Exit()

		logger.Info().Msg("Linux: recieve the signal, shutting down...")
		return nil
	}

	cleanup := func() { upg.Stop() }

	return &UpgraderListener{
		Listener:     ln,
		WaitUpgrader: waitFunc,
	}, cleanup, nil
}
