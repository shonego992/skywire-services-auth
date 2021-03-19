package user

import (
	"fmt"
	"net/mail"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"

	"time"

	"github.com/dchest/uniuri"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const ValidationTokenLength = 40
const PasswordMinLength = 8

// Service provides access to User related data
type Service struct {
	db store
}

// DefaultService prepares new instance of Service
func DefaultService() Service {
	return NewService(DefaultData())
}

// NewService prepares new instance of Service
func NewService(userStore store) Service {
	return Service{
		db: userStore,
	}
}

func (us *Service) GetUsers() ([]Model, error) {
	users, err := us.db.getUsers()
	if err != nil {
		return nil, errCannotLoadUsers
	}
	return users, nil
}

func (us *Service) getAdmins() ([]Model, error) {
	users, err := us.db.getAdmins()
	if err != nil {
		return nil, errCannotLoadUsers
	}
	return users, nil
}

//TODO discuss why this route is hit multiple times (after login for example)
func (us *Service) GenerateOTPForUser(email string) (*otp.Key, error) {
	usr, err := us.FindBy(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return nil, ErrCannotFindUser
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Skycoin",
		AccountName: email,
	})
	if err != nil {
		log.Errorf("Unable to generate opt key for user %v due to error %v", email, err)
		return nil, ErrCannotGenerateKey
	}

	checkOtp, err := us.db.findOtpForUser(email)
	if err == nil {
		if checkOtp.Username != email {
			userOtp := Otp{
				Secret:   key.Secret(),
				Username: usr.Username,
			}
			err = us.db.creteOtp(&userOtp)
			if err != nil {
				return key, err
			}
			checkOtp = userOtp
		}
		checkOtp.Secret = key.Secret()
		err = us.db.updateOtp(&checkOtp)
		if err != nil {
			return key, err
		}
	} else {
		return nil, err
	}
	return key, nil
}

func (us *Service) SetOtp(token string, email string) error {
	usr, err := us.FindBy(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return ErrCannotFindUser
	}
	userOtp, err := us.db.findOtpForUser(usr.Username)
	if err != nil {
		log.Errorf("Unable to find 2FA settings by email %v due to error %v", email, err)
		return ErrCannotFindOtp
	}
	valid := totp.Validate(token, userOtp.Secret)
	if !valid {
		log.Error("Error setting up token for user, wrong code")
		return ErrWrongCode
	}
	usr.UseOtp = true
	err = us.db.updateUser(&usr)
	if err != nil {
		return err
	}
	return nil
}

func (us *Service) ValidateOtpForUser(token string, email string) error {
	userOtp, err := us.db.findOtpForUser(email)
	if err != nil {
		log.Errorf("Unable to find 2FA settings by email %v due to error %v", email, err)
		return ErrCannotFindOtp
	}
	valid := totp.Validate(token, userOtp.Secret)
	if !valid {
		return ErrWrongCode
	}
	return nil
}

// Create enables creation of new User instance upon registration
func (us *Service) Create(newUser *Model) error {
	newUser.Username = strings.Replace(strings.ToLower(newUser.Username), " ", "", -1) //to lower case and remove space
	if strings.Contains(newUser.Username, "@gmail.com") {
		//if gmail remove dots from username part of email
		splinter := strings.Split(newUser.Username, "@")
		dropDot := strings.Replace(splinter[0], ".", "", -1)
		newUser.Username = fmt.Sprintf("%s@%s", dropDot, splinter[1])
	}

	if newUser.Username == "" || newUser.Password == "" {
		return ErrMissingMandatoryFields
	} else if len(newUser.Password) < PasswordMinLength {
		return errPasswordTooShort
	} else if _, err := mail.ParseAddress(newUser.Username); err != nil {
		return errEmailNotValid
	} else if u, err := us.db.findBy(newUser.Username, true); err != nil && err != ErrCannotFindUser || u.Username == newUser.Username {
		return errEmailExists //TODO get confirmation for imported users
	}

	// TODO confirm go version we'll use (1.10 or lower for strings.Builder)
	// var builder strings.Builder
	// builder.WriteString(newUser.Password)
	// builder.WriteString(newUser.Username)
	// builder.String()
	bytes, err := us.createPassword(newUser.Username, newUser.Password)
	if err != nil {
		return errUnableToSave //case for internal bcrypt hashing fail, can't test
	}
	newUser.Password = string(bytes)

	if err := us.db.create(newUser); err != nil {
		return errUnableToSave
	}
	newUser.Password = ""
	return nil
}

