package tokens

import (
	"github.com/kyma-project/kyma/components/connector-service/internal/apperrors"
	"github.com/kyma-project/kyma/components/connector-service/internal/tokens/tokencache"
)

type Generator func(length int) (string, apperrors.AppError)

type TokenParams interface {
	ToJSON() ([]byte, error)
}

type Service interface {
	Save(params TokenParams) (string, apperrors.AppError)
}

type tokenService struct {
	tokenLength int
	store       tokencache.TokenCache
	generator   Generator
}

func NewTokenService(tokenLength int, store tokencache.TokenCache, generator Generator) Service {
	return &tokenService{
		tokenLength: tokenLength,
		store:       store,
		generator:   generator,
	}
}

func (svc *tokenService) Save(params TokenParams) (string, apperrors.AppError) {
	jsonData, err := params.ToJSON()
	if err != nil {
		return "", apperrors.Internal("Faild to serilize token params to JSON, %s", err.Error())
	}

	token, err := svc.generator(svc.tokenLength)
	if err != nil {
		return "", apperrors.Internal("Failed to generate token, %s", err.Error())
	}

	svc.store.Put(string(jsonData), token)

	return token, nil
}
