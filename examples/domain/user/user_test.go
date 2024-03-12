package user

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func Test_impl_Register(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, c Component) {
		registerRes, err := c.Register(ctx, &RegisterRequest{UserName: "testuser", Password: "testpassword"})
		assert.Nil(t, err)
		assert.NotZero(t, registerRes.Uid)

		loginRes, err := c.Login(ctx, &LoginRequest{UserName: "testuser", Password: "testpassword"})
		assert.Nil(t, err)
		assert.NotEmpty(t, loginRes.Token)

		authRes, err := c.Auth(ctx, &AuthRequest{Token: loginRes.Token})
		assert.Nil(t, err)
		assert.Equal(t, registerRes.Uid, authRes.Uid)

		deRegisterRes, err := c.DeRegister(ctx, &DeRegisterRequest{Uid: registerRes.Uid})
		assert.Nil(t, err)
		assert.NotNil(t, deRegisterRes)
	})
}
