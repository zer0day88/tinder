package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/helper"
	"github.com/zer0day88/tinder/pkg/response"
	"strings"
)

type Header struct {
	Authorization string `header:"Authorization"`
}

type mAuth struct {
	rdb *redis.Client
}

func InitMAuth(rdb *redis.Client) *mAuth {
	return &mAuth{rdb: rdb}
}

func (m *mAuth) JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var accessToken string
		var header Header
		err := (&echo.DefaultBinder{}).BindHeaders(c, &header)

		if err != nil {
			return response.WriteJSON(c, response.ErrUnauthorized.WithMsg("error binding header"))
		}

		authorization := header.Authorization

		if !strings.HasPrefix(authorization, "Bearer ") {
			return response.WriteJSON(c, response.ErrUnauthorized)
		}

		accessToken = strings.TrimPrefix(authorization, "Bearer ")

		tokenClaims, err := helper.ValidateToken(accessToken, config.Key.AccessTokenPublicKey)
		if err != nil {
			return response.WriteJSON(c, response.ErrUnauthorized)
		}

		ctx := c.Request().Context()

		jwtKey := fmt.Sprintf("%s:%s", config.Key.JwtRedisKey, tokenClaims.TokenUuid)
		id, err := m.rdb.Get(ctx, jwtKey).Result()
		if errors.Is(err, redis.Nil) {
			return response.WriteJSON(c, response.ErrUnauthorized)
		}

		c.Set("id", id)

		return next(c)
	}
}
