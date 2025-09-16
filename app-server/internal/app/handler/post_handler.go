package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/dto"
	"github.com/maulana1k/forum-app/internal/domain/service"
)

type IPostHandler interface {
	CreatePost(c *fiber.Ctx) error
	GetAllPosts(c *fiber.Ctx) error
	GetPostByID(c *fiber.Ctx) error
	UpdatePost(c *fiber.Ctx) error
	GetUserPosts(c *fiber.Ctx) error
	LikePost(c *fiber.Ctx) error
	UnlikePost(c *fiber.Ctx) error
	// BookmarkPost(c *fiber.Ctx) error
	// UnbookmarkPost(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
}

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) IPostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with content, tags, and optional image
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			post	body		dto.CreatePostRequest	true	"Post data"
//	@Success		201		{object}	dto.PostResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v1/posts/ [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req dto.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Cannot parse request body",
		})
	}

	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// GetAllPosts godoc
//
//	@Summary		Get all posts with pagination
//	@Description	Retrieve all posts with pagination support
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"				default(1)
//	@Param			limit	query		int	false	"Number of posts per page"	default(10)
//	@Success		200		{object}	dto.PaginatedPostsResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v1/posts/ [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	posts, err := h.postService.GetAllPosts(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(posts)
}

// GetPostByID godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieve a specific post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	dto.PostResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id} [get]
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid post ID",
		})
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		if err.Error() == "post not found" {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(post)
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update an existing post (only by author)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Post ID"
//	@Param			post	body		dto.UpdatePostRequest	true	"Updated post data"
//	@Success		200		{object}	dto.PostResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string) // From JWT middleware

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid post ID",
		})
	}

	var req dto.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Cannot parse request body",
		})
	}

	post, err := h.postService.UpdatePost(id, userID, &req)
	if err != nil {
		switch err.Error() {
		case "post not found":
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		case "unauthorized to update this post":
			return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		}
	}

	return c.JSON(post)
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete an existing post (only by author)
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Post ID"
//	@Success		204	"Post deleted successfully"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string) // From JWT middleware

	id := c.Params("id")
	if id != "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid post ID",
		})
	}

	err := h.postService.DeletePost(id, userID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		case "unauthorized to delete this post":
			return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetUserPosts godoc
//
//	@Summary		Get posts by user ID
//	@Description	Retrieve all posts created by a specific user
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Param			page	query		int	false	"Page number"				default(1)
//	@Param			limit	query		int	false	"Number of posts per page"	default(10)
//	@Success		200		{object}	dto.PaginatedPostsResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v1/posts/user/{id} [get]
func (h *PostHandler) GetUserPosts(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	posts, err := h.postService.GetPostsByUserID(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(posts)
}

// LikePost godoc
//
//	@Summary		Like a post
//	@Description	Like a post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id}/like [post]
func (h *PostHandler) LikePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string) // From JWT middleware

	postID := c.Params("id")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid post ID",
		})
	}

	err := h.postService.LikePost(postID, userID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		case "post already liked":
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{"message": "Post liked successfully"})
}

// UnlikePost godoc
//
//	@Summary		Unlike a post
//	@Description	Unlike a previously liked post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id}/unlike [delete]
func (h *PostHandler) UnlikePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string) // From JWT middleware

	postID := c.Params("id")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid post ID",
		})
	}

	err := h.postService.UnlikePost(postID, userID)
	if err != nil {
		switch err.Error() {
		case "post not liked":
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error: err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{"message": "Post unliked successfully"})
}

// BookmarkPost godoc
//
//	@Summary		Bookmark a post
//	@Description	Bookmark a post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id}/bookmark [post]
// func (h *PostHandler) BookmarkPost(c *fiber.Ctx) error {
// 	userID := c.Locals("userID").(string) // From JWT middleware

// 	postID := c.Params("id")
// 	if postID != "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error: "Invalid post ID",
// 		})
// 	}

// 	err := h.postService.BookmarkPost(postID, userID)
// 	if err != nil {
// 		switch err.Error() {
// 		case "post not found":
// 			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		case "post already bookmarked":
// 			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		default:
// 			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		}
// 	}

// 	return c.JSON(fiber.Map{"message": "Post bookmarked successfully"})
// }

// UnbookmarkPost godoc
//
//	@Summary		Remove bookmark from a post
//	@Description	Remove bookmark from a previously bookmarked post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id}/unbookmark [delete]
// func (h *PostHandler) UnbookmarkPost(c *fiber.Ctx) error {
// 	userID := c.Locals("userID").(string) // From JWT middleware

// 	postID := c.Params("id")
// 	if postID != "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error: "Invalid post ID",
// 		})
// 	}

// 	err := h.postService.UnbookmarkPost(postID, userID)
// 	if err != nil {
// 		switch err.Error() {
// 		case "post not bookmarked":
// 			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		default:
// 			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		}
// 	}

// 	return c.JSON(fiber.Map{"message": "Post unbookmarked successfully"})
// }

// RepostPost godoc
//
//	@Summary		Repost a post
//	@Description	Repost a post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id}/repost [post]
// func (h *PostHandler) RepostPost(c *fiber.Ctx) error {
// 	userID := c.Locals("userID").(uint) // From JWT middleware

// 	postID, err := strconv.ParseUint(c.Params("id"), 10, 32)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error: "Invalid post ID",
// 		})
// 	}

// 	err = h.postService.RepostPost(uint(postID), userID)
// 	if err != nil {
// 		switch err.Error() {
// 		case "post not found":
// 			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		case "post already reposted":
// 			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		default:
// 			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		}
// 	}

// 	return c.JSON(fiber.Map{"message": "Post reposted successfully"})
// }

// UnrepostPost godoc
//
//	@Summary		Remove repost from a post
//	@Description	Remove repost from a previously reposted post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/posts/{id}/unrepost [delete]
// func (h *PostHandler) UnrepostPost(c *fiber.Ctx) error {
// 	userID := c.Locals("userID").(uint) // From JWT middleware

// 	postID, err := strconv.ParseUint(c.Params("id"), 10, 32)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 			Error: "Invalid post ID",
// 		})
// 	}

// 	err = h.postService.UnrepostPost(uint(postID), userID)
// 	if err != nil {
// 		switch err.Error() {
// 		case "post not reposted":
// 			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		default:
// 			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
// 				Error: err.Error(),
// 			})
// 		}
// 	}

// 	return c.JSON(fiber.Map{"message": "Post unreposted successfully"})
// }