// CreateAdmin enables creation of new User instance upon registration
func (us *Service) CreateAdmin(newUser *Model) (isExistingUserUpdated bool, err error) {
	if newUser.Username == "" || newUser.Password == "" {
		return false, ErrMissingMandatoryFields
	} else if len(newUser.Password) < PasswordMinLength {
		return false, errPasswordTooShort
	} else if _, err := mail.ParseAddress(newUser.Username); err != nil {
		return false, errEmailNotValid
	} else if u, err := us.db.findBy(newUser.Username, true); err != nil && err != ErrCannotFindUser || u.Username == newUser.Username { //TODO get confirmation for imported users
		if u.IsAdmin() {
			return false, errAdminAlreadyExists
		}
		newUser = &u
	}

	if newUser.ID == 0 {
		// TODO confirm go version we'll use (1.10 or lower for strings.Builder)
		// var builder strings.Builder
		// builder.WriteString(newUser.Password)
		// builder.WriteString(newUser.Username)
		// builder.String()
		bytes, err := us.createPassword(newUser.Username, newUser.Password)
		if err != nil {
			return false, errUnableToSave //case for internal bcrypt hashing fail, can't test
		}
		newUser.Password = string(bytes)
		newUser.SetDisableUser(true)

		if err := us.db.create(newUser); err != nil {
			return false, errUnableToSave
		}
		newUser.Password = ""
		return false, nil
	}

	newUser.SetDisableUser(true)
	if err := us.db.updateUser(newUser); err != nil {
		return true, errUnableToSave
	}

	return true, nil
}

func (us *Service) RemoveUser(email string) error {
	usr, err := us.FindBy(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return ErrCannotFindUser
	}

	if err = us.db.removeUser(&usr); err != nil {
		return errTechnicalError
	}

	return nil
}

// Verify user profile
func (us *Service) Verify(token string) (Model, error) {
	if !verifyTokenLength(token) {
		return Model{}, errIncorrectLengthToken
	}
	link, err := us.db.findLinkByToken(token)
	if err != nil {
		log.Errorf("Cannot find ActionLink by token %v due to error %v", token, err)
		return Model{}, err
	}
	if actionLinkIsExpired(link) {
		return Model{}, errTokenExpired
	}

	if link.Status == Used {
		return Model{}, errAlreadyConfirmed
	}
	u, err := us.db.findUserById(link.UserId)
	if err != nil {
		log.Errorf("Unable to find User by kd %v due to error %v", link.UserId, err)
		return Model{}, err
	}
	u.Confirm()
	link.Status = Used
	for i, element := range u.ActionLinks {
		if element.ID == link.ID {
			u.ActionLinks[i] = link
		}
	}
	err = us.db.updateUser(&u)
	if err != nil {
		return Model{}, err
	}
	u.Password = ""
	return u, nil
}

func (us *Service) ResendValidationTokenForRegistration(username string) (string, error) {
	user, err := us.db.findBy(username, false)
	if err != nil {
		log.Errorf("Unable to find User by id %v due to error %v", username, err)
		return "", err
	}
	for _, element := range user.ActionLinks {
		if element.Type == ConfirmRegistration {
			if element.Status == Used {
				return "", errAlreadyConfirmed
			}
			element.Expiration = time.Now().AddDate(0, 0, viper.GetInt("token.expiration-in-days"))
			err := us.db.updateLink(&element)
			if err != nil {
				return "", errUnableToSave
			}
			return element.Token, nil
		}
	}
	return "", errTechnicalError
}

