package v1

import (
	"errors"
	"strconv"

	"github.com/SaidovZohid/personal-blog-task/api/models"
	"github.com/SaidovZohid/personal-blog-task/config"
	"github.com/SaidovZohid/personal-blog-task/storage"
	"github.com/gin-gonic/gin"
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

func validateGetAllPostsParams(ctx *gin.Context) (*models.GetAllPostsParams, error) {
	var (
		limit      int64 = 10
		page       int64 = 1
		userId     int64 = -1
		err        error
		sortByDate string
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("user_id") != "" {
		userId, err = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("sort") != "" &&
		(ctx.Query("sort") == "desc" || ctx.Query("sort") == "asc") {
		sortByDate = ctx.Query("sort")
	}

	return &models.GetAllPostsParams{
		Limit:      limit,
		Page:       page,
		SortByDate: sortByDate,
		UserId:     userId,
	}, nil
}
