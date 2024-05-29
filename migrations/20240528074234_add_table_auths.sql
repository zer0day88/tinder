-- +goose Up
-- +goose StatementBegin
CREATE TABLE auths (
   id uuid PRIMARY KEY,
   email text NOT NULL unique ,
   enc_password text NOT NULL,
   picture text,
   email_confirmed_at timestamptz,
   confirmation_token text,
   confirmation_sent_at timestamptz,
   password_recovery_token text,
   password_recovery_sent_at timestamptz,
   password_recovery_time_limit timestamptz,
   reset_password_attempt int,
   last_sign_in_at timestamptz,
   max_login_failed_attempt int DEFAULT 5,
   login_failed_attempt int,
   activated BOOLEAN DEFAULT false,
   otp_secret text,
   otp_auth_url text,
   otp_enabled BOOLEAN default false,
   otp_verified BOOLEAN default false,
   blocked BOOLEAN not null DEFAULT FALSE,
   created_at timestamptz not null DEFAULT CURRENT_TIMESTAMP,
   updated_at timestamptz,
   deleted_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table auths;
-- +goose StatementEnd
