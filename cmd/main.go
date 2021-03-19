package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/SkycoinPro/skywire-services-auth/src/app"
	"github.com/SkycoinPro/skywire-services-auth/src/auth"
	"github.com/SkycoinPro/skywire-services-auth/src/config"
	"github.com/SkycoinPro/skywire-services-auth/src/database/postgres"
	"github.com/SkycoinPro/skywire-services-auth/src/rpc/rpc_server"
	"github.com/SkycoinPro/skywire-services-auth/src/template"
	"github.com/SkycoinPro/skywire-services-auth/src/user"
)

// @title Skywire User System API
// @version 1.0
// @description This is a Skywire User System service.

// @host localhost:8080
// @BasePath /api/v1
func main() {
	config.Init("user-config")
	level, err := log.ParseLevel(viper.GetString("server.log-level"))
	if err != nil {
		log.Info("Unable to use configured log level. Using Info instead")
		level = log.InfoLevel
	}
	log.SetLevel(level)
	template.Init()

	tearDown := postgres.Init()
	defer tearDown()

	us := user.DefaultService()
	ts := template.DefaultService()

	if viper.GetBool("template.send-email-after-import") {
		users, err := us.GetUsers()
		if err != nil {
			log.Error("Notifying users after import failed with error", err)
		}

		for _, user := range users {
			email := user.Username
			token, err := us.ForgotPasswordCreateLink(email)
			if err != nil {
				log.Errorf("Creation of password reset link for user %v failed with error %v", email, err)
				continue
			}
			if mailErr := ts.MailUserCreatedAfterImport(email, token); mailErr != nil {
				log.Errorf("Unable to notify user %v about created account after import", email)
				continue
			}
			log.Debugf("Successfully notified user %v about created account after import", email)
		}
	}

	rpc_server.RunRPCServer(us)

	// register all of the controllers here
	log.Debug("Staring controller registration")
	app.NewServer(
		auth.DefaultController(),
		user.NewController(us, ts),
	).Run()
	log.Info("Server finished with the Run() func")
}
