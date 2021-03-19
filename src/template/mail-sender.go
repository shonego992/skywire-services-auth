package template

import (
	"github.com/mattbaird/gochimp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const templateName = "Skycoin template"

// const templateId = "main"

var api *gochimp.MandrillAPI

// Initiate mandril client with the key read from the configuration
func Init() {
	log.Info("Connecting to mandril mail service...")
	apiKey := viper.GetString("template.api-key")
	mandrillApi, err := gochimp.NewMandrill(apiKey)

	if err != nil {
		log.Fatalf("Wrong API key for mandril %v", err)
	}

	_, err = mandrillApi.Ping()
	if err != nil {
		log.Fatalf("Error contacting mandril server %v", err)
	}
	api = mandrillApi
}

func baseSend(a *gochimp.MandrillAPI, content string, subject string, recipient string) error {
	log.Infof("Sending email to user %v with subject %v", recipient, subject)
	log.Debug("Content: ", content)
	if viper.GetBool("template.disable-email-sending") {
		return nil
	}
	//renderedTemplate, err := createTemplate(a, content,template)

	//if err != nil {
	//	log.Errorf("Error rendering template: %v", err)
	//	return err
	//}
	recipients := []gochimp.Recipient{
		gochimp.Recipient{Email: recipient},
	}

	message := gochimp.Message{
		Html:      content,
		Subject:   subject,
		FromEmail: viper.GetString("template.from-email"),
		FromName:  viper.GetString("template.from-name"),
		To:        recipients,
	}

	_, err := a.MessageSend(message, false)

	if err != nil {
		log.Errorf("Error sending message %v", err)
		return err
	}
	return nil
}

// TODO: decide if we want to use template or not - if not, remove the code
//func createTemplate(a *gochimp.MandrillAPI, messageContent,templateId string) (string, error) {
//	contentVar := gochimp.Var{templateId, messageContent}
//	content := []gochimp.Var{contentVar}
//
//	_, err := a.TemplateAdd(templateName, fmt.Sprintf("%s", contentVar.Content), true)
//	if err != nil {
//		log.Errorf("Error adding template: %v", err)
//		return "", err
//	}
//	defer a.TemplateDelete(templateName)
//	renderedTemplate, err := a.TemplateRender(templateName, content, nil)
//	if err != nil {
//		log.Errorf("Error rendering template: %v", err)
//		return "", err
//	}
//	return renderedTemplate, nil
//}
