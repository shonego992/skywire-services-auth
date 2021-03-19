package user

import (
	"time"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
)

type TokenType uint8

const (
	ConfirmRegistration TokenType = 0
	ResetPassword       TokenType = 1
)

type TokenStatus uint8

const (
	NotUsed TokenStatus = 0
	Used    TokenStatus = 1
)

const (
	confirmedUserMask uint8 = 1 // 0000 0001

	createAdminMask uint8 = 128 // 1000 0000
	disableUserMask uint8 = 64  // 0100 0000

	isAdminMask uint8 = 248 // 1111 0000
)

// Model for actions from user - password reset, profile verification
type ActionLink struct {
	ID         uint        `gorm:"primary_key" json:"-"`
	Status     TokenStatus `json:"-"`
	Expiration time.Time   `json:"-"`
	Token      string      `json:"-"`
	Type       TokenType   `json:"-"`
	CreatedAt  time.Time   `json:"-"`
	UpdatedAt  time.Time   `json:"-"`
	DeletedAt  *time.Time  `json:"-"`
	UserId     uint        `json:"-"`
}

type AgentInfo struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Address   string     `json:"-"`
	Client    string     `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	UserId    uint       `json:"-"`
}

// Model for User entity.
type Model struct {
	ID          uint                  `gorm:"primary_key" json:"id" example:"1"`
	Status      uint8                 `json:"status"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   *time.Time            `json:"disabled,omitempty"`
	Username    string                `json:"username" example:"someone@mail.com"`
	Password    string                `json:"password,omitempty"`
	ActionLinks []ActionLink          `json:"-" gorm:"foreignkey:UserId; PRELOAD:false"`
	AgentInfos  []AgentInfo           `json:"-" gorm:"foreignkey:UserId; PRELOAD:false"`
	Rights      []authorization.Right `json:"rights,omitempty" gorm:"-" sql:"-"`
	UseOtp    	bool				  `json:"useOtp"`
}

type Otp struct {
	ID 			uint 		`gorm:"primary_key" json:"id" example:"1"`
	Secret  	string		`json:"-"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"-"`
	DeletedAt   *time.Time  `json:"disabled,omitempty"`
	Username 	string      `json:"-"`
	Expiration  time.Time  	`json:"-"`
}

func (Model) TableName() string {
	return "users"
}

func (m *Model) IsConfirmed() bool {
	return (m.Status & confirmedUserMask) > 0
}

type CreatedAtResponse struct {
	CreatedAt time.Time
}

func (m *Model) Confirm() {
	m.Status |= confirmedUserMask
}

func (m *Model) IsAdmin() bool {
	return (m.Status & isAdminMask) > 0
}

func (m *Model) CanDisableUser() bool {
	return (m.Status & disableUserMask) > 0
}

func (m *Model) SetDisableUser(canDisable bool) {
	m.setRole(disableUserMask, canDisable)
}

func (m *Model) Disable() {
	m.Status &= ^confirmedUserMask
}

func (m *Model) CanCreateAdmin() bool {
	return (m.Status & createAdminMask) > 0
}

func (m *Model) SetCreateAdmin(canCreate bool) {
	m.setRole(createAdminMask, canCreate)
}

func (m *Model) setRole(roleMask uint8, isActive bool) {
	if isActive {
		m.Status |= roleMask
	} else {
		m.Status &= ^roleMask
	}
}
