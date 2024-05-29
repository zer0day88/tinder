package route

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/zer0day88/tinder/internal/app/domain/repository"
	"github.com/zer0day88/tinder/internal/app/handler"
	"github.com/zer0day88/tinder/internal/app/middleware"
	"github.com/zer0day88/tinder/internal/app/service"
)

func InitRoute(e *echo.Echo,
	db *pgxpool.Pool, log zerolog.Logger, rdb *redis.Client) {

	authRepo := repository.NewAuthRepository(db)
	userSrv := service.NewAuthService(log, rdb, *authRepo)
	authHandler := handler.NewAuthHandler(log, userSrv)
	mAuth := middleware.InitMAuth(rdb)

	r := e.Group("/v1")
	{
		r.POST("/signup", authHandler.SignUp)
		r.POST("/login", authHandler.Login)
		cek := r.Group("/cek")
		cek.Use(mAuth.JWT)
		{
			cek.GET("", authHandler.Cek)
		}

	}

}
