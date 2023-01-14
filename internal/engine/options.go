package engine

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"net"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hxia043/nest/pkg/util/crypto"
)

func setAuthenticator(username, password string) authn.Authenticator {
	if username != "" && password != "" {
		auth := &authn.Basic{
			Username: username,
			Password: password,
		}
		return auth
	}

	//fmt.Println("Info: try default authentication.")
	return nil
}

func readCaCert(caCert string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(caCert)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(caCert))
	if err := crypto.DecryptAES(ciphertext, plaintext, crypto.Key, crypto.Iv); err != nil {
		return nil, err
	}

	return plaintext, nil
}

func createClientTLSConfig(cert []byte, verifyCert bool, addr string) (*tls.Config, error) {
	if !verifyCert {
		// nolint:golint,gosec
		return &tls.Config{InsecureSkipVerify: true}, nil
	}

	roots := x509.NewCertPool()
	if len(cert) != 0 {
		ok := roots.AppendCertsFromPEM([]byte(cert))
		if !ok {
			//@TODO: specific the error type with `exception: CertificateParseError` and print the cacert name
			return nil, errors.New("can't parse certificate")
		}
	}

	var serverAddr string
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		serverAddr = host
	} else {
		serverAddr = addr
	}

	return &tls.Config{
		ServerName:         serverAddr,
		RootCAs:            roots,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		// nolint:golint,gosec
		MaxVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.X25519},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			// nolint:golint,gosec
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
	}, nil
}

func (im *imageManager) createOptions() ([]remote.Option, error) {
	var verifyCaCert bool = true
	caCert, err := readCaCert(im.registryOption.caCert)
	if err != nil {
		verifyCaCert = false
	}

	tlsConfig, err := createClientTLSConfig(caCert, verifyCaCert, im.registryOption.registry)
	if err != nil {
		return nil, err
	}

	trs := &http.Transport{TLSClientConfig: tlsConfig}

	options := make([]remote.Option, 0)
	options = append(options, remote.WithTransport(trs))

	auth := setAuthenticator(im.registryOption.username, im.registryOption.password)
	if auth != nil {
		options = append(options, remote.WithAuth(auth))
	} else {
		options = append(options, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	}

	return options, nil
}
