package user

import (
	"reflect"
	"testing"
)

var s Service

func TestCreate(t *testing.T) {
	tests := []test{
		{
			name:         "Creating with regular values should return no error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser, nil},
			newUser:      stbValidUser,
			expected:     stbValidUserResp,
		},
		{
			name:         "Creating with existing email returns error",
			storeRecords: []interface{}{stbValidUser},
			storeErrors:  []error{nil},
			newUser:      stbValidUser,
			err:          errEmailExists,
		},
		{
			name:    "Creating with too short Password returns error",
			newUser: stbShortPassUser,
			err:     errPasswordTooShort,
		},
		{
			name:    "Creating with empty Password returns error",
			newUser: stbEmptyPassUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with missing Password returns error",
			newUser: stbMissingPassUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with empty Username returns error",
			newUser: stbEmptyEmailUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with missing Username returns error",
			newUser: stbMissingEmailUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with invalid Username returns error",
			newUser: stbInvalidEmailUser,
			err:     errEmailNotValid,
		},
	}

	runTests(t, tests, create)

}

func TestCreateAdmin(t *testing.T) {
	tests := []test{
		{
			name:         "Creating with regular values should return no error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser, nil},
			newUser:      stbValidUser,
			expected:     stbValidAdminResp1,
		},
		{
			name:         "Creating when user is already admin returns error",
			storeRecords: []interface{}{stbValidUserAdmin},
			storeErrors:  []error{nil},
			newUser:      stbValidUser,
			err:          errAdminAlreadyExists,
		},
		{
			name:    "Creating with too short Password returns error",
			newUser: stbShortPassUser,
			err:     errPasswordTooShort,
		},
		{
			name:    "Creating with empty Password returns error",
			newUser: stbEmptyPassUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with missing Password returns error",
			newUser: stbMissingPassUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with empty Username returns error",
			newUser: stbEmptyEmailUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with missing Username returns error",
			newUser: stbMissingEmailUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Creating with invalid Username returns error",
			newUser: stbInvalidEmailUser,
			err:     errEmailNotValid,
		},
	}

	runTests(t, tests, createAdmin)

}

func TestRemoveUser(t *testing.T) {
	tests := []test{
		{
			name:         "Remove user with valid username no error",
			storeRecords: []interface{}{stbValidUser},
			storeErrors:  []error{nil, nil},
			newUser:      stbValidUser,
			expected:     stbValidUser,
		},
		{
			name:         "Remove user with invalid email returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			newUser:      stbInvalidEmailUser,
			err:          ErrCannotFindUser,
		},
		{
			name:         "Remove user with valid email returns tehnical error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{nil, errTechnicalError},
			newUser:      stbValidUser,
			err:          errTechnicalError,
		},
	}

	runTests(t, tests, removeUser)
}

func TestUpdatePassword(t *testing.T) {
	tests := []test{
		{
			name:         "Update with valid password should return valid user response",
			storeRecords: []interface{}{stbValidUserWithHashedPassword, stbValidUserResp},
			storeErrors:  []error{nil, nil},
			newUser:      stbValidUser,
			expected:     stbValidUserResp,
		},
		{
			name:    "Update with missing email returns error",
			newUser: stbEmptyEmailUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:    "Update with missing old password returns error",
			newUser: stbEmptyPassUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:         "Update with non existing user returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			newUser:      stbValidUser,
			err:          ErrCannotFindUser,
		},
		{
			name:         "Update with non confirmed user returns error",
			storeRecords: []interface{}{stbNonConfirmedUser},
			storeErrors:  []error{nil},
			newUser:      stbNonConfirmedUser,
			err:          ErrNotConfirmed,
		},
		{
			name:         "Update with wrong old password returns error",
			storeRecords: []interface{}{stbValidUserWithHashedPassword},
			storeErrors:  []error{nil},
			newUser:      stbInvalidOldPasswordUser,
			err:          errPasswordDoesNotMatch,
		},
	}

	runTests(t, tests, updatePassword)
}