func (us *Service) ResetPassword(req resetPassReq) (Model, error) {
	username := req.Email
	token := req.Token
	password := req.Password
	if !verifyTokenLength(token) {
		return Model{}, errIncorrectLengthToken
	}

	user, err := us.FindBy(username)
	if err != nil {
		return Model{}, err
	}

	index := -1
	for i, element := range user.ActionLinks {
		if element.Token == token {
			index = i
		}
	}
	if index == -1 {
		log.Errorf("Invalid combination of email and token")
		return Model{}, errWrongTokenAndMailCombination
	}

	link := user.ActionLinks[index]
	if actionLinkIsExpired(link) {
		return Model{}, errTokenExpired
	}

	if link.Status == Used {
		return Model{}, errAlreadyConfirmed
	}

	link.Status = Used
	user.ActionLinks[index] = link
	user, err = us.changeUserPassword(username, password, user)
	if err != nil {
		return Model{}, err
	}
	return user, nil
}

func (us *Service) UpdatePassword(username, oldPassword, newPassword string) (Model, error) {
	if len(username) == 0 || len(oldPassword) < PasswordMinLength || len(newPassword) < PasswordMinLength {
		return Model{}, ErrMissingMandatoryFields
	}

	dbUser, err := us.FindBy(username)
	if err != nil {
		return Model{}, err
	}

	if !dbUser.IsConfirmed() {
		return Model{}, ErrNotConfirmed
	}

	passMatchErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(oldPassword+username))
	if passMatchErr != nil {
		return Model{}, errPasswordDoesNotMatch
	}

	user, err := us.changeUserPassword(username, newPassword, dbUser)
	if err != nil {
		return Model{}, err
	}
	return user, nil
}

func (us *Service) UpdateRights(dbUser *Model, rights []authorization.Right) error {
	dbUser.SetDisableUser(false)
	dbUser.SetCreateAdmin(false)
	for _, right := range rights {
		if right.Name == "disable_user" {
			dbUser.SetDisableUser(right.Value)
		} else if right.Name == "create_user" {
			dbUser.SetCreateAdmin(right.Value)
		}
	}
	us.db.updateUser(dbUser)

	return nil
}

//FindBy retrieves active user by username
func (us *Service) FindBy(username string) (Model, error) {
	if len(username) == 0 {
		return Model{}, ErrMissingMandatoryFields
	}

	u, err := us.db.findBy(username, false)
	if err != nil {
		if err == ErrCannotFindUser {
			log.Debug("There is no record in DB for username ", username)
		} else {
			log.Errorf("Unable to find active User by username %v due to error %v", username, err)
		}

		return Model{}, err
	}
	if u.IsAdmin() {
		u.Rights = append(u.Rights, authorization.Right{Name: "create_user", Label: "Create Admin", Value: u.CanCreateAdmin()})
		u.Rights = append(u.Rights, authorization.Right{Name: "disable_user", Label: "Disable User", Value: u.CanDisableUser()})
	}

	return u, nil
}

//FindByWithUnscoped retrieves active or disabled user by username
func (us *Service) FindByWithUnscoped(username string) (Model, error) {
	if len(username) == 0 {
		return Model{}, ErrMissingMandatoryFields
	}

	u, err := us.db.findBy(username, true)
	if err != nil {
		log.Errorf("Unable to find User by username %v due to error %v", username, err)
		return Model{}, err
	}
	if u.IsAdmin() && u.DeletedAt == nil {
		u.Rights = append(u.Rights, authorization.Right{Name: "create_user", Label: "Create Admin", Value: u.CanCreateAdmin()})
		u.Rights = append(u.Rights, authorization.Right{Name: "disable_user", Label: "Disable User", Value: u.CanDisableUser()})
	}

	return u, nil
}

func (us *Service) ActivateUser(username string) error {
	if len(username) == 0 {
		return ErrMissingMandatoryFields
	}

	return us.db.activate(username)
}

