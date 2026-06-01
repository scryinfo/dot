package line

import (
	"github.com/scryinfo/dot/line/certificate"
	contextex "github.com/scryinfo/dot/line/context_ex"
)

var (
	CertificateNewBaseCertificate = certificate.NewBaseCertificate
	CertificateNewEcdsa           = certificate.NewEcdsa
	CertificateNewEd25519         = certificate.NewEd25519
	CertificateNewRsa             = certificate.NewRsa
	CertificateNewSm2             = certificate.NewSm2

	ContextexNewContextEx = contextex.NewContextEx
	// t = bad
)