func TestUpdateRights(t *testing.T) {
	tests := []test{
		{
			name:         "Update rights without an error",
			storeRecords: []interface{}{},
			storeErrors:  []error{nil},
			newUser:      stbUpdateRights,
			expected:     stbUpdateRightsResponse,
		},
	}
	runTests(t, tests, updateRights)
}

func TestFindBy(t *testing.T) {
	tests := []test{
		{
			name:         "Finding with no errors",
			storeRecords: []interface{}{stbValidMultipleRights},
			storeErrors:  []error{nil},
			newUser:      stbValidUserQuery,
			expected:     stbValidUserFoundResp,
		},
		{
			name:    "Searching with empty username",
			newUser: stbEmptyEmailUser,
			err:     ErrMissingMandatoryFields,
		},
		{
			name:         "Search for non existing user returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{errUnableToRead},
			newUser:      stbValidUserWithAddress,
			err:          errUnableToRead,
		},
	}
	runTests(t, tests, findBy)
}

func TestForgotPasswordCreateLink(t *testing.T) {
	tests := []test{
		{
			name:         "Sucessful link creation",
			storeRecords: []interface{}{stbValidUser},
			storeErrors:  []error{nil, nil},
			newUser:      stbValidUser,
			err:          nil,
		},
		{
			name:    "Trying to create link with invalid Username returns error",
			newUser: stbInvalidEmailUser,
			err:     errEmailNotValid,
		},
		{
			name:    "Trying to create link with invalid Username returns error",
			newUser: stbInvalidEmailUser,
			err:     errEmailNotValid,
		},
		{
			name:         "Trying to create link for non existing returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			newUser:      stbValidUserWithAddress,
			err:          ErrCannotFindUser,
		},
	}
	runTests(t, tests, forgotPasswordCreateLink)
}

func TestAddUserAgentInfo(t *testing.T) {
	tests := []test{
		{
			name:         "Adding new user agent with no error",
			storeRecords: []interface{}{stbAgent, true},
			storeErrors:  []error{nil, nil},
			newUser:      stbAgentInfoQuery,
			err:          nil,
		},

		{
			name:         "Cannot find user with that id",
			storeRecords: []interface{}{stbEmptyAgentInfo},
			storeErrors:  []error{errUnableToRead},
			newUser:      stbAgentInfoQueryWrongId,
			err:          errUnableToRead,
		},
	}
	runTests(t, tests, addUserAgentInfo)
}

func TestResendValidationTokenForRegistration(t *testing.T) {
	tests := []test{
		{
			name:         "Resending sucessful",
			storeRecords: []interface{}{stbValidUserWithValidToken},
			storeErrors:  []error{nil, nil},
			newUser:      stbValidUserQuery,
			err:          nil,
		},
		{
			name:         "Resending with invalid Username returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{errEmailNotValid},
			newUser:      stbInvalidEmailUser,
			err:          errEmailNotValid,
		},
		{
			name:         "Resending for non existing username returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			newUser:      stbInvalidEmailUser,
			err:          ErrCannotFindUser,
		},
		{
			name:         "Resending for used action link returns error",
			storeRecords: []interface{}{stbValidUserWithUsedActionLink},
			storeErrors:  []error{nil},
			newUser:      stbValidUserQuery,
			err:          errAlreadyConfirmed,
		},
	}
	runTests(t, tests, resendValidationTokenForRegistration)
}

