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
// @Router /posts [post]
// @Summary Create a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param post body models.CreatePostRequest true "Post"
// @Success 201 {object} models.Post
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) CreatePost(ctx *gin.Context) {
	var req models.CreatePostRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required fields :)",
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

	if userInfo.Role != Blogger {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "reader or other type of user can not create post",
			Message: "Sorry, but you can not create post",
		})
		return
	}

	post, err := h.storage.Post().Create(ctx, &repo.Post{
		Header: req.Header,
		Body:   req.Body,
		UserId: userInfo.UserId,
	})
	if err != nil {
		h.logger.Error("unable to unmarshal create post", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Please, try again later :)",
		})
		return
	}

	user, err := h.storage.User().Get(ctx, userInfo.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "id not found",
			Message: "User does not exist",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Post{
		ID:        post.Id,
		Header:    post.Header,
		Body:      post.Body,
		UserID:    post.UserId,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UserInfo: &models.UserData{
			Name: user.Name,
		},
	})
}

// @Router /posts/{id} [get]
// @Summary Get a post with it's id
// @Description Create a post with it's id
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 201 {object} models.GetPostInfo
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) GetPost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required fields :)",
		})
		return
	}

	res, err := h.storage.Post().Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "id not found",
				Message: "Post does not exist",
			})
			return
		}
		h.logger.Error("unable to get post", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	comments, err := h.storage.Comment().GetAll(ctx, res.Id)
	if err != nil {
		h.logger.Error("unable to get comments", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	replies, err := h.storage.Reply().GetAll(ctx, &repo.GetAllRepliesParams{
		PostId: res.Id,
	})
	if err != nil {
		h.logger.Error("unable to get replies", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, h.getPostInfo(res, comments, replies))
}

// @Security ApiKeyAuth
// @Router /posts/{id} [put]
// @Summary Update post with it's id as param
// @Description Update post with it's id and user_id as param
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param post body models.UpdatePostRequest true "Post"
// @Success 200 {object} models.Post
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) UpdatePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, path not provided",
		})
		return
	}

	var req models.UpdatePostRequest

	err = ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, fill all required fields :)",
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

	post, err := h.storage.Post().Update(ctx, &repo.UpdatePost{
		Id:     id,
		Header: req.Header,
		Body:   req.Body,
		UserId: userInfo.UserId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, models.Error{
				Code:    404,
				Error:   "id not found",
				Message: "Post does not exist",
			})
			return
		}
		h.logger.Error("unable to update post", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	user, err := h.storage.User().Get(ctx, userInfo.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   "id not found",
			Message: "User does not exist",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Post{
		ID:        post.Id,
		Header:    post.Header,
		Body:      post.Body,
		UserID:    post.UserId,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UserInfo: &models.UserData{
			Name: user.Name,
		},
	})
}

// @Security ApiKeyAuth
// @Router /posts/{id} [delete]
// @Summary Delete a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseSuccess
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) DeletePost(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, path not provided",
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

	err = h.storage.Post().Delete(ctx, id, userInfo.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, models.Error{
				Code:    400,
				Error:   "id not found",
				Message: "Post does not exist",
			})
			return
		}
		h.logger.Error("unable to update post", zap.Error(err))
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

// @Router /posts [get]
// @Summary Get posts by giving limit, page and sorting asc or desc.
// @Description Get posts by giving limit, page and sorting asc or desc.
// @Tags post
// @Accept json
// @Produce json
// @Param filter query models.GetAllPostsParams false "Filter"
// @Success 200 {object} models.GetAllPostsResponse
// @Failure 500 {object} models.Error
// @Failure 400 {object} models.Error
func (h *handlerV1) GetAllPosts(ctx *gin.Context) {
	params, err := validateGetAllPostsParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Error:   err.Error(),
			Message: "Please, query params not provided",
		})
		return
	}

	result, err := h.storage.Post().GetAll(ctx, &repo.GetAllParams{
		Limit:      params.Limit,
		Page:       params.Page,
		UserId:     params.UserId,
		SortByDate: params.SortByDate,
	})
	if err != nil {
		h.logger.Error("unable to update post", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Error:   "internal server error",
			Message: "Something went wrong :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, getPostsResponse(result))
}

func getPostsResponse(data *repo.GetAllResult) *models.GetAllPostsResponse {
	response := models.GetAllPostsResponse{
		Posts: make([]*models.Post, 0),
		Count: data.Count,
	}

	for _, post := range data.Result {
		u := models.Post{
			ID:        post.Id,
			UserID:    post.UserId,
			Header:    post.Header,
			Body:      post.Body,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UserInfo: &models.UserData{
				Name: post.UserInfo.Name,
				Id:   post.UserInfo.Id,
			},
		}
		response.Posts = append(response.Posts, &u)
	}

	return &response
}

func (h handlerV1) getPostInfo(post *repo.Post, comments *repo.GetAllCommentsResult, replies *repo.GetAllRepliesResult) models.GetPostInfo {
	res := models.GetPostInfo{
		ID:        post.Id,
		Header:    post.Header,
		Body:      post.Body,
		UserID:    post.UserId,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UserInfo: &models.UserData{
			Name: post.UserInfo.Name,
			Id:   post.UserId,
		},
		Comments: &models.PostComments{
			Comments: make([]*models.CommentWithReplies, 0),
			Count:    comments.Count,
		},
	}

	commentReplies := make(map[int64][]*models.Reply)
	for _, reply := range replies.Replies {
		commentReplies[reply.CommentId] = append(commentReplies[reply.CommentId], &models.Reply{
			ID:        reply.ID,
			Content:   reply.Content,
			CommentId: reply.CommentId,
			PostId:    reply.PostId,
			UserID:    reply.UserId,
			CreatedAt: reply.CreatedAt.Format(time.RFC3339),
			Info: &models.UserData{
				Name: reply.UserInfo.Name,
				Id:   reply.UserInfo.Id,
			},
		})
	}

	for _, val := range comments.Comments {
		commentInfo := &models.CommentWithReplies{
			ID:        val.ID,
			Content:   val.Content,
			UserID:    val.UserID,
			PostID:    val.PostID,
			CreatedAt: val.CreatedAt.Format(time.RFC3339),
			Info: &models.UserData{
				Name: val.UserInfo.Name,
				Id:   val.UserInfo.Id,
			},
			Replies: &models.GetAllRepliesResponse{
				Replies: make([]*models.Reply, 0),
				Count:   0,
			},
		}
		if replySlice, ok := commentReplies[val.ID]; ok {
			commentInfo.Replies.Replies = replySlice
			commentInfo.Replies.Count = int64(len(replySlice))
		}
		res.Comments.Comments = append(res.Comments.Comments, commentInfo)
	}

	return res
}
