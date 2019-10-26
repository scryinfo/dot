// Scry Info.  All rights reserved.
// license that can be found in the license file.

package utils

//TlsConfig
//case the CaPem Pem Key are not empty, both tls
//case only the Pem is not empty, server tls,  if the ServerNameOverride exist, then set the ServerNameOverride that it make the pam and key
type TlsConfig struct {
	//public key of ca
	CaPem string `json:"caPem"`
	//public key of client or server
	Pem string `json:"pem"`
	//private key of client
	Key string `json:"key"`
	//if the CaPam and Key are empty, set the ServerNameOverride is the name of the pem
	ServerNameOverride string `json:"serverNameOverride"`
}
