package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	rd "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/helper"
	"github.com/zer0day88/tinder/infra/db"
	"github.com/zer0day88/tinder/internal/app/domain/repository"
	"github.com/zer0day88/tinder/internal/app/model"
	"net/http"
	"testing"
	"time"
)

type AuthTestSuite struct {
	suite.Suite
	ctx         context.Context
	db          *pgxpool.Pool
	pgContainer *postgres.PostgresContainer
	pgConnStr   string
	rdContainer *redis.RedisContainer
	rdConnStr   string
	rdClient    *rd.Client
	authRepo    *repository.AuthRepository
	authService *AuthService
	log         zerolog.Logger
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (suite *AuthTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	config.LoadWithPath("../../../config")

	pgContainer, err := postgres.RunContainer(
		suite.ctx,
		testcontainers.WithImage("docker.io/postgres:15.3-alpine"),
		postgres.WithDatabase("tinder"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("1234"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	suite.NoError(err)

	connStr, err := pgContainer.ConnectionString(suite.ctx, "sslmode=disable")
	suite.NoError(err)

	db, err := initPostgres(connStr)
	suite.NoError(err)

	suite.pgContainer = pgContainer
	suite.pgConnStr = connStr
	suite.db = db

	redisContainer, err := redis.RunContainer(suite.ctx, testcontainers.WithImage("redis:6"))
	suite.NoError(err)
	rdConnStr, err := redisContainer.ConnectionString(suite.ctx)
	suite.NoError(err)

	rdConnOptions, err := rd.ParseURL(rdConnStr)
	suite.NoError(err)

	rdClient := rd.NewClient(rdConnOptions)

	suite.rdContainer = redisContainer
	suite.rdConnStr = connStr
	suite.rdClient = rdClient

	err = suite.rdClient.Ping(suite.ctx).Err()
	suite.NoError(err)

	suite.authRepo = repository.NewAuthRepository(suite.db)
	suite.authService = NewAuthService(suite.log, suite.rdClient, *suite.authRepo)

}

func (suite *AuthTestSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.NoError(err)

	err = suite.rdContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *AuthTestSuite) SetupTest() {
	db.MigratePG(suite.db)
}
func (suite *AuthTestSuite) TearDownTest() {
	//_, err := suite.db.Exec(suite.ctx, "DROP TABLE IF EXISTS auths CASCADE")
	//suite.NoError(err)
	//suite.rdClient.FlushAll(suite.ctx)
}

func (suite *AuthTestSuite) TestSignUp() {

	suite.Run("Signup with invalid email", func() {
		suite.T().Cleanup(func() {
			_, errd := suite.db.Exec(suite.ctx, "DELETE FROM auths")
			suite.NoError(errd)
			suite.rdClient.FlushAll(suite.ctx)
		})

		signup := model.SignUpRequest{
			Email:    "andregmail.com",
			Password: "Abcdefghi$",
		}

		resp := suite.authService.SignUp(suite.ctx, signup)

		suite.Equal(resp.Code, http.StatusBadRequest)

	})

	suite.Run("Signup with bad password", func() {
		suite.T().Cleanup(func() {
			_, errd := suite.db.Exec(suite.ctx, "DELETE FROM auths")
			suite.NoError(errd)
			suite.rdClient.FlushAll(suite.ctx)
		})
		signup := model.SignUpRequest{
			Email:    "andre@gmail.com",
			Password: "Abcdefghi$",
		}

		resp := suite.authService.SignUp(suite.ctx, signup)

		suite.Equal(resp.Code, http.StatusBadRequest)

	})

	suite.Run("Signup with ok password", func() {
		suite.T().Cleanup(func() {
			_, errd := suite.db.Exec(suite.ctx, "DELETE FROM auths")
			suite.NoError(errd)
			suite.rdClient.FlushAll(suite.ctx)
		})
		signup := model.SignUpRequest{
			Email:    "andre@gmail.com",
			Password: "Abcdefghi$1",
		}

		suite.authService.SignUp(suite.ctx, signup)
		auths, err := suite.authRepo.FindOneByEmail(suite.ctx, signup.Email)
		suite.NoError(err)
		suite.NotNil(auths)

	})
}

func (suite *AuthTestSuite) TestLogin() {

	suite.Run("Login with valid data", func() {
		suite.T().Cleanup(func() {
			_, errd := suite.db.Exec(suite.ctx, "DELETE FROM auths")
			suite.NoError(errd)
			suite.rdClient.FlushAll(suite.ctx)
		})

		signup := model.SignUpRequest{
			Email:    "andre@gmail.com",
			Password: "Abcdefghi$1",
		}

		suite.authService.SignUp(suite.ctx, signup)

		auths, err := suite.authRepo.FindOneByEmail(suite.ctx, signup.Email)

		suite.NoError(err)
		suite.NotNil(auths)

		login := model.SigningRequest{
			Email:    "andre@gmail.com",
			Password: "Abcdefghi$1",
		}

		resp, errApi := suite.authService.Login(suite.ctx, login)

		suite.Equal(http.StatusOK, errApi.Code)

		_, err = helper.ValidateToken(*resp.AccessToken.Token, config.Key.AccessTokenPublicKey)
		suite.NoError(err)

	})

	suite.Run("Login with invalid data", func() {
		suite.T().Cleanup(func() {
			_, errd := suite.db.Exec(suite.ctx, "DELETE FROM auths")
			suite.NoError(errd)
			suite.rdClient.FlushAll(suite.ctx)
		})

		signup := model.SignUpRequest{
			Email:    "andre@gmail.com",
			Password: "Abcdefghi$1",
		}

		suite.authService.SignUp(suite.ctx, signup)

		auths, err := suite.authRepo.FindOneByEmail(suite.ctx, signup.Email)

		suite.NoError(err)
		suite.NotNil(auths)

		login := model.SigningRequest{
			Email:    "andre1@gmail.com",
			Password: "Abcdefghi$1",
		}

		_, errApi := suite.authService.Login(suite.ctx, login)

		suite.Equal(http.StatusUnauthorized, errApi.Code)

	})

}

func initPostgres(connStr string) (*pgxpool.Pool, error) {

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping db failed")
	}

	return dbPool, nil
}
