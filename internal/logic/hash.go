package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/maxuanquang/ojs/internal/configs"
	"golang.org/x/crypto/bcrypt"
)

type HashLogic interface {
	HashPassword(ctx context.Context, plainPassword string) (string, error)
	IsHashEqual(ctx context.Context, plainPassword string, hashedPassword string) (bool, error)
}

type hashLogic struct {
	authConfig configs.Auth
}

func NewHashLogic(authConfig configs.Auth) HashLogic {
	return &hashLogic{
		authConfig: authConfig,
	}
}

// HashPassword implements Hash.
func (h *hashLogic) HashPassword(ctx context.Context, plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), h.authConfig.Hash.Cost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(bytes), nil
}

// IsHashEqual implements Hash.
func (h *hashLogic) IsHashEqual(ctx context.Context, plainPassword string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("error comparing passwords: %w", err)
	}
	return true, nil
}
