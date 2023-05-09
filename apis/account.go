package apis

import (
	. "ChatDanBackend/config"
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"time"
)

// Login godoc
// @Summary Login
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/login [post]
// @Param json body LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} Response{data=UserResponse}
// @Failure 401 {object} Response "用户名或密码错误"
// @Failure 500 {object} Response "Internal Server Error"
func Login(c *fiber.Ctx) (err error) {
	// parse and validate body
	var body LoginRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	// get user from database
	var user User
	if err := DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Unauthorized("用户名或密码错误")
		}
		return err
	}

	// check password
	if !CheckPassword(body.Password, user.HashedPassword) {
		return Unauthorized("用户名或密码错误")
	}

	token, err := CreateJwtToken(UserClaims{UserID: user.ID, IsAdmin: user.IsAdmin})
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Path:    "/",
		Domain:  Config.Hostname,
		Expires: time.Now().Add(time.Hour),
	})

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, response)
}

// Register godoc
// @Summary Register
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/register [post]
// @Param json body LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} Response{data=UserResponse}
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
func Register(c *fiber.Ctx) (err error) {
	// parse and validate body
	var body LoginRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	// create user
	user := User{
		Username:       body.Username,
		HashedPassword: MakePassword(body.Password),
	}

	if err = DB.Transaction(func(tx *gorm.DB) error {
		var exists User
		if err := tx.Where("username = ?", user.Username).First(&exists).Error; err == nil {
			return BadRequest("用户已存在")
		}

		return tx.Create(&user).Error
	}); err != nil {
		return
	}

	// create jwt token
	token, err := CreateJwtToken(UserClaims{UserID: user.ID, IsAdmin: user.IsAdmin})
	if err != nil {
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Path:    "/",
		Domain:  Config.Hostname,
		Expires: time.Now().Add(time.Hour),
	})

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, response)
}

// Reset godoc
// @Summary Reset Password
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/reset [post]
// @Param json body ResetRequest true "json"
// @Success 200 {object} Response{data=UserResponse}
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Invalid JWT Token"
// @Failure 500 {object} Response "Internal Server Error"
func Reset(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}
	if err = DB.Take(&user, user.ID).Error; err != nil {
		return err
	}

	// parse and validate body
	var body ResetRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	// check old password
	if !CheckPassword(body.OldPassword, user.HashedPassword) {
		return Unauthorized("原密码错误")
	}

	// update password
	if err = DB.Model(&user).Update("hashed_password", MakePassword(body.NewPassword)).Error; err != nil {
		return err
	}

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, response)
}