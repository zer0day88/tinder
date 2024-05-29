package entities

import (
	"time"
)

type Auth struct {
	ID          string
	Email       string
	EncPassword string

	//for complex cases
	Picture                   *string
	EmailConfirmedAt          *time.Time
	ConfirmationToken         *string
	ConfirmationSentAt        *time.Time
	PasswordRecoveryToken     *string
	PasswordRecoverySentAt    *time.Time
	PasswordRecoveryTimeLimit *time.Time
	ResetPasswordAttempt      *int
	LastSignInAt              *time.Time
	MaxLoginFailedAttempt     *int
	LoginFailedAttempt        *int
	OtpSecret                 *string
	OtpAuthUrl                *string
	OtpEnabled                bool
	OtpVerified               bool
	Blocked                   bool

	//timestamp
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
