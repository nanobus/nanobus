package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/nanobus/nanobus/security/claims"
)

func HTTP(ctx context.Context, req *http.Request) (context.Context, error) {
	authorization := req.Header.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		fmt.Println(authorization)
		token, err := jwt.Parse(authorization[7:], func(token *jwt.Token) (interface{}, error) {
			fmt.Println(token.Header["alg"])
			switch token.Method.(type) {
			case *jwt.SigningMethodRSA:
				pub, err := os.ReadFile("public.pem")
				if err != nil {
					return nil, err
				}
				pubPem, _ := pem.Decode(pub)
				if pubPem == nil {
					return nil, nil
				}
				var parsedKey interface{}
				if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
					return nil, err
				}

				var pubKey *rsa.PublicKey
				var ok bool
				if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
					return nil, nil
				}

				return pubKey, nil
			}

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		})
		if err != nil {
			return nil, err
		}

		if token != nil {
			if c, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println(c)
				ctx = claims.ToContext(ctx, claims.Claims(c))
			}
		}
	}

	return ctx, nil
}
