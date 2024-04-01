package v1

import (
	"errors"
	"net/http"
	"os"

	"github.com/SaidovZohid/personal-blog-task/api/models"
	"github.com/SaidovZohid/personal-blog-task/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *handlerV1) AuthMiddleWare(ctx *gin.Context) {
	accessToken := ctx.GetHeader(os.Getenv("AUTHORIZATION_HEADER_KEY"))

	if len(accessToken) == 0 {
		err := errors.New("authorization header is not provided")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Error{
			Code:    401,
			Error:   err.Error(),
			Message: "Header not provided",
		})
		return
	}
	payload, err := utils.VerifyToken(h.cfg, accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Error{
			Code:    401,
			Error:   err.Error(),
			Message: "Wrong Token",
		})
		return
	}

	ctx.Set(os.Getenv("AUTHORIZATION_PAYLOAD_KEY"), payload)
	ctx.Next()
}

func (h *handlerV1) GetAuthPayload(ctx *gin.Context) (*utils.Payload, error) {
	i, exist := ctx.Get(os.Getenv("AUTHORIZATION_PAYLOAD_KEY"))
	if !exist {
		return nil, errors.New("not found payload")
	}

	payload, ok := i.(*utils.Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return payload, nil
}
