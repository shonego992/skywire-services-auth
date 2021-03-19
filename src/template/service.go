package template

import (
	"fmt"
	"strings"

	"github.com/mattbaird/gochimp"
	"github.com/spf13/viper"
)

type Service struct {
	api *gochimp.MandrillAPI
}

// NewService prepares new instance of Service
func NewService(api *gochimp.MandrillAPI) Service {
	return Service{
		api: api,
	}
}

func DefaultService() Service {
	return NewService(api)
}

// Send an email to the user containing the confirmation link for his registration
func (mailService *Service) MailConfirmationOfUserSignUp(receiver string, link string) error {
	content := viper.GetString("template.sign-up-content") + "<a href='" + link + "'>Confirm your account.</a>"
	subject := viper.GetString("template.sign-up-subject")

	return baseSend(mailService.api, content, subject, receiver)
}

// Send an email to user containing link for password reset
func (mailService *Service) MailForPasswordReset(receiver string, token string) error {
	queryParam := fmt.Sprintf(viper.GetString("server.forgot-password-query"), token, receiver)

	url := fmt.Sprintf("%s%s%s", viper.GetString("server.frontend-endpoint"),
		viper.GetString("server.forgot-password-page"), queryParam)

	link := "<a href='" + url + "'>" + url + "</a>"
	content := viper.GetString("template.forgot-password-content") + link
	subject := viper.GetString("template.forgot-password-subject")

	return baseSend(mailService.api, content, subject, receiver)
}

// Send an email to the user notifying him his password is changed
func (mailService *Service) MailForPasswordChange(receiver string) error {
	content := viper.GetString("template.password-changed-content")
	subject := viper.GetString("template.password-changed-subject")

	return baseSend(mailService.api, content, subject, receiver)
}

// Send an email to the user notifying him his profile skycoin address is changed
func (mailService *Service) MailSkycoinAddressChange(receiver string) error {
	content := viper.GetString("template.skycoin-address-changed-content")
	subject := viper.GetString("template.skycoin-address-changed-subject")

	return baseSend(mailService.api, content, subject, receiver)
}

// Send an email to the user notifying him his status is transfered to admin
func (mailService *Service) MailUserStatusChangedToAdmin(receiver string) error {
	content := viper.GetString("template.user-status-changed-to-admin-content")
	subject := viper.GetString("template.user-status-changed-to-admin-subject")

	return baseSend(mailService.api, content, subject, receiver)
}

// Send an email to the user notifying him his account is disabled
func (mailService *Service) MailAccountIsDisabled(receiver string) error {
	content := viper.GetString("template.account-changed-to-disabled-content")
	subject := viper.GetString("template.account-changed-to-disabled-subject")

	return baseSend(mailService.api, content, subject, receiver)
}

// MailUserCreatedAfterImport is sent after import with link for password reset
func (mailService *Service) MailUserCreatedAfterImport(receiver string, token string) error {
	queryParam := fmt.Sprintf(viper.GetString("server.forgot-password-query"), token, receiver)

	url := fmt.Sprintf("%s%s%s", viper.GetString("server.frontend-endpoint"),
		viper.GetString("server.forgot-password-page"), queryParam)

	content := strings.Replace(viper.GetString("template.account-created-after-import-content"), "[[auth-url]]", url, -1)
	subject := viper.GetString("template.account-created-after-import-subject")

	return baseSend(mailService.api, content, subject, receiver)
}
