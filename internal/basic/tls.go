package basic

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	// TLSTypeSelfSigned is the self-signed certificate type
	TLSTypeSelfSigned = "self-signed"
)

// NewTLSInfo creates a new instance of the TLSInfo struct
func NewTLSInfo(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) *TLSInfo {
	tls := TLSInfo{
		Enabled:      cfg.Listener.HTTPS.Enabled,
		CertFilename: cfg.Listener.HTTPS.CertFile,
		KeyFilename:  cfg.Listener.HTTPS.KeyFile,
		Router:       router,
		Config:       cfg,
		Logger:       log.With(zap.String("module", "tls")),
	}
	if tls.Enabled {
		tls.Initialize()
	}
	return &tls
}

// TLSInfo is the struct for the TLS subsystem
type TLSInfo struct {
	Enabled      bool
	CertFilename string
	KeyFilename  string
	Router       *gin.Engine
	Config       *config.Configuration
	Logger       *zap.Logger
}

// Initialize initializes the TLS subsystem
func (tls *TLSInfo) Initialize() {
	if tls.Config.Listener.HTTPS.Type == TLSTypeSelfSigned {
		// Create default files if not specified and save to config
		if tls.CertFilename == "" || tls.KeyFilename == "" {
			tls.CertFilename = config.DefaultTLSCertName
			tls.KeyFilename = config.DefaultTLSKeyName
		}

		// Ensure the directories exist and are secure
		if err := os.MkdirAll(filepath.Dir(tls.CertFilename), 0700); err != nil {
			tls.Logger.Fatal("failed to create certificate directory", zap.Error(err))
		}
		if err := os.MkdirAll(filepath.Dir(tls.KeyFilename), 0700); err != nil {
			tls.Logger.Fatal("failed to create certificate key directory", zap.Error(err))
		}

		// Ensure the files exist and create them if they do not
		if _, err := os.Stat(tls.CertFilename); os.IsNotExist(err) {
			if _, err := os.Stat(tls.KeyFilename); os.IsNotExist(err) {
				tls.Logger.Info("generating self-signed certificate",
					zap.String("cert", tls.CertFilename),
					zap.String("key", tls.KeyFilename))

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
func GenerateSelfSignedCertKey(cfg *config.Configuration) error {
	// Generate a ecdsa private key
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// Generate a ecdsa private key for CA
	caPriv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// Create a certificate template(s)
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(365) * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	caTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: "DTAC Agent Certificate Authority"},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:         true,
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"DTAC Agent Certificate Authority"},
		},
		DNSNames:              cfg.Listener.HTTPS.Domains,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caPriv.PublicKey, caPriv)
	if err != nil {
		return err
	}

	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	err = os.WriteFile(cfg.Listener.HTTPS.CAFile, caPEM, 0600)
	if err != nil {
		return err
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	err = os.WriteFile(cfg.Listener.HTTPS.CertFile, certPEM, 0600)
	if err != nil {
		return err
	}

	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	err = os.WriteFile(cfg.Listener.HTTPS.KeyFile, keyPEM, 0600)
	if err != nil {
		return err
	}

	return nil
}
