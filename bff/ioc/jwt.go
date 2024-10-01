package ioc

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	ssov1 "post/api/proto/gen/sso/v1"
)

type Jwt struct {
	ssoSvc    ssov1.AuthServiceClient
	publicKey *ecdsa.PublicKey
}

func NewJWTHandler(sso ssov1.AuthServiceClient) *Jwt {
	return &Jwt{
		ssoSvc: sso,
	}
}

func (j *Jwt) InitJwtValidateToken(ctx context.Context) {
	j.publicKey = j.getPublicKey(ctx)
}

func (j *Jwt) getPublicKey(ctx context.Context) *ecdsa.PublicKey {
	key, err := j.ssoSvc.GetPublicKey(ctx, &ssov1.PublicKeyRequest{})
	if err != nil {
		panic(err)
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(key.GetPublicKey())
	if err != nil {
		panic(err)
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}

	ecdsaPubKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		panic(err)
	}
	return ecdsaPubKey
}
