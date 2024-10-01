package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"post/sso/config"
	"post/sso/domain"
	"testing"
)

type jwtTestSuite struct {
	suite.Suite
	info *config.Info
	svc  domain.AuthService
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(jwtTestSuite))
}

func (j *jwtTestSuite) SetupTest() {
	j.info = InitConfig()
	j.svc = NewJWTService(j.info)
}

func (j *jwtTestSuite) TestJWTService() {
	t := j.T()
	testCases := []struct {
		name    string
		issuer  string
		user    *domain.JwtPayload
		wantErr bool
		want    domain.Claims
	}{
		{
			name:   "合法",
			issuer: "issuer",
			user: &domain.JwtPayload{
				Username: "test",
			},
			wantErr: false,
			want: domain.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer: j.info.Config.Jwt.Issuer,
				},
				JwtPayload: &domain.JwtPayload{
					Username: "test",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			token, err := j.svc.GenerateAccessToken(ctx, tc.user)
			require.NoError(t, err)
			claims, err := j.svc.ValidateToken(ctx, token)
			if !tc.wantErr {
				require.NoError(t, err)
			}
			claims.ExpiresAt = nil
			require.Equal(t, tc.want, *claims)

			token, err = j.svc.GenerateRefreshToken(ctx, tc.user)
			require.NoError(t, err)
			claims, err = j.svc.ValidateToken(ctx, token)
			if !tc.wantErr {
				require.NoError(t, err)
			}
			claims.ExpiresAt = nil
			require.Equal(t, tc.want, *claims)
		})
	}
}

func FindFirstYAMLFile() (string, error) {
	var yamlFile string
	err := filepath.Walk("../config", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			yamlFile = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if yamlFile == "" {
		return "", fmt.Errorf("no yaml file found")
	}

	return yamlFile, nil
}

func InitConfig() *config.Info {
	c := &config.Config{}

	yamlFile, err := FindFirstYAMLFile()
	if err != nil {
		panic(fmt.Errorf("find yaml file error: %v", err))
	}

	yamlConf, err := os.ReadFile(yamlFile)
	if err != nil {
		panic(fmt.Errorf("read yaml error: %v\n", err))
	}

	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		panic(fmt.Errorf("config init unmarshal: %v\n", err))
	}

	return &config.Info{
		Config: c,
	}
}
