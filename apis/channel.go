package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// ListChannels
// @Summary 查询帖子的所有回复 thread
// @Tags Channel Module
// @Produce json
// @Router /channels [get]
// @Param json query ChannelListRequest true "page"
// @Success 200 {object} Response{data=ChannelListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ListChannels(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate query
	var query ChannelListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return
	}

	// load posts
	var post Post
	if err = DB.First(&post, query.PostID).Error; err != nil {
		return
	}

	// load channels from database
	var channels []Channel
	if err = query.QuerySet(DB).Where("post_id = ?", query.PostID).Find(&channels).Error; err != nil {
		return
	}

	// construct response
	var response ChannelListResponse
	if err = copier.CopyWithOption(&response.Channels, &channels, CopyOption); err != nil {
		return
	}
	for i := range response.Channels {
		response.Channels[i].IsOwner = channels[i].OwnerID == user.ID
	}

	return Success(c, response)
}

// GetAChannel
// @Summary 获取一条回复 thread 信息
// @Tags Channel Module
// @Produce json
// @Router /channel/{id} [get]
// @Param id path int true "channel id"
// @Success 200 {object} Response{data=ChannelCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func GetAChannel(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get channel id
	var channelID int
	if channelID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// load channel from database
	var channel Channel
	if err = DB.First(&channel, channelID).Error; err != nil {
		return
	}

	// construct response
	var response ChannelCommonResponse
	if err = copier.CopyWithOption(&response, &channel, CopyOption); err != nil {
		return
	}
	response.IsOwner = channel.OwnerID == user.ID

	return Success(c, response)
}

// CreateAChannel
// @Summary 创建一条回复 thread
// @Tags Channel Module
// @Accept json
// @Produce json
// @Router /channel [post]
// @Param json body ChannelCreateRequest true "channel"
// @Success 200 {object} Response{data=ChannelCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func CreateAChannel(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get and validate request
	var body ChannelCreateRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	// load post
	var post Post
	if err = DB.First(&post, body.PostID).Error; err != nil {
		return
	}

	// create channel
	channel := Channel{
		PostID:  body.PostID,
		OwnerID: user.ID,
		Content: body.Content,
	}
	if err = DB.Create(&channel).Error; err != nil {
		return
	}

	// construct response
	var response ChannelCommonResponse
	if err = copier.CopyWithOption(&response, &channel, CopyOption); err != nil {
		return
	}
	response.IsOwner = true

	return Created(c, response)
}

// ModifyAChannel
// @Summary 修改一条回复 thread
// @Tags Channel Module
// @Accept json
// @Produce json
// @Router /channel/{id} [put]
// @Param id path int true "channel id"
// @Param json body ChannelModifyRequest true "channel"
// @Success 200 {object} Response{data=ChannelCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ModifyAChannel(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get channel id
	var channelID int
	if channelID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// get and validate request
	var body ChannelModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	// load channel from database
	var channel Channel
	if err = DB.First(&channel, channelID).Error; err != nil {
		return
	}

	// check if user is owner
	if channel.OwnerID != user.ID {
		return Forbidden("you are not the owner of this channel")
	}

	// update channel
	if err = DB.Model(&channel).Updates(body).Error; err != nil {
		return
	}

	// construct response
	var response ChannelCommonResponse
	if err = copier.CopyWithOption(&response, &channel, CopyOption); err != nil {
		return
	}
	response.IsOwner = true

	return Success(c, response)
}

// DeleteAChannel
// @Summary 删除一条回复 thread
// @Tags Channel Module
// @Produce json
// @Router /channel/{id} [delete]
// @Param id path int true "channel id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func DeleteAChannel(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return
	}

	// get channel id
	var channelID int
	if channelID, err = c.ParamsInt("id"); err != nil {
		return
	}

	// get and validate request
	var body ChannelModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return
	}

	// load channel from database
	var channel Channel
	if err = DB.First(&channel, channelID).Error; err != nil {
		return
	}

	// check if user is owner
	if channel.OwnerID != user.ID {
		return Forbidden("you are not the owner of this channel")
	}

	// update channel
	if err = DB.Delete(&channel).Error; err != nil {
		return
	}

	return Success(c, EmptyStruct{})
}