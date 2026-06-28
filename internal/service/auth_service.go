package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aravindgopan-dev/cli-auth-system/internal/repository"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

const MaxFailedAttempts = 3

// Consumer-driven interface: Service defines the contract it demands from a database
type DataStore interface {
	CreateUser(ctx context.Context, username, passwordHash string) error
	GetUser(ctx context.Context, username string) (*repository.User, error)
	UpdateUserAuth(ctx context.Context, user *repository.User) error
	CreateSession(ctx context.Context, token, username string, expiresAt time.Time) error
	GetSession(ctx context.Context, token string) (*repository.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

// Concrete implementation struct
type AuthService struct {
	Repo            DataStore
	SessionDuration time.Duration
	LockoutDuration time.Duration
}

func NewAuthService(repo DataStore, sessDur, lockDur time.Duration) *AuthService {
	return &AuthService{
		Repo:            repo,
		SessionDuration: sessDur,
		LockoutDuration: lockDur,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.Repo.CreateUser(ctx, username, string(hashedPassword))
}

func (s *AuthService) PreLoginValidate(ctx context.Context, username string) (*repository.User, error) {
	user, err := s.Repo.GetUser(ctx, username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.LockedUntil.Valid && user.LockedUntil.Time.After(time.Now()) {
		return nil, fmt.Errorf("account locked. Try again after %v", user.LockedUntil.Time.Format("15:04:05"))
	}
	return user, nil
}

func (s *AuthService) PasswordLogin(ctx context.Context, user *repository.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		user.FailedAttempts++
		if user.FailedAttempts >= MaxFailedAttempts {
			user.LockedUntil = sql.NullTime{Time: time.Now().Add(s.LockoutDuration), Valid: true}
		}
		_ = s.Repo.UpdateUserAuth(ctx, user)
		return errors.New("invalid username or password")
	}

	user.FailedAttempts = 0
	user.LockedUntil = sql.NullTime{Valid: false}
	user.LastLogin = sql.NullTime{Time: time.Now(), Valid: true}
	_ = s.Repo.UpdateUserAuth(ctx, user)
	return nil
}

func (s *AuthService) VerifyTOTP(user *repository.User, code string) bool {
	return totp.Validate(code, user.TwoFASecret)
}

func (s *AuthService) CreateSession(ctx context.Context, username string) (string, time.Time, error) {
	token := uuid.New().String()
	expiry := time.Now().Add(s.SessionDuration)
	err := s.Repo.CreateSession(ctx, token, username, expiry)
	return token, expiry, err
}

func (s *AuthService) Generate2FASecret(username string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SecureCLI",
		AccountName: username,
	})
	return key.Secret(), key.URL(), err
}

func (s *AuthService) Enable2FA(ctx context.Context, user *repository.User, secret string) error {
	user.TwoFASecret = secret
	user.TwoFAEnabled = true
	return s.Repo.UpdateUserAuth(ctx, user)
}

func (s *AuthService) Disable2FA(ctx context.Context, user *repository.User) error {
	user.TwoFASecret = ""
	user.TwoFAEnabled = false
	return s.Repo.UpdateUserAuth(ctx, user)
}