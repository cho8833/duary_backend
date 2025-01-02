package test

import (
	"github.com/cho8833/duary_lambda/internal/auth/jwtutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_NewToken(t *testing.T) {
	key := "4fM090neAzdhL/+uR/7pJmO92wWo="
	os.Setenv("secretKey", key)
	jwtUtil := jwtutil.Impl{}

	applicationJwt := jwtUtil.NewToken("1-kakao", key)
	id, err := jwtUtil.ValidateApplicationJWT(applicationJwt.AccessToken, key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, "1-kakao", *id)
	id, err = jwtUtil.ValidateApplicationJWT(applicationJwt.RefreshToken, key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, "1-kakao", *id)
}

func Test_ValidateApplicationJWT(t *testing.T) {
	key := "4fM090neAzdhL/+uR/7pJmO92wWo="
	os.Setenv("secretKey", key)
	jwtUtil := jwtutil.Impl{}

	applicationJwt := jwtUtil.NewToken("1-kakao", key)

	id, err := jwtUtil.ValidateApplicationJWT(applicationJwt.AccessToken, key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, "1kakao", *id)
}
