package rpc_server

import (
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/dchest/uniuri"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authentication"
	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"github.com/SkycoinPro/skywire-services-auth/src/user"
)

var userService user.Service

type Handler int

func (h *Handler) Authenticate(req *authentication.Request, resp *authentication.Response) error {
	log.Debug("Received the call ", req)
	email, password := req.Username, req.Password

	usr, err := userService.FindBy(email)
	if err != nil {
		return err
	}

	passMatchErr := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password+email))
	resp.Success = usr.IsConfirmed() && usr.Username == email && passMatchErr == nil

	return nil
}

func (h *Handler) Create(req *authentication.Request, resp *authentication.Response) error {
	log.Debug("Received the call ", req)
	email := req.Username

	_, err := userService.FindBy(email)
	if err != nil {
		if err != user.ErrCannotFindUser {
			return err
		}
		newUser := user.Model{Username: email, Password: uniuri.NewLen(8)}
		// TODO introduce flag for imported user?
		if err := userService.Create(&newUser); err != nil {
			return err
		}
		resp.Success = true
	} else {
		resp.Success = false
	}

	return nil
}
func (h *Handler) GetCreatedAt(req *authorization.GetRequest, resp *time.Time) error {
	log.Debug("Received the call for created_at time ", req)

	usr, err := userService.FindBy(req.Username)
	if err != nil {
		log.Errorf("Unable to collect created_at for user %v due to error %v ", req.Username, err)
		return err
	}
	*resp = usr.CreatedAt
	return nil
}

// VerifyOtpToken verify token from remote services
func (h *Handler) VerifyOtpToken(req *authentication.Request, resp *authentication.Response) error {
	log.Debug("Received the call for otp token verification ", req)
	username := req.Username

	usr, err := userService.FindBy(username)
	if err != nil {
		log.Errorf("Unable to verify otp token for user %v due to error %v ", username, err)
		return err
	}

	if usr.UseOtp {
		if err := userService.ValidateOtpForUser(req.OtpToken, username); err != nil {
			log.Errorf("Unable to verify otp token for user %v due to error %v ", username, err)
			return err
		}
	}

	if !usr.UseOtp && len(req.OtpToken) != 0 {
		log.Errorf("Otp token is sent but 2fa for user: %v is disabled", username)
		return user.Err2FADisabled

	}

	resp.Success = true
	return nil
}

func RunRPCServer(us user.Service) {
	userService = us
	a := new(Handler)
	rpc.Register(a)
	rpc.HandleHTTP()

	s, err := net.Listen(viper.GetString("rpc.protocol"), viper.GetString("rpc.host"))
	if err != nil {
		log.Fatal("Can't initialize RPC server due to error ", err)
	}
	go http.Serve(s, nil)
	log.Info("Listening for RPC requests on ", viper.GetString("rpc.host"))
}
