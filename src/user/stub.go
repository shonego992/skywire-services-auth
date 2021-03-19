package user

import (
	"time"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
)

//stub implements store interface for mocking purposes
type stub struct {
	data []interface{}
	err  []error
}

func newStub(simulatedErrors []error, data []interface{}) *stub {
	return &stub{
		data: data,
		err:  simulatedErrors,
	}
}

var stbAgent = []AgentInfo{{UserId: 42, Client: "client", Address: "0.0.0.0"}}
var stbValidUser = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1}
var stbValidUserResp = Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 1}
var stbValidUserAdmin = Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 193}
var stbValidAdminResp = Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 65}
var stbValidAdminResp1 = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 65}
var stbValidUserWithAddress = Model{Username: "test1@mail.com", Password: "", Status: 1}
var stbValidUserWithAddressResp = Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 1}
var stbShortPassUser = Model{Username: "test1@mail.com", Password: "Pass", Status: 1}
var stbInvalidEmailUser = Model{Username: "test1", Password: "Password", Status: 1}
var stbEmptyPassUser = Model{Username: "test1@mail.com", Password: "", Status: 1}
var stbEmptyPassUserResp = Model{Username: "test1@mail.com", Password: "", ID: 1}
var stbEmptyEmailUser = Model{Username: "", Password: "Password", Status: 1}
var stbMissingPassUser = Model{Username: "test1@mail.com", Status: 1}
var stbMissingEmailUser = Model{Password: "Password", Status: 1}
var stbEmptyEmailAndPassUser = Model{Username: "", Password: "", Status: 1}
var stbMissingEmailAndPassUser = Model{}
var stbNonConfirmedUser = Model{Username: "test1@mail.com", Password: "Password", Status: 0}

/* Password hash composed from user + pass = test1@gmail.comPassword */
var stbValidUserWithHashedPassword = Model{Username: "test1@mail.com", Password: "$2a$14$zQdcE.WkYxGuAWECXE/HNOH4tv9jfaAOqe5kQC9pyh5kx1qyGlcnO", ID: 42, Status: 1}
var stbInvalidOldPasswordUser = Model{Username: "test1@mail.com", Password: "DifferentOldPassword", ID: 42, Status: 1}

var stbUpdateRights = Model{Username: "test1@mail.com", Status: 193, Rights: []authorization.Right{{Name: "create_user", Value: true}, {Name: "disable_user", Value: true}}}
var stbUpdateRightsResponse = Model{Username: "test1@mail.com", Status: 193, Rights: []authorization.Right{{Name: "create_user", Value: true}, {Name: "disable_user", Value: true}}}
var stbValidMultipleRights = Model{Username: "test1@mail.com", Password: "Password", Status: 0xC1, ID: 42}
var stbValidUserQuery = Model{Username: "test1@gmail.com"}

var stbValidUserFoundResp = Model{Username: "test1@mail.com", Password: "Password", Status: 0xC1, ID: 42, Rights: []authorization.Right{{Name: "create_user", Label: "Create Admin", Value: true}, {Name: "disable_user", Label: "Disable User", Value: true}}}
var stbValidUserWithAPIKeys = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1}

var stbInvalidKey = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1}

var stbValidUserWithAgent = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1, AgentInfos: []AgentInfo{{UserId: 43, Client: "client1", Address: "0.0.0.1"}}}
var stbAgentInfoQuery = Model{AgentInfos: []AgentInfo{{UserId: 42, Client: "client", Address: "0.0.0.0"}}}
var stbAgentInfoQueryWrongId = Model{AgentInfos: []AgentInfo{{UserId: 43, Client: "client", Address: "0.0.0.0"}}}

var dummyValidActionLinkToken = "fortyfortyfortyfortyfortyfortyfortyforty"
var stbActionLinkUsed = ActionLink{Token: dummyValidActionLinkToken, Status: Used, Expiration: time.Date(2045, 11, 17, 20, 34, 58, 651387237, time.UTC)}
var stbActionLinkUnused = ActionLink{Token: dummyValidActionLinkToken, Status: NotUsed, Expiration: time.Date(2045, 11, 17, 20, 34, 58, 651387237, time.UTC)}
var stbActionLinkExpired = ActionLink{Token: dummyValidActionLinkToken, Status: NotUsed, Expiration: time.Date(2017, 11, 17, 20, 34, 58, 651387237, time.UTC)}

