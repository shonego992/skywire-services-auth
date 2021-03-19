package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/mattbaird/gochimp"
	"github.com/SkycoinPro/skywire-services-auth/src/api"
	"github.com/SkycoinPro/skywire-services-auth/src/app"
	"github.com/SkycoinPro/skywire-services-auth/src/template"

	"github.com/gin-gonic/gin"
)

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func TestCreateController(t *testing.T) {
	url, method := "/users", "POST"
	tests := []struct {
		name       string
		url        string
		method     string
		body       *Model
		statusCode int
		store      *stub
		response   Model
		errMessage error
	}{
		{
			name:       "Create new user with valid parameters and not used email",
			body:       &stbValidUser,
			statusCode: http.StatusCreated,
			store:      newStub([]error{ErrCannotFindUser, nil}, []interface{}{Model{}}),
			response:   stbValidUserResp,
		},
		{
			name:       "Create new user with valid parameters and used email",
			body:       &stbValidUser,
			statusCode: http.StatusBadRequest,
			store:      newStub([]error{nil}, []interface{}{stbValidUser}),
			errMessage: errEmailExists,
		},
		{
			name:       "Create new user with invalid email",
			body:       &stbInvalidEmailUser,
			statusCode: http.StatusBadRequest,
			errMessage: errEmailNotValid,
		},
		{
			name:       "Create new user with short password",
			body:       &stbShortPassUser,
			statusCode: http.StatusBadRequest,
			errMessage: errPasswordTooShort,
		},
		{
			name:       "Create new user with missing password",
			body:       &stbMissingPassUser,
			statusCode: http.StatusUnprocessableEntity,
			errMessage: ErrMissingMandatoryFields,
		},
		{
			name:       "Create new user with empty password",
			body:       &stbEmptyPassUser,
			statusCode: http.StatusUnprocessableEntity,
			errMessage: ErrMissingMandatoryFields,
		},
		{
			name:       "Create new user with missing email",
			body:       &stbMissingEmailUser,
			statusCode: http.StatusUnprocessableEntity,
			errMessage: ErrMissingMandatoryFields,
		},
		{
			name:       "Create new user with empty email",
			body:       &stbEmptyEmailUser,
			statusCode: http.StatusUnprocessableEntity,
			errMessage: ErrMissingMandatoryFields,
		},
		{
			name:       "Create new user with missing email and password",
			body:       &stbMissingEmailAndPassUser,
			statusCode: http.StatusUnprocessableEntity,
			errMessage: ErrMissingMandatoryFields,
		},
		{
			name:       "Create new user with empty email and password",
			body:       &stbEmptyEmailAndPassUser,
			statusCode: http.StatusUnprocessableEntity,
			errMessage: ErrMissingMandatoryFields,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(NewController(NewService(test.store), template.NewService(&gochimp.MandrillAPI{})))
			requestBody, e := json.Marshal(test.body)
			if e != nil {
				t.Errorf("Unable to prepare request body %v", e)
			}
			req, err := http.NewRequest(method, "/api/v1"+url, bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("Unable to prepare request %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {
				actualResponse := Model{}
				if err = json.Unmarshal(body, &actualResponse); err != nil {
					t.Errorf("Failed to unmarshal response body %v due to error : %s", body, err)
				}

				if !reflect.DeepEqual(test.response, actualResponse) {
					t.Errorf("Expected %v, actual %v", test.response, actualResponse)
				}
			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error(); msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			}
		})
	}
}
func TestVerifyController(t *testing.T) {
	url, method := "/users/verify", "GET"
	tests := []struct {
		name         string
		url          string
		method       string
		params       map[string]string
		statusCode   int
		storeRecords []interface{}
		storeErrors  []error
		response     Model
		errMessage   api.ErrorResponse
	}{
		{name: "Create admin with no error",
			params: map[string]string{
				"token": "fortyfortyfortyfortyfortyfortyfortyforty",
			},
			statusCode:   http.StatusOK,
			storeErrors:  []error{nil, nil, nil},
			storeRecords: []interface{}{stbActionLinkUnused, stbValidUserWithValidToken},
			response:     stbValidUserResp,
		},

		{
			name: "Empty token",

			statusCode: http.StatusBadRequest,
			errMessage: api.ErrorResponse{Error: "Incorrect value for confirmation token sent"},
		},
		{
			name: "Already confirmed user error for given token",
			params: map[string]string{
				"token": "fortyfortyfortyfortyfortyfortyfortyforty",
			},
			storeRecords: []interface{}{stbActionLinkUsed},
			storeErrors:  []error{nil},
			statusCode:   http.StatusInternalServerError,
			errMessage:   api.ErrorResponse{Error: errAlreadyConfirmed.Error()},
		},
		{
			name:         "Searching for action link expired",
			storeRecords: []interface{}{ActionLink{}},
			storeErrors:  []error{nil},
			params: map[string]string{
				"token": "fortyfortyfortyfortyfortyfortyfortyforty",
			},
			statusCode: http.StatusInternalServerError,
			errMessage: api.ErrorResponse{Error: errTokenExpired.Error()},
		},
		{
			name: "Searching for action link for given token error",
			params: map[string]string{
				"token": "fortyfortyfortyfortyfortyfortyfortyforty",
			},
			storeRecords: []interface{}{ActionLink{}},
			storeErrors:  []error{errUnableToRead},
			statusCode:   http.StatusInternalServerError,
			errMessage:   api.ErrorResponse{Error: errUnableToRead.Error()},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(NewController(NewService(newStub(test.storeErrors, test.storeRecords)), template.NewService(&gochimp.MandrillAPI{})))

			req, err := http.NewRequest(method, "/api/v1"+url, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			for k, v := range test.params {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Content-Type", "application/json")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {
				actualResponse := Model{}
				if err = json.Unmarshal(body, &actualResponse); err != nil {
					t.Errorf("Failed to unmarshal response body %v due to error : %s", body, err)
				}

				if !reflect.DeepEqual(test.response, actualResponse) {
					t.Errorf("Expected %v, actual %v", test.response, actualResponse)
				}
			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			} else if w.Code == http.StatusInternalServerError {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}

			}
		})
	}
}

