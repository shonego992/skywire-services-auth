package user

import (
	"github.com/SkycoinPro/skywire-services-auth/src/database/postgres"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// store is user related interface for dealing with database operations
type store interface {
	create(newUser *Model) error
	creteOtp(otp *Otp) error
	removeUser(user *Model) error
	activate(email string) error
	findBy(email string, includeDisabled bool) (Model, error)
	findOtpForUser(email string) (Otp, error)
	findUserById(id uint) (Model, error)
	findLinkByToken(token string) (ActionLink, error)
	updateUser(user *Model) error
	updateOtp(otp *Otp) error
	updateLink(link *ActionLink) error
	createUserAgent(info AgentInfo) error
	findUserAgentsByUserId(id uint) ([]AgentInfo, error)
	updateUserAgent(agent *AgentInfo) error
	getUsers() ([]Model, error)
	getAdmins() ([]Model, error)
}

// data implements store interface which uses GORM library
type data struct {
	db *gorm.DB
}

func DefaultData() data {
	return NewData(postgres.DB)
}

func NewData(database *gorm.DB) data {
	return data{
		db: database,
	}
}

const adminStatusStart uint8 = 16

func (u data) getUsers() ([]Model, error) {
	var users []Model
	record := u.db.Unscoped().Where("status < ?", adminStatusStart).Find(&users)
	if record.RecordNotFound() {
		return nil, errCannotLoadUsers
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while qu", err)
		}
		return nil, errUnableToRead
	}
	return users, nil
}

func (u data) getAdmins() ([]Model, error) {
	var users []Model
	record := u.db.Unscoped().Where("status >= ?", adminStatusStart).Find(&users)
	if record.RecordNotFound() {
		return nil, errCannotLoadUsers
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Error("Error while qu", err)
		}
		return nil, errUnableToRead
	}
	return users, nil
}

func (u data) create(newUser *Model) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Create(newUser).GetErrors() {
		dbError = err
		log.Errorf("Error while persisting new user in DB %v", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) creteOtp(otp *Otp) error{
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Create(otp).GetErrors() {
		dbError = err
		log.Errorf("Error while persisting new otp in DB %v", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) removeUser(user *Model) error {
	if errs := u.db.Delete(user).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while removing user %v - %v", user.Username, err)
		}
		return errUnableToSave
	}

	return nil
}

func (u data) activate(email string) (err error) {
	if errs := u.db.Exec("UPDATE users SET deleted_at=NULL WHERE username = ?", email).GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while removing user %v - %v", email, err)
		}
		err = errUnableToSave
	}

	return
}

func (u data) updateLink(link *ActionLink) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(link).GetErrors() {
		dbError = err
		log.Errorf("Error while persisting link in DB %v", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) findBy(email string, includeDisabled bool) (Model, error) {
	var user Model
	var record *gorm.DB
	if includeDisabled {
		record = u.db.Unscoped().Where("username = ?", email).Preload("ActionLinks").Find(&user)
	} else {
		record = u.db.Where("username = ?", email).Preload("ActionLinks").Find(&user)
	}

	if record.RecordNotFound() {
		return Model{}, ErrCannotFindUser
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user by email %v - %v", email, err)
		}
		return Model{}, errUnableToRead
	}
	return user, nil
}

func (u data) findUserById(value uint) (Model, error) {
	var user Model
	record := u.db.Where("id = ?", value).Preload("ActionLinks").Find(&user)
	if record.RecordNotFound() {
		return Model{}, nil
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user by email %v - %v", value, err)
		}
		return Model{}, errUnableToRead
	}
	return user, nil
}

func (u data) findOtpForUser(email string) (Otp, error){
	var otp Otp
	record := u.db.Where("username= ?", email).First(&otp)
	if record.RecordNotFound() {
		return Otp{}, nil
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching otp for user by email %v - %v", email, err)
		}
		return Otp{}, errUnableToRead
	}
	return otp, nil
}

func (u data) findLinkByToken(token string) (ActionLink, error) {
	var link ActionLink
	record := u.db.Where("token = ?", token).First(&link)
	if record.RecordNotFound() {
		return ActionLink{}, nil
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching action link by token %v - %v", token, err)
		}
		return ActionLink{}, errUnableToRead
	}
	return link, nil
}

func (u data) updateUser(user *Model) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&user).GetErrors() {
		dbError = err
		log.Error("Error while updating user in DB ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) updateOtp(otp *Otp) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&otp).GetErrors() {
		dbError = err
		log.Error("Error while updating otp in DB ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) updateUserAgent(agent *AgentInfo) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&agent).GetErrors() {
		dbError = err
		log.Error("Error while updating agent in DB ", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) createUserAgent(info AgentInfo) error {
	db := u.db.Begin()
	var dbError error
	for _, err := range db.Save(&info).GetErrors() {
		dbError = err
		log.Errorf("Error while persisting agent info in DB %v", err)
	}
	if dbError != nil {
		db.Rollback()
		return dbError
	}
	db.Commit()

	return nil
}

func (u data) findUserAgentsByUserId(id uint) ([]AgentInfo, error) {
	var agents []AgentInfo
	record := u.db.Where("user_id = ?", id).Find(&agents)
	if record.RecordNotFound() {
		return agents, nil
	}
	if errs := record.GetErrors(); len(errs) > 0 {
		for err := range errs {
			log.Errorf("Error occurred while fetching user agents by userId %v - %v", id, err)
		}
		return nil, errUnableToRead
	}
	return agents, nil
}
