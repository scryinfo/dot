// Scry Info.  All rights reserved.
// license that can be found in the license file.

package utils

// TlsConfig
// case the CaPem Pem Key are not empty, both tls
// case only the Pem is not empty, server tls,  if the ServerNameOverride exist, then set the ServerNameOverride that it make the pam and key
type TlsConfig struct {
	//public key of ca
	CaPem string `toml:"ca_pem" json:"ca_pem" yaml:"ca_pem" mapstructure:"ca_pem"` //public key of ca
	//public key of client or server
	Pem string `toml:"pem" json:"pem" yaml:"pem" mapstructure:"pem"` //public key of client or server
	//private key of client
	Key string `toml:"key" json:"key" yaml:"key" mapstructure:"key"` //private key of client
	//if the CaPam and Key are empty, set the ServerNameOverride is the name of the pem
	ServerNameOverride string `toml:"server_name_override" json:"server_name_override" yaml:"server_name_override" mapstructure:"server_name_override"` //if the CaPam and Key are empty, set the ServerNameOverride is the name of the pem
}
