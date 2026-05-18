// see: https://connectrpc.com/docs/web/getting-started/
// "note:
// Though the Connect protocol supports all types of streaming RPCs,
// web browsers do not support streaming from the client side across the board.
// The fetch API does specify streaming request bodies, but unfortunately,
// browser vendors have not come to an agreement to support streams from the client—see this WHATWG issue on GitHub.
// This means you can use streaming from the browser, but only server streaming."
package rpcdot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/certificate"
)

type RpcTls string

const (
	RpcTlsNone     RpcTls = "none"
	RpcTlsInsecure RpcTls = "insecure" //InsecureSkipVerify = true
	RpcTlsSecure   RpcTls = "secure"   //InsecureSkipVerify = false, read name from server cert
	RpcTlsBoth     RpcTls = "both"     // client and server cert are verified
)

type TlsConfig struct {
	Mode     RpcTls `json:"mode" yaml:"mode" toml:"mode"`
	Cert     string `json:"cert" yaml:"cert" toml:"cert"`
	Key      string `json:"key" yaml:"key" toml:"key"`
	RootCert string `json:"rootCert" yaml:"rootCert" toml:"rootCert"`
	PeerCert string `json:"peerCert" yaml:"peerCert" toml:"peerCert"`
}

func (p *TlsConfig) FullPath(sconf dot.SConfig) error {
	if p.Cert != "" {
		cert, err := sconf.FullPath(p.Cert)
		if err != nil {
			return err
		}
		p.Cert = cert
	}
	if p.Key != "" {
		key, err := sconf.FullPath(p.Key)
		if err != nil {
			return err
		}
		p.Key = key
	}
	if p.RootCert != "" {
		rootCert, err := sconf.FullPath(p.RootCert)
		if err != nil {
			return err
		}
		p.RootCert = rootCert
	}
	if p.PeerCert != "" {
		peerCert, err := sconf.FullPath(p.PeerCert)
		if err != nil {
			return err
		}
		p.PeerCert = peerCert
	}
	return nil
}

// needs tls
func (p *TlsConfig) NeedsTls() bool {
	return p.Mode != RpcTlsNone || p.Key != ""
}

func (p *TlsConfig) MakeTlsConfig(sconf dot.SConfig, baseCert *certificate.BaseCertificate) (*tls.Config, error) {
	switch p.Mode {
	case RpcTlsNone:
		return nil, nil
	case RpcTlsInsecure:
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	case RpcTlsSecure:
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
		}
		{
			pool := x509.NewCertPool()
			if p.RootCert != "" {
				rootCertFile, err := sconf.FullPath(p.RootCert)
				if err != nil {
					return nil, fmt.Errorf("failed to get root cert path: %w", err)
				}
				rootCert, err := baseCert.LoadCertificate(rootCertFile)
				if err != nil {
					return nil, fmt.Errorf("failed to load root cert: %w", err)
				}
				pool.AddCert(rootCert)
			}
			peerCertFile, err := sconf.FullPath(p.PeerCert)
			if err != nil {
				return nil, fmt.Errorf("failed to get peer cert path: %w", err)
			}
			if peerCertFile == "" {
				return nil, fmt.Errorf("peer cert is required")
			}
			peerCert, err := baseCert.LoadCertificate(peerCertFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load peer cert: %w", err)
			}
			pool.AddCert(peerCert)

			tlsConfig.RootCAs = pool
			tlsConfig.ServerName = baseCert.ServerName(peerCert)
			if tlsConfig.ServerName == "" {
				return nil, fmt.Errorf("cant get server name from peer certificate")
			}
		}
		return tlsConfig, nil
	case RpcTlsBoth:
		return &tls.Config{
			InsecureSkipVerify: false,
		}, nil
	default:
		return nil, nil
	}
}