func (us *Service) ForgotPasswordCreateLink(email string) (string, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", errEmailNotValid
	}

	usr, err := us.FindBy(email)
	if err != nil {
		log.Errorf("Unable to find User by email %v due to error %v", email, err)
		return "", ErrCannotFindUser
	}
	link := ActionLink{}
	expiration := time.Now().AddDate(0, 0, viper.GetInt("token.expiration-in-days"))
	token := uniuri.NewLen(40)
	for i, element := range usr.ActionLinks {
		if element.Type == ResetPassword {
			element.Status = NotUsed
			element.Token = token
			element.Expiration = expiration
			usr.ActionLinks[i] = element
			goto update
		}
	}
	link = ActionLink{
		Token:      token,
		Status:     NotUsed,
		Type:       ResetPassword,
		Expiration: expiration,
	}
	usr.ActionLinks = append(usr.ActionLinks, link)
update:
	err = us.db.updateUser(&usr)
	if err != nil {
		log.Errorf("Error while updating user")
		return "", errUnableToSave
	}

	return token, nil
}

// check and return if the ip/agent combination is new. If so, return true.
// return error in case of error while saving in db
func (us *Service) AddUserAgentInfo(client, ipAddress string, userId uint) (bool, error) {
	agents, err := us.db.findUserAgentsByUserId(userId)
	if err != nil {
		log.Errorf("Cannot find UserAgents by user ID %v", userId)
	}
	if agents == nil {
		return false, errUnableToRead
	}
	for _, agent := range agents {
		if agent.Address == ipAddress {
			if agent.Client == client {
				agent.UpdatedAt = time.Now()
				// TODO: check how do we handle the situation where we cannot persist user agent info
				err := us.db.updateUserAgent(&agent)
				return false, err
			}
		}
	}
	info := AgentInfo{
		Client:  client,
		Address: ipAddress,
		UserId:  userId,
	}
	err = us.db.createUserAgent(info)
	if err != nil {
		log.Errorf("Error while creating new user agent")
		return true, errCannotSaveUserAgent
	}
	return true, nil
}

func (us *Service) createPassword(username, password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+username), 14)
}

func verifyTokenLength(token string) bool {
	if len(token) != ValidationTokenLength {
		log.Errorf("Invalid token length for user action token %s", token)
		return false
	}
	return true
}

func (us *Service) changeUserPassword(username, newPassword string, dbUser Model) (Model, error) {
	newPass, err := us.createPassword(username, newPassword)
	if err != nil {
		return Model{}, errUnableToSave
	}
	dbUser.Password = string(newPass)
	us.db.updateUser(&dbUser)
	dbUser.Password = ""

	return dbUser, nil
}

func (us *Service) disableOtp(mail string, password string, token string) error {
	usr, err := us.db.findBy(mail, false)
	if err != nil {
		return err
	}
	passMatchErr := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password+mail))
	if passMatchErr == nil {
		userOtp, err := us.db.findOtpForUser(mail)
		if err != nil {
			log.Errorf("Unable to find 2FA settings by email %v due to error %v", mail, err)
			return ErrCannotFindOtp
		}
		valid := totp.Validate(token, userOtp.Secret)
		if !valid {
			return ErrWrongCode
		}

		usr.UseOtp = false
		err = us.db.updateUser(&usr)
		if err != nil {
			return ErrCannotDisable2FA
		}
		return nil
	}
	return ErrCannotDisable2FA
}

// create a new action link that expires in amount of time read from configuration
func createNotUsedActionLink() ActionLink {
	token := uniuri.NewLen(40)
	return ActionLink{
		Status:     NotUsed,
		Token:      token,
		Type:       ConfirmRegistration,
		Expiration: time.Now().AddDate(0, 0, viper.GetInt("token.expiration-in-days")),
	}
}

func actionLinkIsExpired(link ActionLink) bool {
	if link.Expiration.Before(time.Now()) {
		log.Errorf("Token is expired")
		return true
	}
	return false
}

func containsAccessRight(rights []string, right string) bool {
	for _, r := range rights {
		if r == right {
			return true
		}
	}
	return false
}
