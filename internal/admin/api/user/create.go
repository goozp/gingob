package user

import (
	"github.com/puti-projects/puti/internal/admin/api"
	"github.com/puti-projects/puti/internal/admin/service"
	"github.com/puti-projects/puti/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Create user create handler
func Create(c *gin.Context) {
	var r service.UserCreateRequest
	if err := c.Bind(&r); err != nil {
		api.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// check params
	if err := checkParam(&r); err != nil {
		api.SendResponse(c, err, nil)
		return
	}

	svc := service.New(c.Request.Context())
	username, nickname, err := svc.CreateUser(&r)
	if err != nil {
		api.SendResponse(c, err, nil)
	}

	rsp := &service.UserCreateResponse{
		Account:  username,
		Nickname: nickname,
	}

	// Show the user information.
	api.SendResponse(c, nil, rsp)
}

func checkParam(r *service.UserCreateRequest) error {
	if r.Account == "" {
		return errno.New(errno.ErrValidation, nil).Add("account is empty.")
	}

	if r.Password == "" {
		return errno.New(errno.ErrValidation, nil).Add("password is empty.")
	}

	if r.PasswordAgain == "" {
		return errno.New(errno.ErrValidation, nil).Add("check password is empty.")
	}

	if r.Password != r.PasswordAgain {
		return errno.New(errno.ErrValidation, nil).Add("check password is incorrect.")
	}

	if r.Email == "" {
		return errno.New(errno.ErrValidation, nil).Add("Email is empty.")
	}

	if r.Role == "" {
		return errno.New(errno.ErrValidation, nil).Add("role is empty.")
	}

	if r.Role != "administrator" && r.Role != "writer" && r.Role != "subscriber" {
		return errno.New(errno.ErrValidation, nil).Add("role is incorrect.")
	}

	return nil
}
