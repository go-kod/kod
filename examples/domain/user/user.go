package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/examples/domain/snowflake"
	redisC "github.com/go-kod/kod/examples/infra/redis"
	"github.com/go-kod/kod/ext/client/kredis"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/go-kod/kod/interceptor/kvalidate"
)

type impl struct {
	kod.Implements[Component]
	kod.WithConfig[config]
	snowflake kod.Ref[snowflake.Component]
	redisComp kod.Ref[redisC.Component]

	redis *redis.Client
}

type config struct {
	Redis     kredis.Config
	SecretKey string
}

type claims struct {
	Uid      uint64 `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (ins *impl) Init(ctx context.Context) error {
	ins.redis = ins.redisComp.Get().Client()

	if ins.Config().SecretKey == "" {
		ins.Config().SecretKey = "my-secret"
	}

	return nil
}

func (ins *impl) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		krecovery.Interceptor(),
		ktrace.Interceptor(),
		kmetric.Interceptor(),
		kvalidate.Interceptor(),
	}
}

type RegisterRequest struct {
	UserName string `redis:"username" validate:"required"`
	Password string `redis:"password" validate:"required"`
}

type RegisterResponse struct {
	Uid uint64
}

// Register register a user
func (ins *impl) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	uid, err := ins.snowflake.Get().NextID(ctx, &snowflake.NextIDRequest{})
	if err != nil {
		return nil, err
	}

	err = ins.redis.HSet(ctx, fmt.Sprintf("user:%d", uid.ID), req).Err()
	if err != nil {
		return nil, err
	}

	err = ins.redis.Set(ctx, fmt.Sprintf("user:n2id:%s", req.UserName), uid.ID, 0).Err()
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{Uid: uid.ID}, nil
}

type DeRegisterRequest struct {
	Uid uint64 `validate:"required"`
}

type DeRegisterResponse struct{}

// DeRegister deregister a user
func (ins *impl) DeRegister(ctx context.Context, req *DeRegisterRequest) (*DeRegisterResponse, error) {
	data := new(RegisterRequest)
	err := ins.redis.HGetAll(ctx, fmt.Sprintf("user:%d", req.Uid)).Scan(data)
	if err != nil {
		return nil, err
	}

	if data.UserName == "" {
		return nil, errors.New("user not found")
	}

	err = ins.redis.Del(ctx, fmt.Sprintf("user:%d", req.Uid)).Err()
	if err != nil {
		return nil, err
	}

	err = ins.redis.Del(ctx, fmt.Sprintf("user:n2id:%s", data.UserName)).Err()
	if err != nil {
		return nil, err
	}

	return &DeRegisterResponse{}, nil
}

type LoginRequest struct {
	UserName string `validate:"required"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	Token string
}

// Login login a user with jwt
func (ins *impl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	uid, err := ins.redis.Get(ctx, fmt.Sprintf("user:n2id:%s", req.UserName)).Uint64()
	if err != nil {
		return nil, err
	}

	if uid == 0 {
		return nil, errors.New("user not found")
	}

	claims := &claims{
		Uid:      uid,
		Username: req.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwt.SignedString([]byte(ins.Config().SecretKey))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
	}, nil
}

type AuthRequest struct {
	Token string `validate:"required"`
}

type AuthResponse struct {
	Valid bool
	Uid   uint64
}

// Auth auth a user with jwt
func (ins *impl) Auth(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	claim := new(claims)
	token, err := jwt.ParseWithClaims(req.Token, claim, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(ins.Config().SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*claims); ok {
		return &AuthResponse{Valid: false, Uid: claims.Uid}, nil
	}

	return &AuthResponse{}, nil
}
