package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// ListPosts godoc
// @Summary List posts
// @Description List posts
// @Tags Post
// @Accept json
// @Produce json
// @Router /posts [get]
// @Param body query PostListRequest true "page"
// @Success 200 {object} Response{data=PostListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListPosts(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate query
	var query PostListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// load box
	var box Box
	if err = DB.First(&box, query.BoxID).Error; err != nil {
		return
	}

	// construct querySet
	querySet := query.QuerySet(DB).Where("box_id = ?", query.BoxID)
	if user.ID != box.OwnerID {
		querySet = querySet.Where("is_public = ?", true)
	}

	// load posts from database
	var posts []Post
	if err = querySet.Find(&posts).Error; err != nil {
		return
	}

	// construct response
	var response PostListResponse
	if err = copier.CopyWithOption(&response.Posts, &posts, CopyOption); err != nil {
		return
	}

	return Success(c, response)
}

// GetAPost godoc
// @Summary Get a post
// @Description Get a post
// @Tags Post
// @Accept json
// @Produce json
// @Router /post/{id} [get]
// @Param id path string true "id"
// @Success 200 {object} Response{data=PostGetResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func GetAPost(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get post id
	var postID int
	if postID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load post and associated box from database
	var post Post
	if err = DB.Preload("Box").First(&post, postID).Error; err != nil {
		return
	}

	// check if user is authorized to view this post
	if user.ID != post.Box.OwnerID && user.ID != post.PosterID && !post.IsPublic {
		return Forbidden()
	}

	// load channels' content of the post
	var channelsContent []string
	if err = DB.Model(&Channel{}).Where("post_id = ?", post.ID).Pluck("content", &channelsContent).Error; err != nil {
		return err
	}

	// construct response
	var response PostGetResponse
	if err = copier.CopyWithOption(&response, &post, CopyOption); err != nil {
		return
	}
	response.Channels = channelsContent

	return Success(c, response)
}

// CreateAPost godoc
// @Summary Create a post
// @Description Create a post
// @Tags Post
// @Accept json
// @Produce json
// @Param post body PostCreateRequest true "post"
// @Success 200 {object} Response{data=PostCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /post [post]
func CreateAPost(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// parse and validate body
	var body PostCreateRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	// load box
	var box Box
	if err = DB.First(&box, body.BoxID).Error; err != nil {
		return
	}

	// construct post
	var post Post
	if err = copier.CopyWithOption(&post, &body, CopyOption); err != nil {
		return
	}
	post.PosterID = user.ID

	// create the post to database
	if err = DB.Create(&post).Error; err != nil {
		return
	}

	// construct response
	var response PostCommonResponse
	if err = copier.CopyWithOption(&response, &post, CopyOption); err != nil {
		return
	}

	return Created(c, response)
}

// ModifyAPost godoc
// @Summary Modify a post
// @Description Modify a post, owner only
// @Tags Post
// @Accept json
// @Produce json
// @Router /post/{id} [put]
// @Param id path string true "id"
// @Param json body PostModifyRequest true "post"
// @Success 200 {object} Response{data=PostCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ModifyAPost(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}

	// get post id
	var postID int
	if postID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// parse and validate body
	var body PostModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}
	if body.IsEmpty() {
		return BadRequest("empty body")
	}

	// load post from database
	var post Post
	if err = DB.First(&post, postID).Error; err != nil {
		return
	}

	// check if user is authorized to modify this post
	if user.ID != post.PosterID {
		return Forbidden()
	}

	// update post
	if err = copier.CopyWithOption(&post, &body, CopyOption); err != nil {
		return
	}
	if err = DB.Model(&post).Select("Content", "Visibility").Updates(&post).Error; err != nil {
		return
	}

	// construct response
	var response PostCommonResponse
	if err = copier.Copy(&response, &post); err != nil {
		return
	}

	return Success(c, response)
}

// DeleteAPost godoc
// @Summary Delete a post
// @Description Delete a post, owner only
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=PostDeleteResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /post/{id} [delete]
func DeleteAPost(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}

	// get post id
	var postID int
	if postID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load post from database
	var post Post
	if err = DB.First(&post, postID).Error; err != nil {
		return
	}

	// check if user is authorized to delete this post
	if user.ID != post.PosterID {
		return Forbidden()
	}

	// delete post
	if err = DB.Delete(&post).Error; err != nil {
		return
	}

	return Success(c, Map{})
}