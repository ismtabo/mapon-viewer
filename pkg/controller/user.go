package controller

import (
	"net/http"

	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/service"
)

// UserController implements an HTTP Controller for users endpoints.
type UserController interface {
	// LoginUser logs in an user in the application
	LoginUser(rw http.ResponseWriter, r *http.Request)
	// LogoutUser logs out an user session in the application
	LogoutUser(rw http.ResponseWriter, r *http.Request)
	// RegisterUser creates a new user in the application then it logs it in
	RegisterUser(rw http.ResponseWriter, r *http.Request)
}

type userController struct {
	userSvc     service.UserService
	securitySvc service.SessionsService
}

// NewUserController creates a new instance of UserController.
func NewUserController(userSvc service.UserService, securitySvc service.SessionsService) UserController {
	return &userController{userSvc: userSvc, securitySvc: securitySvc}
}

// LoginUser validates user authentication and it starts a session in the incoming HTTP request.
func (c userController) LoginUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		err := errors.NewBadRequestError("malformed POST form body")
		RenderError(ctx, rw, err)
		return
	}
	form := r.Form
	email := form.Get("email")
	if email == "" {
		err := errors.NewBadRequestError("missing email form field")
		RenderError(ctx, rw, err)
		return
	}
	password := form.Get("password")
	if password == "" {
		err := errors.NewBadRequestError("missing password form field")
		RenderError(ctx, rw, err)
		return
	}
	if err := c.userSvc.ValidateUserPassword(ctx, email, password); err != nil {
		RenderError(ctx, rw, err)
		return
	}
	c.securitySvc.Login(rw, r)
}

// LogoutUser finalize creates a session in the incoming HTTP request.
func (c userController) LogoutUser(rw http.ResponseWriter, r *http.Request) {
	c.securitySvc.Logout(rw, r)
}

// Register creates a new user and it starts a session in the incoming HTTP request.
func (c userController) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		err := errors.NewBadRequestError("malformed POST form body")
		RenderError(ctx, rw, err)
		return
	}
	form := r.Form
	email := form.Get("email")
	if email == "" {
		err := errors.NewBadRequestError("missing email form field")
		RenderError(ctx, rw, err)
		return
	}
	password := form.Get("password")
	if password == "" {
		err := errors.NewBadRequestError("missing password form field")
		RenderError(ctx, rw, err)
		return
	}
	if err := c.userSvc.CreateUser(ctx, email, password); err != nil {
		RenderError(ctx, rw, err)
		return
	}
	c.securitySvc.Login(rw, r)
}
