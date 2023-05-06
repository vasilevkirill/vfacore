package vfacore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

//генерируем закрытый ключи и сертфикат, для доменного имени и слхраняем на диск
func generateCertificate() error {
	certFileKey := "./cert.key"
	certFileCer := "./cert.cer"

	if fileExist(certFileKey) {
		err := os.Remove(certFileKey)
		if err != nil {
			return err
		}
	}

	if fileExist(certFileCer) {
		err := os.Remove(certFileCer)
		if err != nil {
			return err
		}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: configGlobalS.Telegram.HookDomain},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}
	certFile, err := os.Create(certFileCer)
	if err != nil {
		return err
	}
	defer closeFileForce(certFile)
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return err
	}
	keyFile, err := os.Create(certFileKey)
	if err != nil {
		return err
	}
	defer closeFileForce(keyFile)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return err
	}
	configGlobalS.Telegram.HookCertKey = certFileKey
	configGlobalS.Telegram.HookCertPub = certFileCer
	return nil
}