func TestVerify(t *testing.T) {
	tests := []test{

		{
			name:         "No error",
			storeRecords: []interface{}{stbActionLinkUnused, stbValidUserWithValidToken},
			storeErrors:  []error{nil, nil, nil},
			newUser:      stbValidUserWithValidToken,
			err:          nil,
			expected:     stbValidUserWithValidToken,
		},
		{
			name:         "Incorrect token length error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{errIncorrectLengthToken},
			newUser:      stbValidUserWithBadToken,
			err:          errIncorrectLengthToken,
		},
		{
			name:         "Searching for action link for given token error",
			storeRecords: []interface{}{ActionLink{}},
			storeErrors:  []error{errUnableToRead},
			newUser:      stbValidUserWithValidToken,
			err:          errUnableToRead,
			expected:     Model{},
		},
		{
			name:         "Searching for action link expired",
			storeRecords: []interface{}{ActionLink{}},
			storeErrors:  []error{nil},
			newUser:      stbValidUserWithValidToken,
			err:          errTokenExpired,
			expected:     Model{},
		},
		{
			name:         "Already confirmed user error for given token",
			storeRecords: []interface{}{stbActionLinkUsed},
			storeErrors:  []error{nil},
			newUser:      stbValidUserWithValidToken,
			err:          errAlreadyConfirmed,
		},
	}
	runTests(t, tests, verify)
}
func TestResetPassword(t *testing.T) {
	tests := []test{

		{
			name:         "Reset with valid action link no error",
			storeRecords: []interface{}{stbValidUserWithValidToken, stbValidUserWithValidToken},
			storeErrors:  []error{nil, nil},
			newUser:      stbValidUserWithValidTokenNewPassword,
			err:          nil,
			expected:     stbValidUserWithValidToken,
		},
		{
			name:         "Reset with incorrect token length returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{errIncorrectLengthToken},
			newUser:      stbValidUserWithBadToken,
			err:          errIncorrectLengthToken,
		},
		{
			name:         "Reset for non existing user returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			newUser:      stbInvalidEmailUserValidToken,
			err:          ErrCannotFindUser,
		},
		{
			name:         "Reset for non confirmed user returns error",
			storeRecords: []interface{}{stbNonConfirmedUserWithToken},
			storeErrors:  []error{ErrNotConfirmed},
			newUser:      stbNonConfirmedUserWithToken,
			err:          ErrNotConfirmed,
		},
		{
			name:         "Reset with wrong token returns error",
			storeRecords: []interface{}{stbValidUserWithValidToken},
			storeErrors:  []error{nil},
			newUser:      stbValidUserWithNonExistingToken,
			err:          errWrongTokenAndMailCombination,
		},
		{
			name:         "Reset with expired token returns error",
			storeRecords: []interface{}{stbValidUserWithExpiredActionLink},
			storeErrors:  []error{nil},
			newUser:      stbValidUserWithExpiredActionLink,
			err:          errTokenExpired,
		},
		{
			name:         "Reset with used token returns error",
			storeRecords: []interface{}{stbValidUserWithUsedActionLink},
			storeErrors:  []error{nil},
			newUser:      stbValidUserWithUsedActionLink,
			err:          errAlreadyConfirmed,
		},
	}
	runTests(t, tests, resetPassword)
}
func TestGetAdmins(t *testing.T) {
	tests := []test{
		{
			name:         "No error",
			storeRecords: []interface{}{stbAdmins},
			storeErrors:  []error{nil},
			expected:     stbValidAdminResp,
			err:          nil,
		},
		{
			name:         "Failed due to : unable to load users",
			storeRecords: []interface{}{[]Model{}},
			storeErrors:  []error{errCannotLoadUsers},

			err: errCannotLoadUsers,
		},
		{
			name:         "Failed due to:  unable to read",
			storeRecords: []interface{}{[]Model{}},
			storeErrors:  []error{errUnableToRead},

			err: errCannotLoadUsers,
		},
	}
	runTests(t, tests, getAdmins)
}

func TestGetUsers(t *testing.T) {
	tests := []test{
		{
			name:         "No error",
			storeRecords: []interface{}{stbUsers},
			storeErrors:  []error{nil},
			expected:     stbValidUserResp,
			err:          nil,
		},
		{
			name:         "Failed due to : unable to load users",
			storeRecords: []interface{}{[]Model{}},
			storeErrors:  []error{errCannotLoadUsers},

			err: errCannotLoadUsers,
		},
		{
			name:         "Failed due to:  unable to read",
			storeRecords: []interface{}{[]Model{}},
			storeErrors:  []error{errUnableToRead},

			err: errCannotLoadUsers,
		},
	}
	runTests(t, tests, getUsers)
}
func TestActivateUser(t *testing.T) {
	tests := []test{
		{
			name:        "No issues in activating user",
			storeErrors: []error{nil},
			newUser:     stbValidUser,
			expected:    Model{},
			err:         nil,
		},

		{
			name:    "Creating with empty Username returns error",
			newUser: stbEmptyEmailUser,
			err:     ErrMissingMandatoryFields,
		},

		{
			name:        "Database saving error",
			storeErrors: []error{errUnableToSave},
			newUser:     stbValidUser,
			err:         errUnableToSave,
		},
	}
	runTests(t, tests, activateUser)
}

