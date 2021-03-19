package user

import (
	"fmt"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"

	"github.com/SkycoinPro/skywire-services-auth/src/api"
	"github.com/SkycoinPro/skywire-services-auth/src/rpc/rpc_client"
	"github.com/SkycoinPro/skywire-services-auth/src/template"

	"bytes"
	"encoding/base64"
	"image/jpeg"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Token = "token"
const Email = "email"
const Id = "id"
const Opt = "otp"

// Controller is handling reguests regarding Model
type Controller struct {
	userService Service
	mailService template.Service
}

func DefaultController() Controller {
	return NewController(DefaultService(), template.DefaultService())
}

func NewController(us Service, ms template.Service) Controller {
	return Controller{
		userService: us,
		mailService: ms,
	}
}

func (ctrl Controller) RegisterAPIs(public *gin.RouterGroup, closed *gin.RouterGroup) {
	publicUserGroup := public.Group("/users")
	closedUserGroup := closed.Group("/users")
	adminGroup := closed.Group("/admin")

	publicUserGroup.POST("", ctrl.create)
	publicUserGroup.GET("/verify", ctrl.verifyRegistration)
	publicUserGroup.GET("/forgotPassword", ctrl.forgotPasswordCreateLink)
	publicUserGroup.POST("/resetPassword", ctrl.resetPassword)
	publicUserGroup.GET("/resendValidationToken", ctrl.resendValidationToken)

	closedUserGroup.PATCH("/password", ctrl.updatePassword)
	closedUserGroup.GET("/setupOTP", ctrl.setupOTP)
	closedUserGroup.GET("/confirmOTP", ctrl.confirmOTP)
	closedUserGroup.POST("/disableOTP", ctrl.disableOTP)

	//admin
	adminGroup.POST("", ctrl.canCreateUsersMiddleware, ctrl.createAdmin)
	adminGroup.GET("/admins", ctrl.canManipulateUsersMiddleware, ctrl.getAllAdmins)
	adminGroup.GET("/users", ctrl.canManipulateUsersMiddleware, ctrl.getAllUsers)
	adminGroup.GET("/users/:username", ctrl.canManipulateUsersMiddleware, ctrl.getByUsername)
	adminGroup.GET("/users/:username/activate", ctrl.canManipulateUsersMiddleware, ctrl.activateUser)
	adminGroup.DELETE("/users/:username", ctrl.canDisableUsersMiddleware, ctrl.deleteByUsername)
	adminGroup.POST("/admins/:username/rights", ctrl.canCreateUsersMiddleware, ctrl.updateAdminRights)

}

func (ctrl Controller) canCreateUsersMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !usr.CanCreateAdmin() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: "Auth: No admin privileges"})
		return
	}
}

func (ctrl Controller) canDisableUsersMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !usr.CanDisableUser() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: "Auth: No admin privileges"})
		return
	}
}

func (ctrl Controller) canManipulateUsersMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !(usr.CanCreateAdmin() || usr.CanDisableUser()) {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: "Auth: No admin privileges"})
		return
	}
}