func TestResendValidationTokenController(t *testing.T) {
	url, method := "/users/resendValidationToken", "GET"
	tests := []struct {
		name         string
		url          string
		method       string
		params       map[string]string
		statusCode   int
		storeRecords []interface{}
		storeErrors  []error
		response     Model
		errMessage   api.ErrorResponse
	}{
		{
			name:         "Resending sucessful",
			storeRecords: []interface{}{stbValidUserWithValidToken},
			storeErrors:  []error{nil, nil},
			statusCode:   http.StatusOK,
			params: map[string]string{
				"email": "test1@gmail.com",
			},
		},
		{
			name:         "Resending for non existing username returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			statusCode:   http.StatusBadRequest,
			params: map[string]string{
				"email": "test2@gmail.com",
			},
			errMessage: api.ErrorResponse{Error: "Incorrect value for user email sent"},
		},
		{
			name:         "Resending with invalid Username returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{errEmailNotValid},
			statusCode:   http.StatusBadRequest,
			params: map[string]string{
				"email": "test2@gmail.com",
			},
			errMessage: api.ErrorResponse{Error: "Incorrect value for user email sent"},
		},
		{
			name:         "Resending for used action link returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{errAlreadyConfirmed},
			statusCode:   http.StatusBadRequest,
			params: map[string]string{
				"email": "test2",
			},
			errMessage: api.ErrorResponse{Error: "Incorrect value for user email sent"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(NewController(NewService(newStub(test.storeErrors, test.storeRecords)), template.NewService(&gochimp.MandrillAPI{})))

			req, err := http.NewRequest(method, "/api/v1"+url, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			for k, v := range test.params {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Content-Type", "application/json")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {

			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			} else if w.Code == http.StatusInternalServerError {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}

			}
		})
	}
}
func TestForgotPasswordController(t *testing.T) {
	url, method := "/users/forgotPassword", "GET"
	tests := []struct {
		name         string
		url          string
		method       string
		params       map[string]string
		statusCode   int
		storeRecords []interface{}
		storeErrors  []error
		response     Model
		errMessage   api.ErrorResponse
	}{
		{
			name:         "Reset with valid action link no error",
			storeRecords: []interface{}{stbValidUser},
			storeErrors:  []error{nil, nil},
			statusCode:   http.StatusOK,
			params: map[string]string{
				"email": "test1@gmail.com",
			},
		},
		{
			name:         "Trying to create link for non existing returns error",
			storeRecords: []interface{}{Model{}},
			storeErrors:  []error{ErrCannotFindUser},
			statusCode:   http.StatusInternalServerError,
			params: map[string]string{
				"email": "test1@gmail.com",
			},
			errMessage: api.ErrorResponse{Error: ErrCannotFindUser.Error()},
		},
		{
			name: "Trying to create link for non existing returns error",

			statusCode: http.StatusInternalServerError,
			params: map[string]string{
				"email": "test1",
			},
			errMessage: api.ErrorResponse{Error: errEmailNotValid.Error()},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(NewController(NewService(newStub(test.storeErrors, test.storeRecords)), template.NewService(&gochimp.MandrillAPI{})))

			req, err := http.NewRequest(method, "/api/v1"+url, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			for k, v := range test.params {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Content-Type", "application/json")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {

			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			} else if w.Code == http.StatusInternalServerError {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}

			}
		})
	}
}

