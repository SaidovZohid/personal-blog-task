package v1

import (
	"errors"

	"github.com/SaidovZohid/personal-blog-task/api/models"
	"github.com/SaidovZohid/personal-blog-task/config"
	"github.com/SaidovZohid/personal-blog-task/storage"
	"go.uber.org/zap"
)

var (
	Blogger = "blogger"
	Reader  = "reader"
)

type handlerV1 struct {
	cfg      *config.Config
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
	logger   *zap.Logger
}

type HandlerV1Options struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
	Logger   *zap.Logger
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:      options.Cfg,
		storage:  options.Storage,
		inMemory: options.InMemory,
		logger:   options.Logger,
	}
}

func checkIsValid(req *models.SignUpReq) error {
	// TODO: if you want, add some other validation on password :)

	if req.Role != Reader && req.Role != Blogger {
		return errors.New("role should be either blogger or reader :)")
	}

	return nil
}