type usrSrvc func(svc Service, usr *Model) (Model, error)

func create(svc Service, usr *Model) (Model, error) {
	return *usr, svc.Create(usr)
}

func createAdmin(svc Service, usr *Model) (Model, error) {
	_, err := svc.CreateAdmin(usr)
	return *usr, err
}

func removeUser(svc Service, usr *Model) (Model, error) {
	return *usr, svc.RemoveUser(usr.Username)
}

func updatePassword(svc Service, usr *Model) (Model, error) {
	return svc.UpdatePassword(usr.Username, usr.Password, "SomeNewPassword")
}

func updateRights(svc Service, usr *Model) (Model, error) {
	return *usr, svc.UpdateRights(usr, usr.Rights)
}

func findBy(svc Service, usr *Model) (Model, error) {
	return svc.FindBy(usr.Username)
}

func forgotPasswordCreateLink(svc Service, usr *Model) (Model, error) {
	_, err := svc.ForgotPasswordCreateLink(usr.Username)
	return Model{}, err
}

func getAdmins(svc Service, usr *Model) (Model, error) {
	var mod Model
	models, err := svc.getAdmins()
	if len(models) > 0 {
		mod = models[0]
	}
	return mod, err
}
func getUsers(svc Service, usr *Model) (Model, error) {
	var mod Model
	models, err := svc.GetUsers()
	if len(models) > 0 {
		mod = models[0]
	}
	return mod, err
}

func addUserAgentInfo(svc Service, usr *Model) (Model, error) {
	var client, ipAddress string
	var uid uint
	if len(usr.AgentInfos) > 0 {
		client = usr.AgentInfos[0].Client
		ipAddress = usr.AgentInfos[0].Address
		uid = usr.AgentInfos[0].UserId
	}
	_, err := svc.AddUserAgentInfo(client, ipAddress, uid)
	return Model{}, err
}

func resendValidationTokenForRegistration(svc Service, usr *Model) (Model, error) {
	_, err := svc.ResendValidationTokenForRegistration(usr.Username)
	return Model{}, err
}

func verify(svc Service, usr *Model) (Model, error) {
	var str string
	if len(usr.ActionLinks) > 0 {
		str = usr.ActionLinks[0].Token
	}
	resp, err := svc.Verify(str)
	return resp, err
}

func resetPassword(svc Service, usr *Model) (Model, error) {
	req := resetPassReq{Password: usr.Password, Email: usr.Username}
	if len(usr.ActionLinks) > 0 {
		req.Token = usr.ActionLinks[0].Token
	}
	resp, err := svc.ResetPassword(req)
	return resp, err
}

func activateUser(svc Service, usr *Model) (Model, error) {
	err := svc.ActivateUser(usr.Username)
	return Model{}, err

}

func runTests(t *testing.T, tests []test, run usrSrvc) {
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			svc := NewService(newStub(test.storeErrors, test.storeRecords))
			resp, err := run(svc, &test.newUser)
			if err != nil && test.err == nil {
				t.Errorf("%s failed, expected no error but received error %v", test.name, err)
			} else if test.err != nil && err == nil {
				t.Errorf("%s failed, expected error %v but no error was received", test.name, test.err)
			} else if test.err != nil && err != nil {
				if err != test.err {
					t.Errorf("%s failed, expected error %v but received error %v", test.name, test.err, err)
				}
			} else {
				if !reflect.DeepEqual(test.expected, resp) {
					t.Errorf("%s failed, expected: %#v - actual %#v", test.name, test.expected, resp)
				}
			}
		})
	}
}

type test struct {
	name         string
	storeRecords []interface{}
	storeErrors  []error
	newUser      Model
	err          error
	expected     Model
}
