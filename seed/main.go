package main

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/zer0day88/tinder/config"
	"github.com/zer0day88/tinder/infra/db"
	"github.com/zer0day88/tinder/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func main() {
	config.Load()
	log := logger.New()

	pg, err := db.InitPostgres()
	if err != nil {
		panic(err)
	}

	length := 10

	for i := 0; i < length; i++ {
		errIn := insertAuth(pg, i)
		if errIn != nil {
			log.Err(err).Send()
		}
	}

}

func insertAuth(db *pgxpool.Pool, counter int) error {
	q := `INSERT INTO auths (id, email, enc_password) VALUES ($1,$2,$3)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Abcdefghi$"+strconv.Itoa(counter)), 8)
	if err != nil {
		log.Err(err).Send()
	}
	id := uuid.NewString()
	_, err = db.Exec(context.Background(), q, id, faker.Email(), hashedPassword)
	if err != nil {
		return fmt.Errorf("unable to insert row: %v", err)
	}

	return nil
}