func TestResetPasswordController(t *testing.T) {
	url, method := "/users/resetPassword", "POST"
	tests := []struct {
		name       string
		url        string
		method     string
		body       map[string]string
		statusCode int
		store      *stub

		errMessage error
	}{
		{
			name:       "Reset with incorrect token length returns error",
			body:       map[string]string{"email": "test1@gmail.com", "token": "0", "newpassword": "SomeNewPassword"},
			statusCode: http.StatusInternalServerError,
			store:      newStub([]error{errIncorrectLengthToken}, []interface{}{Model{}}),

			errMessage: errIncorrectLengthToken,
		},
		{
			name:       "Reset for non existing user returns error",
			body:       map[string]string{"email": "test1@gmail.com", "token": "fortyfortyfortyfortyfortyfortyfortyforty", "newpassword": "SomeNewPassword"},
			statusCode: http.StatusInternalServerError,
			store:      newStub([]error{ErrCannotFindUser}, []interface{}{stbNonConfirmedUserWithToken}),
			errMessage: ErrCannotFindUser,
		},
		{
			name:       "Reset for non confirmed user returns error",
			body:       map[string]string{"email": "test1@gmail.com", "token": "fortyfortyfortyfortyfortyfortyfortyforty", "newpassword": "SomeNewPassword"},
			statusCode: http.StatusInternalServerError,
			store:      newStub([]error{ErrNotConfirmed}, []interface{}{stbNonConfirmedUserWithToken}),
			errMessage: ErrNotConfirmed,
		},
		{
			name:       "Reset with wrong token returns error",
			body:       map[string]string{"email": "test1@gmail.com", "token": "fortyfortyfortyfortyfortyfortyfortyfdiff", "newpassword": "SomeNewPassword"},
			statusCode: http.StatusInternalServerError,
			store:      newStub([]error{nil}, []interface{}{stbValidUserWithValidToken}),
			errMessage: errWrongTokenAndMailCombination,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(NewController(NewService(test.store), template.NewService(&gochimp.MandrillAPI{})))
			requestBody, e := json.Marshal(test.body)
			if e != nil {
				t.Errorf("Unable to prepare request body %v", e)
			}
			req, err := http.NewRequest(method, "/api/v1"+url, bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("Unable to prepare request %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {

			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error(); msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			} else if w.Code == http.StatusInternalServerError {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error(); msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			}
		})
	}
} /* TODO: Needs stubbed jwt verification mechanism to complete this set of tests
func TestUpdatePasswordController(t *testing.T) {
	url, method := "/users/password", "PATCH"
	tests := []struct {
		name       string
		url        string
		method     string
		body       map[string]string
		statusCode int
		store      *stub
		response   Model
		errMessage error
	}{
		{
			name:       "Update with valid password should return valid user response",
			body:       map[string]string{"email": "test1@gmail.com", "oldPassword": "Password", "newPassword": "SomeNewPassword"},
			statusCode: http.StatusCreated,
			store:      newStub([]error{nil, nil}, []interface{}{stbValidUserWithHashedPassword, stbValidUserResp}),
			response:   stbValidUserResp,
			errMessage: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(NewController(NewService(test.store), template.NewService(&gochimp.MandrillAPI{})))
			requestBody, e := json.Marshal(test.body)
			if e != nil {
				t.Errorf("Unable to prepare request body %v", e)
			}
			req, err := http.NewRequest(method, "/api/v1"+url, bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("Unable to prepare request %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusCreated {
				actualResponse := Model{}
				if err = json.Unmarshal(body, &actualResponse); err != nil {
					t.Errorf("Failed to unmarshal response body %v due to error : %s", body, err)
				}

				if !reflect.DeepEqual(test.response, actualResponse) {
					t.Errorf("Expected %v, actual %v", test.response, actualResponse)
				}

			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error(); msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			} else if w.Code == http.StatusInternalServerError {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error(); msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			}
		})
	}
}
*/
/* TODO: Needs stubbed jwt verification mechanism to complete this set of tests
func TestGetAdminsController(t *testing.T) {
	url, method := "/admin/admins", "GET"
	tests := []struct {
		name         string
		url          string
		method       string
		params       map[string]string
		statusCode   int
		storeRecords []interface{}
		storeErrors  []error
		response     []Model
		errMessage   api.ErrorResponse
	}{
		{
			name:         "No error",
			storeRecords: []interface{}{stbAdmins[0], stbAdmins},
			storeErrors:  []error{nil},
			statusCode:   http.StatusOK,
			response:     stbAdmins,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()
			server := app.NewServer(NewController(NewService(newStub(test.storeErrors, test.storeRecords)), template.NewService(&gochimp.MandrillAPI{})))

			req, err := http.NewRequest(method, "/api/v1"+url, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			for k, v := range test.params {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InRlc3QxQGdtYWlsLmNvbSJ9.Yq_Wr8HWryv4r86VbDIUflamUZEXD-Tf9s8Ir2H4byE")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {

			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			} else if w.Code == http.StatusInternalServerError {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error; msg != errorResponse.Error {
					t.Errorf("Expected %v, actual %v", msg, errorResponse.Error)
				}

			}
		})
	}
}
*/
