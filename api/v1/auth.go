package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SaidovZohid/personal-blog-task/api/models"
	"github.com/SaidovZohid/personal-blog-task/pkg/email"
	"github.com/SaidovZohid/personal-blog-task/pkg/utils"
	"github.com/SaidovZohid/personal-blog-task/storage"
	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @ID SignUp
// @Router /auth/signup [post]
// @Summary Sign up for blogger and reader
// @Description user type = {blogger, reader}
// @Tags register
// @Accept json
// @Produce json
// @Param data body models.SignUpReq true "Data"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) SignUp(ctx *gin.Context) {
	var req models.SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required fields!",
		})
		return
	}

	takenUser, err := h.storage.User().GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "email already taken",
			Message: "Please, login to your email. if this email belongs to you :)",
		})
		return
	} else if takenUser != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "email already taken",
			Message: "Please, login to your email. if this email belongs to you :)",
		})
		return
	}

	if err := checkIsValid(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, your input seems like not valid :(",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("unable to hash password", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}

	req.Password = hashedPassword

	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		h.logger.Error("unable to generate code", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}
	hashedCode, err := utils.HashPassword(code)
	if err != nil {
		h.logger.Error("unable to hash code", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}

	go func(e string, code string) {
		if err := email.SendEmail(h.cfg, &email.SendEmailRequest{
			To: []string{e},
			Body: map[string]string{
				"code": code,
			},
			Type:    email.VerificationEmail,
			Subject: "Verification link for Sign Up",
		}); err != nil {
			h.logger.Error(fmt.Sprintf("unable to send email user %s", e), zap.Error(err))
		}
	}(req.Email, code)

	code = hashedCode
	userData, err := json.Marshal(storage.RedisData{
		Password: req.Password,
		Code:     hashedCode,
		Email:    req.Email,
		Role:     req.Role,
	})
	if err != nil {
		h.logger.Error("unable to set marshal data", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}

	err = h.inMemory.Set("sign_up_"+req.Email, string(userData), h.cfg.SignUpLinkTokenTime)
	if err != nil {
		h.logger.Error("unable to set user data to redis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Code:    200,
		Message: "Successfully sent code :), please check your inbox!",
	})
}

// @ID VerifyEmail
// @Router /auth/verify [post]
// @Summary verifyEmail
// @Description Verification Email with Code After successfull request from SignUp
// @Tags register
// @Accept json
// @Produce json
// @Param data body models.VerifyEmailRequest true "Data"
// @Success 200 {object} models.UserLoginAndValidateReq
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) VerifyEmail(ctx *gin.Context) {
	var req models.VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required fields!",
		})
		return
	}

	userData, err := h.inMemory.Get("sign_up_" + req.Email)
	if err != nil {
		h.logger.Error("unable to get user data to redis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}

	var user storage.RedisData
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		h.logger.Error("unable to unmarshal user data to redis", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "some minor errors, we are working on it :)",
			Message: "Please, try again later",
		})
		return
	}

	if err := utils.CheckPassword(req.Code, user.Code); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "code is not valid as we sent to you :(",
		})
		return
	}

	splittedName := strings.Split(req.Email, "@")
	userInfo, err := h.storage.User().Create(ctx, &repo.User{
		Name:     splittedName[0],
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   err.Error(),
			Message: "Something went wrong to create :( info!",
		})
		return
	}

	accessToken, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		Email:    req.Email,
		Duration: h.cfg.Jwt.AccessTokenDuration,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   err.Error(),
			Message: "Something went wrong to create token, try again in login!",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.UserLoginAndValidateReq{
		Info: models.UserInfo{
			Id:        userInfo.Id,
			Name:      userInfo.Name,
			Email:     userInfo.Email,
			Role:      userInfo.Role,
			CreatedAt: userInfo.CreatedAt.Format(time.RFC3339),
		},
		AcceccToken: accessToken,
	})
}