// @Summary Disable OTP code for your account
// @Description
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Router /users/disableOtp [post]
func (ctrl Controller) disableOTP(c *gin.Context) {
	var tokenReq TokenDisableReq
	email := currentUser(c)
	if err := c.BindJSON(&tokenReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	err := ctrl.userService.disableOtp(email, tokenReq.Password, tokenReq.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Request a setup for OTP code for your account
// @Description Request otp to be set up for your account
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Router /users/setupOTP [get]
func (ctrl Controller) setupOTP(c *gin.Context) {
	email := currentUser(c)

	usr, err := ctrl.userService.FindBy(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	if usr.UseOtp {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: "user service: 2fa already set"})
		return
	}

	key, err := ctrl.userService.GenerateOTPForUser(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	buffer := new(bytes.Buffer)
	img, err := key.Image(200, 200)

	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}
	imgBase64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())

	code := OtpCode{Code: key.Secret(), Image: imgBase64Str}
	c.JSON(http.StatusOK, code)
}

// @Summary Confirm a setup for OTP code on your account
// @Description Send a otp code to the backend, and have it set up on your account, so it will be required for future actions
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Router /users/confirmOTP [get]
func (ctrl Controller) confirmOTP(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Opt]) <= 0 {
		log.Errorf("Incorrect value for otp in URL : %s", c.Request.URL.RequestURI())
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: "Auth: Incorrect value for confirmation token sent"})
		return
	}
	err := ctrl.userService.SetOtp(params[Opt][0], currentUser(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary List all users
// @Description Returns the list of all current users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} user.Model
// @Failure 500 {object} api.ErrorResponse
// @Router /admin/users [get]
func (ctrl Controller) getAllUsers(c *gin.Context) {
	users, err := ctrl.userService.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary List all admins
// @Description Method for admins to get list of all admins
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} user.Model
// @Failure 500 {object} api.ErrorResponse
// @Router /admin/admins [get]
func (ctrl Controller) getAllAdmins(c *gin.Context) {
	users, err := ctrl.userService.getAdmins()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Activate User
// @Description Method for admins to activate User that was deactivated
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /admin/users/:username/activate [get]
func (ctrl Controller) activateUser(c *gin.Context) {
	email := c.Param("username")

	if err := ctrl.userService.ActivateUser(email); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}
}

// @Summary Returns user
// @Description Returns user found by username provided
// @Tags users
// @Accept json
// @Produce json
// @Param username query string true "User's email"
// @Success 200 {object} user.Model
// @Failure 400 {object} api.ErrorResponse
// @Router /admin/users/:username [get]
func (ctrl Controller) getByUsername(c *gin.Context) {
	email := c.Param("username")

	usr, err := ctrl.userService.FindBy(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	rights := rpc_client.FetchRightsFromRemoteServices(email)

	if len(rights) > 0 {
		usr.Rights = append(usr.Rights, rights...)
	}

	c.JSON(http.StatusOK, usr)
}

// @Summary Removes user
// @Description Removes user for given username
// @Tags users
// @Accept json
// @Produce json
// @Param username query string true "Mail of user to be removed"
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Router /admin/users/:username [delete]
func (ctrl Controller) deleteByUsername(c *gin.Context) {
	email := c.Param("username")
	if err := ctrl.userService.RemoveUser(email); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	go ctrl.deleteNodesForUserOnRemoteServices(email)

	ctrl.mailService.MailAccountIsDisabled(email)
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Update admin rights
// @Description Changes rights of a user according to the request
// @Tags users
// @Accept json
// @Produce json
// @Param updated body user.Model true "Model with new rights"
// @Success 200 {object} user.Model
// @Failure 400 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /admin/admins/:username/rights [post]
func (ctrl Controller) updateAdminRights(c *gin.Context) {
	email := c.Param("username")
	var updated Model
	if err := c.BindJSON(&updated); err != nil || updated.Username != email {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}
	dbUser, err := ctrl.userService.FindBy(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	if err := ctrl.userService.UpdateRights(&dbUser, updated.Rights); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: errUnableToSave.Error()})
		return
	}
	go ctrl.UpdateRightsOnRemoteServices(email, updated.Rights)
	dbUser.Password = "" //TODO remove pass always?
	c.JSON(http.StatusOK, dbUser)
}

// @Summary Create a new User in the system
// @Description Collect provided User attributes from the body and create new User in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param newUser body user.Model true "New User"
// @Success 201 {object} user.Model
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users [post]
func (ctrl Controller) create(c *gin.Context) {
	var newUser Model
	if err := c.BindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	actionLink := createNotUsedActionLink()
	newUser.ActionLinks = []ActionLink{actionLink}

	if err := ctrl.userService.Create(&newUser); err != nil {
		statusCode := http.StatusBadRequest
		if msg := err.Error(); strings.Contains(msg, ErrMissingMandatoryFields.Error()) {
			statusCode = http.StatusUnprocessableEntity
		} else if msg := err.Error(); strings.Contains(msg, errUnableToSave.Error()) {
			statusCode = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(statusCode, api.ErrorResponse{Error: err.Error()})
		return
	}
	ctrl.sendEmailForProfileConfirmation(newUser.Username, actionLink.Token)
	c.JSON(http.StatusCreated, newUser)
}

// @Summary Create a new Admin in the system
// @Description Collect provided Admin attributes from the body and create new Admin in the system
// @Tags admins
// @Accept  json
// @Produce  json
// @Param newAdmin body user.Model true "New Admin"
// @Success 201 {object} user.Model
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /admin [post]
func (ctrl Controller) createAdmin(c *gin.Context) {
	var newAdmin Model
	if err := c.BindJSON(&newAdmin); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	actionLink := createNotUsedActionLink()
	newAdmin.ActionLinks = []ActionLink{actionLink}

	isExistingUserUpdated, err := ctrl.userService.CreateAdmin(&newAdmin)
	if err != nil {
		statusCode := http.StatusBadRequest
		if msg := err.Error(); strings.Contains(msg, ErrMissingMandatoryFields.Error()) {
			statusCode = http.StatusUnprocessableEntity
		} else if msg := err.Error(); strings.Contains(msg, errUnableToSave.Error()) {
			statusCode = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(statusCode, api.ErrorResponse{Error: err.Error()})
		return
	}
	if isExistingUserUpdated {
		ctrl.mailService.MailUserStatusChangedToAdmin(newAdmin.Username)
	} else {
		ctrl.sendEmailForProfileConfirmation(newAdmin.Username, actionLink.Token)
	}

	c.JSON(http.StatusCreated, newAdmin)
}

// @Summary Re-send welcome mail containing verification token
// @Description In case User misplaces original welcome mail new one is sent to the registered email address
// @Tags users
// @Accept  json
// @Produce  json
// @Param email query string true "User's email address"
// @Success 200
// @Failure 400 {object} api.ErrorResponse
// @Router /users/resendValidationToken [get]
func (ctrl Controller) resendValidationToken(c *gin.Context) {
	//TODO use one from JWT
	params := c.Request.URL.Query()
	if len(params[Email]) <= 0 {
		log.Errorf("Incorrect value for email in URL : %s", c.Request.URL.RequestURI())
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: "Auth: Incorrect value for user email sent"})
		return
	}
	email := params[Email][0]
	token, err := ctrl.userService.ResendValidationTokenForRegistration(email)
	if err != nil {
		if err != errTechnicalError {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: "Auth: Incorrect value for user email sent"})
			return
		}
		actionLink := createNotUsedActionLink()
		user, err := ctrl.userService.FindBy(email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: "Auth: Incorrect value for user email sent"})
			return
		}
		user.ActionLinks = append(user.ActionLinks, actionLink)
		ctrl.userService.db.updateUser(&user) //TODO rework this

		token = actionLink.Token
	}
	ctrl.sendEmailForProfileConfirmation(email, token)

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Verify User's email address based on token provided via email
// @Description After User creates account first available action is to verify account's email address.
// @Tags users
// @Accept  json
// @Produce  json
// @Param token query string true "User's token for email validation"
// @Success 200 {object} user.Model
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users/verify [get]
func (ctrl Controller) verifyRegistration(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Token]) <= 0 {
		log.Errorf("Incorrect value for token in URL : %s", c.Request.URL.RequestURI())
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: "Auth: Incorrect value for confirmation token sent"})
		return
	}
	user, err := ctrl.userService.Verify(params[Token][0])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Process forgot password request