var stbValidUserWithValidToken = Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 1, ActionLinks: []ActionLink{stbActionLinkUnused, stbActionLinkUnused}}
var stbValidUserWithValidTokenNewPassword = Model{Username: "test1@mail.com", Password: "$2a$14$zQdcE.WkYxGuAWECXE/HNOH4tv9jfaAOqe5kQC9pyh5kx1qyGlcnO", ID: 42, Status: 1, ActionLinks: []ActionLink{stbActionLinkUnused}}
var stbValidUserWithBadToken = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1, ActionLinks: []ActionLink{{Token: "wewewew"}}}
var stbValidUserWithNonExistingToken = Model{Username: "test1@gmail.com", Password: "Password", ID: 42, Status: 1, ActionLinks: []ActionLink{{Token: "fortyfortyfortyfortyfortyfortyfortyfdiff"}}}
var stbValidUserWithExpiredActionLink = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1, ActionLinks: []ActionLink{stbActionLinkExpired}}
var stbValidUserWithUsedActionLink = Model{Username: "test1@mail.com", Password: "Password", ID: 42, Status: 1, ActionLinks: []ActionLink{stbActionLinkUsed}}
var stbInvalidEmailUserValidToken = Model{Username: "test1", Password: "Password", ID: 42, Status: 1, ActionLinks: []ActionLink{{Token: dummyValidActionLinkToken}}}
var stbNonConfirmedUserWithToken = Model{Username: "test1@mail.com", Password: "Password", ActionLinks: []ActionLink{{Token: dummyValidActionLinkToken}}}
var stbEmptyAgentInfo []AgentInfo

var stbAdmins = []Model{Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 65}, Model{Username: "test2@mail.com", Password: "", ID: 43, Status: 65}}
var stbUsers = []Model{Model{Username: "test1@mail.com", Password: "", ID: 42, Status: 1}, Model{Username: "test2@mail.com", Password: "", ID: 43, Status: 1}}

type EmailStr struct {
	email string
}

func (cs *stub) create(newUser *Model) error {
	newUser.ID = 42
	return cs.returnError()
}

func (cs *stub) update(newUser *Model) error {
	newUser.ID = 42
	return cs.returnError()
}

func (cs *stub) findBy(email string) (Model, error) {
	return cs.returnRecord().(Model), cs.returnError()
}

func (cs *stub) removeUser(user *Model) error {
	return cs.returnError()
}

func (cs *stub) findUserById(id uint) (Model, error) {
	return cs.returnRecord().(Model), cs.returnError()
}

func (cs *stub) findLinkByToken(token string) (ActionLink, error) {
	return cs.returnRecord().(ActionLink), cs.returnError()
}

func (cs *stub) updateUser(user *Model) error {
	return cs.returnError()
}

func (cs *stub) updateLink(link *ActionLink) error {
	return cs.returnError()
}

func (cs *stub) findUserWithAPIKeys(email string) (Model, error) {
	return cs.returnRecord().(Model), cs.returnError()
}

func (cs *stub) createUserAgent(info AgentInfo) error {
	return cs.returnError()
}

func (cs *stub) findUserAgentsByUserId(id uint) ([]AgentInfo, error) {
	return cs.returnRecord().([]AgentInfo), cs.returnError()
}

func (cs *stub) updateUserAgent(agent *AgentInfo) error {
	return cs.returnError()
}

func (cs *stub) findUserWithApplications(email string) (Model, error) {
	return cs.returnRecord().(Model), cs.returnError()
}

func (cs *stub) getUsers() ([]Model, error) {
	return cs.returnRecord().([]Model), cs.returnError()
}

func (cs *stub) getAdmins() ([]Model, error) {
	return cs.returnRecord().([]Model), cs.returnError()
}

func (cs *stub) activate(email string) error {
	return cs.returnError()

}
func (cs *stub) returnRecord() interface{} {
	response := cs.data[0]
	cs.data = cs.data[1:]
	return response
}

func (cs *stub) returnError() error {
	err := cs.err[0]
	cs.err = cs.err[1:]
	return err
}
