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
// @Router /comments [post]
// @Summary Create a comment
// @Description Create a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment body models.CreateCommentRequest true "Comment"
// @Success 201 {object} models.Comment
// @Failure 500 {object} models.Error
func (h *handlerV1) CreateComment(ctx *gin.Context) {
	var req models.CreateCommentRequest

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

	comment, err := h.storage.Comment().Create(ctx, &repo.Comment{
		Content: req.Content,
		UserID:  payload.UserId,
		PostID:  req.PostId,
	})
	if err != nil {
		h.logger.Error("unable to create comment", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Please, try again later :)",
		})
		return
	}
	c := models.Comment{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
	}
	ctx.JSON(http.StatusOK, c)
}

// @Security ApiKeyAuth
// @Router /comments/{id} [put]
// @Summary Update comment with it's id as param
// @Description Update comment with it's id as param
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param comment body models.UpdateComment true "Comment"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.Error
func (h *handlerV1) UpdateComment(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, provide path of id",
		})
		return
	}

	var req models.UpdateComment
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

	err = h.storage.Comment().Update(ctx, &repo.UpdateComment{
		ID:      id,
		Content: req.Content,
		UserID:  userInfo.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "id not found",
				Message: "Comment does not exist",
			})
			return
		}
		h.logger.Error("unable to update comment", zap.Error(err))
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
// @Router /comments/{id} [delete]
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.Error
func (h *handlerV1) DeleteComment(ctx *gin.Context) {
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

	err = h.storage.Comment().Delete(ctx, id, userInfo.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "id not found",
				Message: "Comment does not exist",
			})
			return
		}
		h.logger.Error("unable to delete comment", zap.Error(err))
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

// @Router /comments [get]
// @Summary Get comments
// @Description Get comments
// @Tags comment
// @Accept json
// @Produce json
// @Param filter query models.GetAllCommentsParams false "Filter"
// @Success 200 {object} models.GetAllCommentsResponse
// @Failure 500 {object} models.Error
func (h *handlerV1) GetAllComments(ctx *gin.Context) {
	id := ctx.Query("post_id")
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "post id not provided",
			Message: "You did not provide the post_id in query param",
		})
		return
	}

	result, err := h.storage.Comment().GetAll(ctx, postId)
	if err != nil {
		h.logger.Error("unable to get all reply", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "We are working on problem",
		})
		return
	}

	ctx.JSON(http.StatusOK, getCommentsResponse(result))
}

func getCommentsResponse(data *repo.GetAllCommentsResult) *models.GetAllCommentsResponse {
	response := models.GetAllCommentsResponse{
		Comments: make([]*models.Comment, 0),
		Count:    data.Count,
	}

	for _, comment := range data.Comments {
		u := models.Comment{
			ID:        comment.ID,
			UserID:    comment.UserID,
			PostID:    comment.PostID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			Info: &models.UserData{
				Name: comment.UserInfo.Name,
				Id:   comment.UserInfo.Id,
			},
		}
		response.Comments = append(response.Comments, &u)
	}

	return &response
}
