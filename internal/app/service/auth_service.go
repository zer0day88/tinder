package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/helper"
	"github.com/zer0day88/tinder/internal/app/domain/entities"
	"github.com/zer0day88/tinder/internal/app/domain/repository"
	"github.com/zer0day88/tinder/internal/app/model"
	"github.com/zer0day88/tinder/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	log      zerolog.Logger
	authRepo repository.AuthRepository
	rdb      *redis.Client
}

func NewAuthService(log zerolog.Logger, rdb *redis.Client, userRepo repository.AuthRepository) *AuthService {
	return &AuthService{
		log:      log,
		rdb:      rdb,
		authRepo: userRepo,
	}
}

func (s *AuthService) SignUp(ctx context.Context, signup model.SignUpRequest) response.ApiJSON {
	errv := helper.ValidateStruct(signup)

	if errv != nil {
		return response.ErrBadRequest.WithMsg(*errv)
	}

	email := helper.SanitizeEmail(signup.Email)

	if ok, errs := helper.ValidatePassword(signup.Password); !ok {
		return response.ErrBadRequest.WithMsg(errs)
	}

	row, err := s.authRepo.FindOneByEmail(ctx, email)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return response.ErrSystem
	}

	if row != nil {
		return response.ErrBadRequest.WithMsg("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signup.Password), 8)
	if err != nil {
		return response.ErrSystem.WithErr(err)
	}

	id := uuid.NewString()

	err = s.authRepo.Insert(ctx, &entities.Auth{ID: id, Email: signup.Email, EncPassword: string(hashedPassword)})

	if err != nil {
		return response.ErrSystem.WithErr(err)
	}
	return response.OKNoError
}

func (s *AuthService) Login(ctx context.Context, signing model.SigningRequest) (*helper.JwtToken, response.ApiJSON) {
	//validate
	errv := helper.ValidateStruct(signing)

	if errv != nil {
		return nil, response.ErrBadRequest.WithMsg(*errv)
	}

	authData, err := s.authRepo.FindOneByEmail(ctx, signing.Email)

	if err != nil {
		return nil, response.ErrUnauthorized.WithErr(err)
	}

	//compare password
	err = bcrypt.CompareHashAndPassword([]byte(authData.EncPassword), []byte(signing.Password))
	if err != nil {

		//TODO: increase login failed attempt
		//if reached limit, in this case 5, block until specific time

		return nil, response.ErrUnauthorized.WithErr(err)
	}

	var (
		accessTokenExpiresIn   = time.Minute * config.Key.AccessTokenExpiredIn
		accessTokenPrivateKey  = config.Key.AccessTokenPrivateKey
		refreshTokenExpiresIn  = time.Minute * config.Key.RefreshTokenExpiredIn
		refreshTokenPrivateKey = config.Key.RefreshTokenPrivateKey
	)

	//create Token
	accessTokenDetail, err := helper.CreateToken(authData.ID, accessTokenExpiresIn, accessTokenPrivateKey)
	if err != nil {
		return nil, response.ErrSystem.WithErr(err)
	}

	refreshTokenDetail, err := helper.CreateToken(authData.ID, refreshTokenExpiresIn, refreshTokenPrivateKey)
	if err != nil {
		return nil, response.ErrSystem.WithErr(err)
	}

	jwtToken := helper.JwtToken{
		AccessToken:  accessTokenDetail,
		RefreshToken: refreshTokenDetail,
	}

	// for none blocking task such as send email, update last login, etc...
	// or passing the task to event streaming using kafka, NATS or pub/sub
	go func() {

		ctxBg := context.Background()

		jwtKeyAccess := fmt.Sprintf("%s:%s", config.Key.JwtRedisKey, accessTokenDetail.TokenUuid)
		jwtKeyRefresh := fmt.Sprintf("%s:%s", config.Key.JwtRedisKey, refreshTokenDetail.TokenUuid)

		//Save token data
		err = s.rdb.Set(ctxBg, jwtKeyAccess, authData.ID, config.Key.AccessTokenExpiredIn*time.Minute).Err()

		if err != nil {
			s.log.Err(err).Send()
		}

		err = s.rdb.Set(ctxBg, jwtKeyRefresh, authData.ID, config.Key.RefreshTokenExpiredIn*time.Minute).Err()
		if err != nil {
			s.log.Err(err).Send()
		}

		//update last login time
		err = s.authRepo.UpdateLastSigning(ctxBg, authData.ID)
		if err != nil {
			s.log.Err(err).Msgf("Failed to update last login, id=%s", authData.ID)
		}
	}()

	return &jwtToken, response.OKNoError
}