// @Description If User has forgotten his password and received password reset link this endpoint validates input and stores new password.
// @Tags users
// @Accept  json
// @Produce  json
// @Param forgotPasswordRequest body resetPassReq true "User's forgot password input"
// @Success 200
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users/resetPassword [post]
func (ctrl Controller) resetPassword(c *gin.Context) {
	var changePassReq resetPassReq
	if err := c.BindJSON(&changePassReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}
	_, err := ctrl.userService.ResetPassword(changePassReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Start forgot password flow
// @Description If User has forgotten his password this endpoint enables him to reset password using link sent to provided email address.
// @Tags users
// @Accept  json
// @Produce  json
// @Param email query string true "User's email"
// @Success 200
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users/forgotPassword [get]
func (ctrl Controller) forgotPasswordCreateLink(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Email]) <= 0 {
		log.Errorf("Incorrect value for email in URL : %s", c.Request.URL.RequestURI())
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: "Auth: Incorrect value for user email sent"})
		return
	}
	email := params[Email][0]
	token, err := ctrl.userService.ForgotPasswordCreateLink(email)
	//dont notify user that email provided doesnt exist in database
	if err != nil && err != ErrCannotFindUser {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	//send mail only if no errors occured in previous step
	if err == nil {
		ctrl.mailService.MailForPasswordReset(email, token)
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Update Users's password
// @Description Collect, validate and store User's new Skycoin address.
// @Tags users
// @Accept  json
// @Produce  json
// @Param newAddress body addressUpdateReq true "User's new Skycoin address"
// @Success 200 {object} user.Model
// @Failure 400 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users/password [patch]
func (ctrl Controller) updatePassword(c *gin.Context) {
	var updatePass passwordUpdateReq
	if err := c.BindJSON(&updatePass); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	email := currentUser(c)
	updatedUser, err := ctrl.userService.UpdatePassword(email, updatePass.OldPassword, updatePass.Password)
	if err != nil {
		statusCode := http.StatusBadRequest
		if msg := err.Error(); strings.Contains(msg, ErrMissingMandatoryFields.Error()) {
			statusCode = http.StatusUnprocessableEntity
		} else if msg := err.Error(); strings.Contains(msg, errUnableToSave.Error()) {
			statusCode = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(statusCode, api.ErrorResponse{Error: err.Error()})
		return
	}

	ctrl.mailService.MailForPasswordChange(email)
	c.JSON(http.StatusCreated, updatedUser)
}

func (ctrl Controller) sendEmailForProfileConfirmation(username, token string) {
	var confirmationLink = fmt.Sprintf(
		"%s%s%s", viper.GetString("server.frontend-endpoint"),
		viper.GetString("server.user-confirmation-page"),
		token)
	//TODO: figure out a strategy for failed sent emails
	ctrl.mailService.MailConfirmationOfUserSignUp(username, confirmationLink)
}

// UpdateRightsOnRemoteServices - update rights for whitelist and chb
func (ctrl Controller) UpdateRightsOnRemoteServices(username string, rights []authorization.Right) {
	var whitelistRights []authorization.Right
	var chbRights []authorization.Right

	for _, right := range rights {
		if right.Name == "read_transactions" || right.Name == "manage_chb_users" {
			chbRights = append(chbRights, right)
		} else {
			whitelistRights = append(whitelistRights, right)
		}
	}

	err := updateRightsOnRemoteServices(username, "rpc.whitelist.protocol", "rpc.whitelist.address", whitelistRights)
	if err != nil {
		log.Errorf("Failed to update rights on Whitelist for user %v, due to error: %v", username, err)
	}
	err = updateRightsOnRemoteServices(username, "rpc.chb.protocol", "rpc.chb.address", chbRights)
	if err != nil {
		log.Errorf("Failed to update rights on Coinhour bank for user %v, due to error: %v", username, err)
	}
}

func updateRightsOnRemoteServices(username, protocol, address string, rights []authorization.Right) error {
	client, err := rpc.DialHTTP(viper.GetString(protocol), viper.GetString(address))
	if err != nil {
		log.Error("dialing:", err)
		return err
	}

	args := &authorization.SetRequest{Username: username, Rights: rights}
	var reply authorization.SetResponse
	err = client.Call("Handler.SetUserAuthorization", args, &reply)
	if err != nil {
		log.Error("authorization access rights change error: ", err)
		return err
	} else {
		fmt.Printf("Authorization update for for %v was success", args.Username)
	}

	return nil
}

func (ctrl Controller) deleteNodesForUserOnRemoteServices(username string) error {
	client, err := rpc.DialHTTP(viper.GetString("rpc.whitelist.protocol"), viper.GetString("rpc.whitelist.address"))
	if err != nil {
		log.Error("dialing:", err)
		return err
	}

	args := &authorization.SetRequest{Username: username, Rights: []authorization.Right{}}
	var reply authorization.SetResponse
	err = client.Call("Handler.DeleteNodesForUser", args, &reply)
	if err != nil {
		log.Errorf("Deleting nodes for user: %v faled due to error: %v ", username, err)
		return err
	}

	log.Infof("Deleting nodes for user: %v was success", username)
	return nil
}

// returns username(email address) of current user
func currentUser(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	return claims["id"].(string)
}

type passwordUpdateReq struct {
	OldPassword string
	Password    string
}

type resetPassReq struct {
	Email    string
	Token    string
	Password string
}

type TokenDisableReq struct {
	Token    string
	Password string
}

type OtpCode struct {
	Code  string
	Image string
}
