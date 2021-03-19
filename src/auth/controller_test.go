package auth

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	//prevents cors and jwt errors
	viper.SetDefault("cors.allow-all", true)
	viper.SetDefault("jwt.realm", "test realm")
	viper.SetDefault("jwt.key", "test key")

	// Run the other tests
	os.Exit(m.Run())
}

//TODO: Needs stubbed jwt verification mechanism to complete this set of tests
/*
func TestInfoController(t *testing.T) {
	url, method := "/auth/info", "GET"
	tests := []struct {
		name       string
		url        string
		method     string
		statusCode int
		response   user.Model
		errMessage error
	}{
		{
			name:       "Get Info with no valid token should result in error",
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "Get info with ",
			statusCode: http.StatusUnauthorized,
		},
		//TODO add more tests
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := app.NewServer(DefaultController()) //TODO mock here
			requestBody, e := json.Marshal("")
			if e != nil {
				t.Fatalf("Unable to prepare request body %v", e)
			}
			req, err := http.NewRequest(method, "/api/v1"+url, bytes.NewReader(requestBody))
			if err != nil {
				t.Fatalf("Unable to prepare request %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Origin", "http://localhost")

			//w is used to capture the response from the server
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Fatalf("Response code should be %d, was: %d", test.statusCode, w.Code)
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}

			if w.Code == http.StatusOK {
				actualResponse := user.Model{}
				if err = json.Unmarshal(body, &actualResponse); err != nil {
					t.Fatalf("Failed to unmarshal response body %v due to error : %s", body, err)
				}

				if !reflect.DeepEqual(test.response, actualResponse) {
					t.Fatalf("Expected %v, actual %v", test.response, actualResponse)
				}
			} else if w.Code == http.StatusBadRequest {
				var errorResponse api.ErrorResponse
				if err = json.Unmarshal(body, &errorResponse); err != nil {
					t.Fatalf("Failed to unmarshal response %v due to error : %s", errorResponse, err)
				}
				if msg := test.errMessage.Error(); msg != errorResponse.Error {
					t.Fatalf("Expected %v, actual %v", msg, errorResponse.Error)
				}
			}
		})
	}
}
*/
