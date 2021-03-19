package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/SkycoinPro/skywire-services-auth/src/api"
	"github.com/SkycoinPro/skywire-services-auth/src/rpc/rpc_client"
	"github.com/SkycoinPro/skywire-services-auth/src/user"

	"golang.org/x/crypto/bcrypt"

	"net"

	"fmt"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Controller is handling reguests regarding Model
type Controller struct {
	userService user.Service
}

func DefaultController() Controller {
	return NewController(user.DefaultService())
}

func NewController(us user.Service) Controller {
	return Controller{
		userService: us,
	}
}

var userDisabledReqParam = "userDisabled"

func (ctrl Controller) RegisterAPIs(public *gin.RouterGroup, closed *gin.RouterGroup) {
	authorization := ctrl.initJWT()
	closed.Use(authorization.MiddlewareFunc())

	closedAuthGroup := closed.Group("/auth")
	closedAuthGroup.GET("/info", ctrl.info)
	closedAuthGroup.GET("/refresh", authorization.RefreshHandler)

	public.POST("/auth/login", authorization.LoginHandler)
}

// ErrUnauthorized is the error returned when user can't be found in JWT
var ErrUnauthorized = errors.New("auth controller: user is not recognized")

func (ctrl *Controller) initJWT() *jwt.GinJWTMiddleware {
	var conf jwtConfig
	viper.UnmarshalKey("jwt", &conf)
	log.Infof("Initializing JWT with params %+v", conf)
	return &jwt.GinJWTMiddleware{
		Realm:            conf.Realm,
		Key:              []byte(conf.Key),
		SigningAlgorithm: conf.Algorithm,
		Timeout:          conf.Timeout,
		MaxRefresh:       conf.MaxRefresh,
		Authenticator: func(email string, password string, c *gin.Context) (string, bool) {
			email = strings.Replace(strings.ToLower(email), " ", "", -1) //to lower case and remove space
			if strings.Contains(email, "@gmail.com") {
				//if gmail remove dots from username part of email
				splinter := strings.Split(email, "@")
				dropDot := strings.Replace(splinter[0], ".", "", -1)
				email = fmt.Sprintf("%s@%s", dropDot, splinter[1])
			}

			usr, err := ctrl.userService.FindByWithUnscoped(email)
			if err != nil {
				return email, false //TODO consider adding error to gin context
			}
			if usr.DeletedAt != nil {
				c.Set(userDisabledReqParam, true)
				return email, false
			}

			passMatchErr := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password+email))
			if usr.Username == email && passMatchErr == nil {
				twoFactorEnabled := viper.GetBool("server.two-factor-enabled")
				if twoFactorEnabled {
					//TODO: check if we save the login attempts with failed 2fa
					if usr.UseOtp {
						token := c.Request.Header.Get("2fa")
						errValidToken := ctrl.userService.ValidateOtpForUser(token, usr.Username)
						if errValidToken != nil {
							c.Set("2faError", true)
							return email, false
						}
					}
				}

				ipAddress := getIPAddress(c.Request)
				userAgent := c.Request.UserAgent()
				ctrl.saveUserAgentInfo(userAgent, ipAddress, usr.ID)
				return email, true
			}

			return email, false
		},
		Authorizator: func(u string, c *gin.Context) bool {
			//TODO implement this case for admins
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {

			if c.Request.URL.Path == "/api/v1/auth/login" {
				code = http.StatusBadRequest
				if c.GetBool("2faError") {
					c.JSON(code, api.ErrorResponse{Error: "auth service: missing 2FA code"})
					return

				}
				if c.GetBool(userDisabledReqParam) {
					c.JSON(code, api.ErrorResponse{Error: "auth service: User has been disabled. Cannot login."})
					return
				}
			}

			c.JSON(code, api.ErrorResponse{Error: "auth service: Username and/or Password do not match any user"})

		},
		PayloadFunc: func(userID string) map[string]interface{} {
			claims := make(map[string]interface{})
			usr, err := ctrl.userService.FindBy(userID)
			if err != nil {
				log.Error("Unable to fetch current user from the DB", err)
				return claims
			}

			if usr.IsAdmin() {
				claims["can_disable"] = usr.CanDisableUser()
				claims["can_create"] = usr.CanCreateAdmin()
			}
			if !usr.IsConfirmed() {
				claims["missing_confirmation"] = true
			}

			if usr.UseOtp {
				claims["use_otp"] = true
			}

			rights := rpc_client.FetchRightsFromRemoteServices(userID)
			for _, right := range rights {
				claims[right.Name] = right.Value
			}

			return claims
		},
		//TODO (security) consider customizign some of parameters
	}
}

// TODO: handle the case where the address is new one, if needed
func (ctrl Controller) saveUserAgentInfo(userAgent, ipAddress string, userId uint) error {
	_, err := ctrl.userService.AddUserAgentInfo(userAgent, ipAddress, userId)
	if err != nil {
		return err
	}
	return nil
}

// TODO: create a shared module and move this?
func getIPAddress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		if len(addresses[0]) == 0 {
			return stripPortFromIp(r.RemoteAddr)
		}
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			}
			return stripPortFromIp(ip)
		}
	}
	return ""
}

func stripPortFromIp(address string) string {
	return strings.Split(address, ":")[0]
}

// @Summary Retrieve signed in User's info
// @Description Information about currently signed in user is collected and returned as response.
// @Tags authorization
// @Accept  json
// @Produce  json
// @Success 200 {object} user.Model
// @Failure 401 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /auth/info [get]
func (ctrl *Controller) info(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userEmail := claims["id"].(string)

	if len(userEmail) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse{Error: ErrUnauthorized.Error()})
		return
	}

	usr, err := ctrl.userService.FindBy(userEmail)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if msg := err.Error(); strings.Contains(msg, user.ErrMissingMandatoryFields.Error()) {
			statusCode = http.StatusUnprocessableEntity
		}

		c.AbortWithStatusJSON(statusCode, api.ErrorResponse{Error: err.Error()})
		return
	}
	usr.Password = ""
	c.JSON(http.StatusOK, usr)
}

type jwtConfig struct {
	Realm      string
	Algorithm  string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration `mapstructure:"max-refresh"`
}
