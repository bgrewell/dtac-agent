package basic

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/bgrewell/dtac-agent/internal/config"
	"go.uber.org/zap"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	// TLSTypeSelfSigned is the self-signed certificate type
	TLSTypeSelfSigned = "self-signed"

	// TLSTypeCASigned is the CA signed certificate type
	TLSTypeCASigned = "ca-signed"
)

// NewTLSInfo creates a new instance of the TLSInfo struct
func NewTLSInfo(cfg *config.Configuration, log *zap.Logger) *map[string]TLSInfo {
	tlsMap := make(map[string]TLSInfo)

	for k, c := range cfg.TLS {
		tls := TLSInfo{
			Enabled:      c.Enabled,
			CertFilename: c.CertFile,
			KeyFilename:  c.KeyFile,
			CAFilename:   c.CAFile,
			Config:       c,
			Logger:       log,
		}
		if tls.Enabled {
			tls.Initialize()
		}
		tlsMap[k] = tls
	}

	return &tlsMap
}

// TLSInfo is the struct for the TLS subsystem
type TLSInfo struct {
	Enabled      bool
	CertFilename string
	KeyFilename  string
	CAFilename   string
	Config       config.TLSConfigurationEntry
	Logger       *zap.Logger
}

// Initialize initializes the TLS subsystem
func (tls *TLSInfo) Initialize() {
	if tls.Config.Type == TLSTypeSelfSigned {
		// Create default files if not specified and save to config
		if tls.CertFilename == "" || tls.KeyFilename == "" || (tls.CAFilename == "" && tls.Config.Type == TLSTypeSelfSigned) {
			tls.CertFilename = config.DefaultTLSCertName
			tls.KeyFilename = config.DefaultTLSKeyName
			tls.CAFilename = config.DefaultTLSCACertName
		}

		// Ensure the directories exist and are secure
		if err := os.MkdirAll(filepath.Dir(tls.CertFilename), 0700); err != nil {
			tls.Logger.Fatal("failed to create certificate directory", zap.Error(err))
		}
		if err := os.MkdirAll(filepath.Dir(tls.KeyFilename), 0700); err != nil {
			tls.Logger.Fatal("failed to create certificate key directory", zap.Error(err))
		}
		if err := os.MkdirAll(filepath.Dir(tls.CAFilename), 0700); err != nil {
			tls.Logger.Fatal("failed to create certificate CA directory", zap.Error(err))
		}

		// Ensure the files exist and create them if they do not
		if _, err := os.Stat(tls.CertFilename); os.IsNotExist(err) {
			if _, err := os.Stat(tls.KeyFilename); os.IsNotExist(err) {
				tls.Logger.Info("generating self-signed certificate",
					zap.String("cert", tls.CertFilename),
					zap.String("key", tls.KeyFilename),
					zap.String("ca", tls.CAFilename))

				if err := GenerateSelfSignedCertKey(tls.Config); err != nil {
					tls.Logger.Fatal("failed to generate self-signed certificate", zap.Error(err))
				}
			} else if err != nil {
				tls.Logger.Fatal("failed to access key file", zap.Error(err), zap.String("key", tls.KeyFilename))
			}
		} else if err != nil {
			tls.Logger.Fatal("failed to access cert file", zap.Error(err), zap.String("cert", tls.KeyFilename))
		}
	}
}

// GenerateSelfSignedCertKey generates a self-signed certificate and key.
// The certificate and key are written to certPath and keyPath respectively.
func GenerateSelfSignedCertKey(cfg config.TLSConfigurationEntry) error {

	// Create CA Certificate Template
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			Organization:  []string{"DTAC Agent Certificate Authority"},
			Country:       []string{"US"},
			Province:      []string{"Oregon"},
			Locality:      []string{"Hillsboro"},
			StreetAddress: []string{"2111 NE 25th Ave"},
			PostalCode:    []string{"97124"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Create Server Certificate Template
	certTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			Organization:  []string{"DTAC Agent Certificate Authority"},
			Country:       []string{"US"},
			Province:      []string{"Oregon"},
			Locality:      []string{"Hillsboro"},
			StreetAddress: []string{"2111 NE 25th Ave"},
			PostalCode:    []string{"97124"},
		},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:              cfg.Domains,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Generate an ecdsa key for CA
	// Note: this key is intentionally not stored so that the CA can not sign any more certificates in the future
	caPriv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// Create the CA
	caDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caPriv.PublicKey, caPriv)
	if err != nil {
		return err
	}

	// PEM encode and write to disk
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	err = os.WriteFile(cfg.CAFile, caPEM, 0600)
	if err != nil {
		return err
	}

	// Generate a ecdsa private key for the server cert
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// Create the server cert with the CA cert as the parent and signed by the CA key
	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &caTemplate, &priv.PublicKey, caPriv)
	if err != nil {
		return err
	}

	// PEM encode cert and write to disk
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	err = os.WriteFile(cfg.CertFile, certPEM, 0600)
	if err != nil {
		return err
	}

	// Convert the private key to DER format
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}

	// PEM encode key and write to disk
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	err = os.WriteFile(cfg.KeyFile, keyPEM, 0600)
	if err != nil {
		return err
	}

	// Try to mangle the CA private key to help prevent in-memory retrieval
	clearValue, err := rand.Int(rand.Reader, priv.Y)
	if err != nil {
		return err
	}
	priv.Y.Set(clearValue)
	priv.X.Set(clearValue)
	priv.D.Set(clearValue)

	return nil
}
