// see: https://connectrpc.com/docs/web/getting-started/
// "note:
// Though the Connect protocol supports all types of streaming RPCs,
// web browsers do not support streaming from the client side across the board.
// The fetch API does specify streaming request bodies, but unfortunately,
// browser vendors have not come to an agreement to support streams from the client—see this WHATWG issue on GitHub.
// This means you can use streaming from the browser, but only server streaming."
package rpcdot

import "github.com/scryinfo/dot/dot"

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
