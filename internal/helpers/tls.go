package helpers

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
	TLS_TYPE_SELF_SIGNED = "self-signed"
)

func NewTlsInfo(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) *TlsInfo {
	tls := TlsInfo{
		Enabled:      cfg.Listener.Https.Enabled,
		CertFilename: cfg.Listener.Https.CertFile,
		KeyFilename:  cfg.Listener.Https.KeyFile,
		Router:       router,
		Config:       cfg,
		Logger:       log.With(zap.String("module", "tls")),
	}
	if tls.Enabled {
		tls.Initialize()
	}
	return &tls
}

type TlsInfo struct {
	Enabled      bool
	CertFilename string
	KeyFilename  string
	Router       *gin.Engine
	Config       *config.Configuration
	Logger       *zap.Logger
}

func (tls *TlsInfo) Initialize() {
	if tls.Config.Listener.Https.Type == TLS_TYPE_SELF_SIGNED {
		// Create default files if not specified and save to config
		if tls.CertFilename == "" || tls.KeyFilename == "" {
			tls.CertFilename = config.DEFAULT_TLS_CERT_NAME
			tls.KeyFilename = config.DEFAULT_TLS_KEY_NAME
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
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(365) * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"DTAC Agent Certificate Authority"},
		},
		DNSNames:              cfg.Listener.Https.Domains,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	err = os.WriteFile(cfg.Listener.Https.CertFile, certPEM, 0600)
	if err != nil {
		return err
	}

	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	err = os.WriteFile(cfg.Listener.Https.KeyFile, keyPEM, 0600)
	if err != nil {
		return err
	}

	return nil
}
