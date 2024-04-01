package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SaidovZohid/personal-blog-task/api/models"
	"github.com/SaidovZohid/personal-blog-task/storage/repo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Security ApiKeyAuth
// @Router /replies [post]
// @Summary Create a reply
// @Description Create a reply to comment
// @Tags reply
// @Accept json
// @Produce json
// @Param reply body models.CreateReplyRequest true "Reply"
// @Success 201 {object} models.Reply
// @Failure 500 {object} models.Error
func (h *handlerV1) CreateReply(ctx *gin.Context) {
	var req models.CreateReplyRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required fields :)",
		})
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    401,
			Error:   err.Error(),
			Message: "Not authorized user",
		})
		return
	}

	reply, err := h.storage.Reply().Create(ctx, &repo.Reply{
		Content:   req.Content,
		UserId:    payload.UserId,
		CommentId: req.CommentId,
		PostId:    req.PostId,
	})
	if err != nil {
		h.logger.Error("unable to create reply", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Please, try again later :)",
		})
		return
	}
	r := models.Reply{
		ID:        reply.ID,
		Content:   reply.Content,
		UserID:    reply.UserId,
		CommentId: reply.CommentId,
		PostId:    reply.PostId,
		CreatedAt: reply.CreatedAt.Format(time.RFC3339),
		Info: &models.UserData{
			Name: payload.Name,
			Id:   payload.UserId,
		},
	}
	ctx.JSON(http.StatusOK, r)
}

// @Security ApiKeyAuth
// @Router /replies/{id} [put]
// @Summary Update reply with it's id as param
// @Description Update reply with it's id as param
// @Tags reply
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param reply body models.UpdateReply true "Reply"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.Error
func (h *handlerV1) UpdateReply(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, provide path of id",
		})
		return
	}

	var req models.UpdateReply
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required :)",
		})
		return
	}

	userInfo, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    401,
			Error:   err.Error(),
			Message: "Not authorized user",
		})
		return
	}

	err = h.storage.Reply().Update(ctx, &repo.UpdateReply{
		Id:      id,
		Content: req.Content,
		UserId:  userInfo.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "id not found",
				Message: "Reply does not exist",
			})
			return
		}
		h.logger.Error("unable to update reply", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Code:    200,
		Message: "Successfully updated comment :)",
	})
}

// @Security ApiKeyAuth
// @Router /replies/{id} [delete]
// @Summary Delete a reply
// @Description Delete a reply
// @Tags reply
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.Error
func (h *handlerV1) DeleteReply(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    401,
			Error:   err.Error(),
			Message: "Not authorized user",
		})
		return
	}

	userInfo, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    401,
			Error:   err.Error(),
			Message: "Not authorized user",
		})
		return
	}

	err = h.storage.Reply().Delete(ctx, id, userInfo.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "id not found",
				Message: "Reply does not exist",
			})
			return
		}
		h.logger.Error("unable to delete reply", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseSuccess{
		Code:    200,
		Message: "Successfully deleted!",
	})
}

// @Router /replies [get]
// @Summary Get replies
// @Description Get replies
// @Tags reply
// @Accept json
// @Produce json
// @Param filter query models.GetAllRepliesParams false "Filter"
// @Success 200 {object} models.GetAllRepliesResponse
// @Failure 500 {object} models.Error
func (h *handlerV1) GetAllReplies(ctx *gin.Context) {
	id := ctx.Query("comment_id")
	commentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "comment_id id not provided",
			Message: "You did not provide the comment_id in query param",
		})
		return
	}
	pId := ctx.Query("post_id")
	var postId int64 = -1
	if pId != "" {
		postId, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "post_id is not integer value",
				Message: "You did not provide the post_id in query param properly",
			})
			return
		}
	}

	result, err := h.storage.Reply().GetAll(ctx, &repo.GetAllRepliesParams{
		CommentId: commentID,
		PostId:    postId,
	})
	if err != nil {
		h.logger.Error("unable to get all reply", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "We are working on problem",
		})
		return
	}

	ctx.JSON(http.StatusOK, getRepliesResponse(result))
}

func getRepliesResponse(data *repo.GetAllRepliesResult) *models.GetAllRepliesResponse {
	response := models.GetAllRepliesResponse{
		Replies: make([]*models.Reply, 0),
		Count:   data.Count,
	}

	for _, reply := range data.Replies {
		u := models.Reply{
			ID:        reply.ID,
			UserID:    reply.UserId,
			CommentId: reply.CommentId,
			Content:   reply.Content,
			PostId:    reply.PostId,
			CreatedAt: reply.CreatedAt.Format(time.RFC3339),
			Info: &models.UserData{
				Name: reply.UserInfo.Name,
				Id:   reply.UserInfo.Id,
			},
		}
		response.Replies = append(response.Replies, &u)
	}

	return &response
}
