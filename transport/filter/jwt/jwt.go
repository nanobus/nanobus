package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/security/claims"
	"github.com/nanobus/nanobus/transport/filter"
)

type Config struct {
	RSAPublicKeyFile     string `mapstructure:"rsaPublicKeyFile"`
	RSAPublicKeyString   string `mapstructure:"rsaPublicKeyString"`
	ECDSAPublicKeyFile   string `mapstructure:"ecdsaPublicKeyFile"`
	ECDSAPublicKeyString string `mapstructure:"ecdsaPublicKeyString"`
	HMACSecretKeyFile    string `mapstructure:"hmacSecretKeyFile"`
	HMACSecretKeyBase64  bool   `mapstructure:"hmacSecretKeyBase64"`
	HMACSecretKeyString  string `mapstructure:"hmacSecretKeyString"`
}

type Settings struct {
	RSAPublicKey   *rsa.PublicKey
	ECDSAPublicKey *ecdsa.PublicKey
	HMACSecretKey  []byte
}

// HTTP is the NamedLoader for the assign action.
func HTTP() (string, filter.Loader) {
	return "jwt", HTTPLoader
}

func HTTPLoader(with interface{}, resolver resolve.ResolveAs) (filter.Filter, error) {
	var settings Settings
	var c Config
	err := config.Decode(with, &c)
	if err != nil {
		return nil, err
	}

	var rsaPublicKeyBytes []byte
	if c.RSAPublicKeyFile != "" {
		rsaPublicKeyBytes, err = os.ReadFile(c.RSAPublicKeyFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read public key file: %w", err)
		}
	} else if c.RSAPublicKeyString != "" {
		rsaPublicKeyBytes = []byte(c.RSAPublicKeyString)
	}
	if rsaPublicKeyBytes != nil {
		pubPem, _ := pem.Decode(rsaPublicKeyBytes)
		if pubPem == nil {
			return nil, nil
		}
		var parsedKey interface{}
		if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
			return nil, err
		}

		var ok bool
		if settings.RSAPublicKey, ok = parsedKey.(*rsa.PublicKey); !ok {
			return nil, errors.New("parsed key was not a RSA public key")
		}
	}

	var ecdsaPublicKeyBytes []byte
	if c.ECDSAPublicKeyFile != "" {
		ecdsaPublicKeyBytes, err = os.ReadFile(c.ECDSAPublicKeyFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read public key file: %w", err)
		}
	} else if c.ECDSAPublicKeyString != "" {
		ecdsaPublicKeyBytes = []byte(c.ECDSAPublicKeyString)
	}
	if ecdsaPublicKeyBytes != nil {
		pubPem, _ := pem.Decode(ecdsaPublicKeyBytes)
		if pubPem == nil {
			return nil, nil
		}
		var parsedKey interface{}
		if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
			return nil, err
		}

		var ok bool
		if settings.ECDSAPublicKey, ok = parsedKey.(*ecdsa.PublicKey); !ok {
			return nil, errors.New("parsed key was not a ECDSA public key")
		}
	}

	if c.HMACSecretKeyFile != "" {
		settings.HMACSecretKey, err = os.ReadFile(c.HMACSecretKeyFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read secret key file: %w", err)
		}
		if c.HMACSecretKeyBase64 {
			settings.HMACSecretKey, err = base64.StdEncoding.DecodeString(string(settings.HMACSecretKey))
			if err != nil {
				return nil, fmt.Errorf("cannot base64 decode secret key file: %w", err)
			}
		}
	}

	return HTTPFilter(&settings), nil
}

func HTTPFilter(settings *Settings) filter.Filter {
	return func(ctx context.Context, header filter.Header) (context.Context, error) {
		authorization := header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			return ctx, nil
		}

		token, err := jwt.Parse(authorization[7:], func(token *jwt.Token) (interface{}, error) {
			switch token.Method.(type) {
			case *jwt.SigningMethodRSA:
				return settings.RSAPublicKey, nil
			case *jwt.SigningMethodECDSA:
				return settings.ECDSAPublicKey, nil
			case *jwt.SigningMethodHMAC:
				return settings.HMACSecretKey, nil
			}

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		})
		if err != nil {
			return nil, err
		}

		if token != nil {
			if c, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx = claims.ToContext(ctx, claims.Claims(c))
			}
		}

		return ctx, nil
	}
}
