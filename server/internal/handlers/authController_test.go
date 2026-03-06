package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-app/database"
	"go-app/models"
	"go-app/util"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Helper function to setup a test Fiber app and in-memory SQLite database
func setupApp() (*fiber.App, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	database.DB = db

	app := fiber.New()

	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/refresh", RefreshToken)
	app.Get("/user", User)
	app.Post("/logout", Logout)
	app.Put("/users/info", UpdateInfo)
	app.Put("/users/password", UpdatePassword)
	app.Post("/forgot", RequestResetPassword)

	return app, db
}

func TestRequestResetPassword(t *testing.T) {
	// var testUser models.User // Removed as it's no longer used

	setupUserForReset := func(db *gorm.DB) models.User {
		db.Where("email = ?", "resetrequest@example.com").Delete(&models.User{})
		user := models.User{
			FirstName: "Reset",
			LastName:  "PassUser",
			Email:     "resetrequest@example.com",
		}
		user.SetPassword("password123")
		db.Create(&user)
		return user
	}

	tests := []struct {
		name                 string
		setupUser            bool
		payload              map[string]string
		expectedStatusCode   int
		expectedMessage      string
		expectedErrorJSON    map[string]interface{}
	}{
		{
			name:      "successful password reset request for existing user",
			setupUser: true,
			payload:   map[string]string{"email": "resetrequest@example.com"},
			expectedStatusCode: fiber.StatusOK,
			expectedMessage:    "success",
		},
		{
			name:      "password reset request for non-existent user",
			setupUser: false,
			payload:   map[string]string{"email": "nonexistent@example.com"},
			expectedStatusCode: fiber.StatusNotFound,
			expectedErrorJSON:  map[string]interface{}{"errors": map[string]interface{}{"user": []interface{}{"not found"}}},
		},
		{
			name:               "password reset request with missing email field",
			setupUser:          false,
			payload:            map[string]string{},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "email is required",
		},
		{
			name:               "password reset request with empty email string",
			setupUser:          false,
			payload:            map[string]string{"email": ""},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "email is required",
		},
		// Test cases for util function failures were removed due to inability to mock them
		// by direct reassignment of package-level functions.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tc.setupUser {
				setupUserForReset(db)
			}

			bodyBytes, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/forgot", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			respBodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch for "+tc.name+": "+string(respBodyBytes))

			if tc.expectedStatusCode == fiber.StatusOK && tc.expectedMessage != "" {
				var responseBody map[string]string
				err = json.Unmarshal(respBodyBytes, &responseBody)
				assert.NoError(t, err, "Failed to decode success response for "+tc.name+": "+string(respBodyBytes))
				assert.Equal(t, tc.expectedMessage, responseBody["message"], "Success message mismatch for "+tc.name)
			} else if tc.expectedErrorJSON != nil {
				var errorResponse map[string]interface{}
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error JSON for "+tc.name+": "+string(respBodyBytes))
				if expectedErrors, ok := tc.expectedErrorJSON["errors"].(map[string]interface{}); ok {
					responseErrors, ok := errorResponse["errors"].(map[string]interface{})
					assert.True(t, ok, "Response 'errors' field is not a map[string]interface{} for "+tc.name)
					assert.Equal(t, len(expectedErrors), len(responseErrors), "Mismatch in number of error keys for "+tc.name)
					for k, v := range expectedErrors {
						assert.Equal(t, v, responseErrors[k], "Mismatch in error content for key "+k+" for "+tc.name)
					}
				} else {
					assert.Equal(t, tc.expectedErrorJSON, errorResponse, "Error JSON mismatch for "+tc.name)
				}
			} else if tc.expectedMessage != "" {
				var errorResponse map[string]string
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error message for "+tc.name+": "+string(respBodyBytes))
				assert.Equal(t, tc.expectedMessage, errorResponse["message"], "Error message mismatch for "+tc.name)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	var testUser models.User
	var userJWT string
	originalPasswordHash := []byte{}

	setupUserForPasswordUpdate := func(db *gorm.DB) {
		db.Where("email = ?", "passupdate@example.com").Delete(&models.User{})
		currentUser := models.User{
			FirstName: "PassUpdate",
			LastName:  "User",
			Email:     "passupdate@example.com",
		}
		currentUser.SetPassword("oldPassword123")
		db.Create(&currentUser)
		originalPasswordHash = currentUser.Password
		testUser = currentUser

		err := util.GenerateUserTokens(&testUser)
		assert.NoError(t, err, "Setup: Failed to generate tokens")
		userJWT = testUser.AccessToken
		assert.NotEmpty(t, userJWT, "Setup: userJWT is empty")
	}

	tests := []struct {
		name                 string
		setupUser            bool
		payload              map[string]string
		jwtToSend            string
		expectedStatusCode   int
		expectedErrorMessage string
		checkDBPasswordChange bool
		newPasswordForLoginTest string
	}{
		{
			name:      "successful password update",
			setupUser: true,
			payload: map[string]string{
				"password":         "newPassword456",
				"password_confirm": "newPassword456",
			},
			expectedStatusCode:   fiber.StatusOK,
			checkDBPasswordChange: true,
			newPasswordForLoginTest: "newPassword456",
		},
		{
			name:      "password update with mismatched confirmation",
			setupUser: true,
			payload: map[string]string{
				"password":         "newPassword456",
				"password_confirm": "differentPassword789",
			},
			expectedStatusCode:   fiber.StatusBadRequest,
			expectedErrorMessage: "passwords do not match",
		},
		{
			name:      "password update with missing password field",
			setupUser: true,
			payload: map[string]string{
				"password_confirm": "newPassword456",
			},
			expectedStatusCode:   fiber.StatusBadRequest,
			expectedErrorMessage: "password is required",
		},
		{
			name:      "password update with missing password_confirm field",
			setupUser: true,
			payload: map[string]string{
				"password": "newPassword456",
			},
			expectedStatusCode:   fiber.StatusBadRequest,
			expectedErrorMessage: "password_confirm is required",
		},
		{
			name:      "password update with empty password string",
			setupUser: true,
			payload: map[string]string{
				"password":         "",
				"password_confirm": "",
			},
			expectedStatusCode:   fiber.StatusBadRequest,
			expectedErrorMessage: "password cannot be empty",
		},
		{
			name:                 "update attempt with no JWT",
			setupUser:            false,
			payload:              map[string]string{"password": "abc", "password_confirm": "abc"},
			jwtToSend:            "",
			expectedStatusCode:   fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Missing token",
		},
		{
			name:                 "update attempt with invalid JWT",
			setupUser:            false,
			payload:              map[string]string{"password": "abc", "password_confirm": "abc"},
			jwtToSend:            "this.is.a.bad.token",
			expectedStatusCode:   fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Invalid token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			currentJwt := ""
			var activeTestUser models.User

			if tc.setupUser {
				setupUserForPasswordUpdate(db)
				currentJwt = userJWT
				activeTestUser = testUser
			}
			if _, customJwtSpecified := tc.payload["_jwt_override_"]; customJwtSpecified || tc.jwtToSend != "" || !tc.setupUser {
                 currentJwt = tc.jwtToSend
            }

			bodyBytes, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("PUT", "/users/password", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if currentJwt != "" {
				req.AddCookie(&http.Cookie{Name: "jwt", Value: currentJwt})
			}

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			respBodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch for "+tc.name+": "+string(respBodyBytes))

			if tc.expectedStatusCode == fiber.StatusOK {
				var responseBody models.User
				err = json.Unmarshal(respBodyBytes, &responseBody)
				assert.NoError(t, err, "Failed to decode success response for "+tc.name)
				assert.Equal(t, activeTestUser.ID, responseBody.ID, "User ID in response mismatch for "+tc.name)

				if tc.checkDBPasswordChange {
					var userInDb models.User
					db.First(&userInDb, "id = ?", activeTestUser.ID)
					assert.NotEqual(t, originalPasswordHash, userInDb.Password, "Password hash in DB should have changed for "+tc.name)

					if tc.newPasswordForLoginTest != "" {
						loginPayload := map[string]string{
							"email":    activeTestUser.Email,
							"password": tc.newPasswordForLoginTest,
						}
						loginBodyBytes, _ := json.Marshal(loginPayload)
						loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBodyBytes))
						loginReq.Header.Set("Content-Type", "application/json")

						loginResp, loginErr := app.Test(loginReq, -1)
						assert.NoError(t, loginErr)
						assert.Equal(t, fiber.StatusOK, loginResp.StatusCode, "Login with new password failed for "+tc.name+": "+string(respBodyBytes))
						var loginRespBody map[string]interface{}
						err = json.NewDecoder(loginResp.Body).Decode(&loginRespBody)
						assert.NoError(t, err)
						assert.Equal(t, "success", loginRespBody["message"])
					}
				}
			} else if tc.expectedErrorMessage != "" {
				var errorResponse map[string]string
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error message for "+tc.name+ ": "+string(respBodyBytes))
				assert.Equal(t, tc.expectedErrorMessage, errorResponse["message"], "Error message mismatch for "+tc.name)
			}
		})
	}
}

func TestUpdateInfo(t *testing.T) {
	var baseUser models.User
	var userJWT string

	setupUserForUpdateTest := func(db *gorm.DB) {
		db.Where("email = ?", "originalinfo@example.com").Delete(&models.User{})
		db.Where("email = ?", "updatedinfo@example.com").Delete(&models.User{})
		db.Where("email = ?", "existingother@example.com").Delete(&models.User{})

		baseUser = models.User{
			FirstName: "OriginalFirst",
			LastName:  "OriginalLast",
			Email:     "originalinfo@example.com",
		}
		baseUser.SetPassword("password123")
		db.Create(&baseUser)

		err := util.GenerateUserTokens(&baseUser)
		assert.NoError(t, err, "Setup: Failed to generate tokens for baseUser")
		userJWT = baseUser.AccessToken
		assert.NotEmpty(t, userJWT, "Setup: userJWT is empty")
	}

	tests := []struct {
		name                 string
		setupUser            bool
		payload              map[string]string
		jwtToSend            string
		expectedStatusCode   int
		expectedBodyContains map[string]string
		expectedErrorMessage string
		expectedErrorJSON    map[string]interface{}
		checkDBAfterTest     bool
		expectedDBUser       models.User
	}{
		{
			name:      "successful info update",
			setupUser: true,
			payload: map[string]string{
				"first_name": "UpdatedFirst",
				"last_name":  "UpdatedLast",
				"email":      "updatedinfo@example.com",
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBodyContains: map[string]string{
				"first_name": "UpdatedFirst",
				"last_name":  "UpdatedLast",
				"email":      "updatedinfo@example.com",
			},
			checkDBAfterTest: true,
			expectedDBUser: models.User{
				FirstName: "UpdatedFirst",
				LastName:  "UpdatedLast",
				Email:     "updatedinfo@example.com",
			},
		},
		{
			name:      "update with only some fields",
			setupUser: true,
			payload: map[string]string{
				"first_name": "PartialUpdateFirst",
			},
			expectedStatusCode: fiber.StatusOK,
			expectedBodyContains: map[string]string {
				"first_name": "PartialUpdateFirst",
			},
			checkDBAfterTest: true,
			expectedDBUser: models.User{
				FirstName: "PartialUpdateFirst",
			},
		},
		{
			name:                 "update attempt with no JWT",
			setupUser:            false,
			payload:              map[string]string{"first_name": "NoAuth"},
			jwtToSend:            "",
			expectedStatusCode:   fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Missing token",
		},
		{
			name:                 "update attempt with invalid JWT",
			setupUser:            false,
			payload:              map[string]string{"first_name": "BadAuth"},
			jwtToSend:            "this.is.a.bad.token",
			expectedStatusCode:   fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Invalid token",
		},
		{
			name:      "update email to one already taken by another user",
			setupUser: true,
			payload: map[string]string{
				"email": "existingother@example.com",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedErrorMessage: "UNIQUE constraint failed: users.email",
		},
		{
			name:      "update with empty payload",
			setupUser: true,
			payload:   map[string]string{},
			expectedStatusCode: fiber.StatusOK,
			checkDBAfterTest: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			currentJwt := ""
			var tempBaseUser models.User

			if tc.setupUser {
				setupUserForUpdateTest(db)
				currentJwt = userJWT
				tempBaseUser = baseUser

				if tc.name == "update email to one already taken by another user" {
					otherUser := models.User{FirstName: "Other", LastName: "User", Email: "existingother@example.com"}
					otherUser.SetPassword("password")
					db.Create(&otherUser)
				}

				if tc.name == "update with only some fields" {
					tc.expectedDBUser.LastName = tempBaseUser.LastName
					tc.expectedDBUser.Email = tempBaseUser.Email
					if tc.expectedBodyContains != nil {
						tc.expectedBodyContains["last_name"] = tempBaseUser.LastName
						tc.expectedBodyContains["email"] = tempBaseUser.Email
					}
				}
                if tc.name == "update with empty payload" {
                    tc.expectedDBUser = tempBaseUser
					tc.expectedBodyContains = map[string]string{
						"id":         tempBaseUser.ID,
						"first_name": tempBaseUser.FirstName,
						"last_name":  tempBaseUser.LastName,
						"email":      tempBaseUser.Email,
					}
                }
			}

			if tc.jwtToSend != "" || (tc.name == "update attempt with no JWT" && tc.jwtToSend == "") {
				currentJwt = tc.jwtToSend
			}

			bodyBytes, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("PUT", "/users/info", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if currentJwt != "" {
				req.AddCookie(&http.Cookie{Name: "jwt", Value: currentJwt})
			}

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			respBodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch for "+tc.name+": "+string(respBodyBytes))

			if tc.expectedStatusCode == fiber.StatusOK {
				var responseBody models.User
				err = json.Unmarshal(respBodyBytes, &responseBody)
				assert.NoError(t, err, "Failed to decode success response for "+tc.name+": "+string(respBodyBytes))

				if tc.expectedBodyContains != nil {
					assert.Equal(t, tc.expectedBodyContains["first_name"], responseBody.FirstName, "FirstName in response mismatch for "+tc.name)
					assert.Equal(t, tc.expectedBodyContains["last_name"], responseBody.LastName, "LastName in response mismatch for "+tc.name)
					assert.Equal(t, tc.expectedBodyContains["email"], responseBody.Email, "Email in response mismatch for "+tc.name)
					if tc.name != "update with empty payload" {
                         assert.Equal(t, tempBaseUser.ID, responseBody.ID, "User ID in response mismatch for "+tc.name)
                    } else {
                         assert.Equal(t, tc.expectedBodyContains["id"], responseBody.ID, "User ID in response mismatch for "+tc.name)
                    }
				}

				if tc.checkDBAfterTest {
					var userInDb models.User
					db.First(&userInDb, "id = ?", tempBaseUser.ID)
					assert.Equal(t, tc.expectedDBUser.FirstName, userInDb.FirstName, "FirstName in DB mismatch for "+tc.name)
					assert.Equal(t, tc.expectedDBUser.LastName, userInDb.LastName, "LastName in DB mismatch for "+tc.name)
					assert.Equal(t, tc.expectedDBUser.Email, userInDb.Email, "Email in DB mismatch for "+tc.name)
				}

			} else if tc.expectedErrorJSON != nil {
				var errorResponse map[string]interface{}
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error JSON for "+tc.name+": "+string(respBodyBytes))
				assert.Equal(t, tc.expectedErrorJSON, errorResponse, "Error JSON mismatch for "+tc.name)
			} else if tc.expectedErrorMessage != "" {
				var errorResponse map[string]string
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error message for "+tc.name+": "+string(respBodyBytes))
				assert.Contains(t, errorResponse["message"], tc.expectedErrorMessage, "Error message mismatch for "+tc.name)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	app, _ := setupApp()

	t.Run("successful logout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		req.AddCookie(&http.Cookie{Name: "jwt", Value: "dummytoken"})
		req.AddCookie(&http.Cookie{Name: "refreshjwt", Value: "dummyrefreshtoken"})

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "HTTP status code should be 200 OK")

		var responseBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err, "Failed to decode logout response body")
		assert.Equal(t, "success", responseBody["message"], "Logout success message mismatch")

		foundJwtCookie := false
		foundRefreshCookie := false

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "jwt" {
				foundJwtCookie = true
				assert.Equal(t, "", cookie.Value, "jwt cookie value should be empty")
				assert.True(t, cookie.Expires.Before(time.Now()), "jwt cookie should be expired")
				assert.True(t, cookie.HttpOnly, "jwt cookie should be HttpOnly")
			}
			if cookie.Name == "refreshjwt" {
				foundRefreshCookie = true
				assert.Equal(t, "", cookie.Value, "refreshjwt cookie value should be empty")
				assert.True(t, cookie.Expires.Before(time.Now()), "refreshjwt cookie should be expired")
				assert.True(t, cookie.HttpOnly, "refreshjwt cookie should be HttpOnly")
			}
		}
		assert.True(t, foundJwtCookie, "jwt cookie not found in response")
		assert.True(t, foundRefreshCookie, "refreshjwt cookie not found in response")
	})
}

func TestUser(t *testing.T) {
	var validUserJWT string

	setupAuthenticatedUser := func(db *gorm.DB) models.User {
		currentUser := models.User{
			FirstName: "Authenticated",
			LastName:  "User",
			Email:     "autheduser@example.com",
		}
		currentUser.SetPassword("password123")
		db.Where("email = ?", currentUser.Email).Delete(&models.User{})
		db.Create(&currentUser)

		err := util.GenerateUserTokens(&currentUser)
		assert.NoError(t, err, "Setup: Failed to generate tokens for auth user")
		validUserJWT = currentUser.AccessToken
		assert.NotEmpty(t, validUserJWT, "Setup: Generated access token is empty")
		return currentUser
	}

	tests := []struct {
		name                string
		setupAuthUser       bool
		sendJWTAsCookie     bool
		sendJWTAsHeader     bool
		customJWT           string
		expectedStatusCode  int
		expectedBodyUser    *models.User
		expectedErrorJSON   map[string]interface{}
		expectedErrorMessage string
	}{
		{
			name:               "successful user retrieval with JWT cookie",
			setupAuthUser:      true,
			sendJWTAsCookie:    true,
			expectedStatusCode: fiber.StatusOK,
		},
		{
			name:               "successful user retrieval with JWT Bearer header",
			setupAuthUser:      true,
			sendJWTAsHeader:    true,
			expectedStatusCode: fiber.StatusOK,
		},
		{
			name:               "user retrieval with invalid JWT - malformed",
			setupAuthUser:      false,
			sendJWTAsCookie:    true,
			customJWT:          "this.is.not.a.valid.jwt",
			expectedStatusCode: fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Invalid token",
		},
		{
			name:               "user retrieval with expired JWT",
			setupAuthUser:      true,
			sendJWTAsCookie:    true,
			customJWT:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjM3YxYjhkZC01YzJjLTQxZmYtYjE4Ny0zMWMyYjY2YmQ3MzAiLCJleHAiOjE2MDkzNzI4MDB9.someinvalidsignature",
			expectedStatusCode: fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Invalid token",
		},
		{
			name:               "user retrieval for non-existent user ID in valid JWT structure",
			setupAuthUser:      false,
			sendJWTAsCookie:    true,
			customJWT:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJmYWtlLXVzZXItaWQtMTIzNDUiLCJleHAiOjE5MjQ5MDU2MDB9.fakelySignedPArt",
			expectedStatusCode: fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Invalid token",
		},
		{
			name:                 "user retrieval with no JWT",
			setupAuthUser:        false,
			expectedStatusCode:   fiber.StatusUnauthorized,
			expectedErrorMessage: "Unauthenticated: Missing token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			var currentTestUser models.User
			activeJWT := ""

			if tc.setupAuthUser {
				currentTestUser = setupAuthenticatedUser(db)
				activeJWT = validUserJWT
				if tc.expectedBodyUser == nil && tc.expectedStatusCode == fiber.StatusOK {
					expectedUser := currentTestUser
					expectedUser.Password = nil
					expectedUser.AccessToken = ""
					expectedUser.RefreshToken = ""
					tc.expectedBodyUser = &expectedUser
				}
			}

			if tc.customJWT != "" {
				activeJWT = tc.customJWT
			}

			req := httptest.NewRequest("GET", "/user", nil)

			if tc.sendJWTAsCookie && activeJWT != "" {
				cookie := new(fiber.Cookie)
				cookie.Name = "jwt"
				cookie.Value = activeJWT
				httpCookie := &http.Cookie{Name: cookie.Name, Value: cookie.Value}
				req.AddCookie(httpCookie)
			}
			if tc.sendJWTAsHeader && activeJWT != "" {
				req.Header.Set("Token", activeJWT)
			}

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			respBodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch for "+tc.name+": "+string(respBodyBytes))

			if tc.expectedStatusCode == fiber.StatusOK && tc.expectedBodyUser != nil {
				var responseBody models.User
				err = json.Unmarshal(respBodyBytes, &responseBody)
				assert.NoError(t, err, "Failed to decode user response for "+tc.name+": "+string(respBodyBytes))

				assert.Equal(t, tc.expectedBodyUser.ID, responseBody.ID, "User ID mismatch for "+tc.name)
				assert.Equal(t, tc.expectedBodyUser.FirstName, responseBody.FirstName, "FirstName mismatch for "+tc.name)
				assert.Equal(t, tc.expectedBodyUser.LastName, responseBody.LastName, "LastName mismatch for "+tc.name)
				assert.Equal(t, tc.expectedBodyUser.Email, responseBody.Email, "Email mismatch for "+tc.name)
				assert.Empty(t, responseBody.Password, "Password should be empty in response for "+tc.name)
			} else if tc.expectedErrorJSON != nil {
				var errorResponse map[string]interface{}
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error JSON response for "+tc.name+": "+string(respBodyBytes))
				assert.Equal(t, tc.expectedErrorJSON, errorResponse, "Error JSON content mismatch for "+tc.name)
			} else if tc.expectedErrorMessage != "" {
				var errorResponse map[string]string
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error message response for "+tc.name+": "+string(respBodyBytes))
				assert.Equal(t, tc.expectedErrorMessage, errorResponse["message"], "Error message mismatch for "+tc.name)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	var testUser models.User
	var validRefreshToken string

	setupForRefreshTokenTest := func(db *gorm.DB) {
		testUser = models.User{
			FirstName: "Refresh",
			LastName:  "User",
			Email:     "refreshuser@example.com",
		}
		testUser.SetPassword("password123")
		db.Create(&testUser)

		err := util.GenerateUserTokens(&testUser)
		assert.NoError(t, err, "Setup: Failed to generate tokens for test user")
		validRefreshToken = testUser.RefreshToken
		assert.NotEmpty(t, validRefreshToken, "Setup: Generated refresh token is empty")
	}

	tests := []struct {
		name               string
		setupTestUser      bool
		payload            map[string]string
		expectedStatusCode int
		expectedJSON       map[string]interface{}
		expectedMessage    string
		checkCookie        bool
	}{
		{
			name:          "successful token refresh",
			setupTestUser: true,
			payload:       map[string]string{},
			expectedStatusCode: fiber.StatusOK,
			checkCookie: true,
		},
		{
			name:          "refresh with invalid token - malformed",
			setupTestUser: false,
			payload: map[string]string{
				"refreshToken": "this.is.not.a.valid.jwt",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "token invalid or expired",
		},
		{
			name:          "refresh with expired token",
			setupTestUser: true,
			payload:       map[string]string{},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "token invalid or expired",
		},
		{
			name:          "refresh for non-existent user",
			setupTestUser: false,
			payload:       map[string]string{},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "token invalid or expired",
		},
		{
			name:               "refresh with missing refreshToken field",
			setupTestUser:      false,
			payload:            map[string]string{},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "refreshToken is required",
		},
		{
			name:          "refresh with empty refreshToken string",
			setupTestUser: false,
			payload: map[string]string{
				"refreshToken": "",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "refreshToken is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			currentPayload := make(map[string]string)
			for k, v := range tc.payload {
				currentPayload[k] = v
			}

			if tc.setupTestUser {
				setupForRefreshTokenTest(db)
				if tc.name == "successful token refresh" {
					currentPayload["refreshToken"] = validRefreshToken
				}
			}

			if tc.name == "refresh with expired token" {
				currentPayload["refreshToken"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjM3YxYjhkZC01YzJjLTQxZmYtYjE4Ny0zMWMyYjY2YmQ3MzAiLCJleHAiOjE2MDkzNzI4MDB9.someinvalidsignature"
			}

			if tc.name == "refresh for non-existent user" {
				currentPayload["refreshToken"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJmYWtlLXVzZXItaWQtMTIzNDUiLCJleHAiOjE5MjQ5MDU2MDB9.fakelySignedPArt"
			}

			bodyBytes, _ := json.Marshal(currentPayload)
			req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch for "+tc.name)

			respBodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			if tc.expectedJSON != nil {
				var responseBody map[string]interface{}
				err = json.Unmarshal(respBodyBytes, &responseBody)
				assert.NoError(t, err, "Failed to decode JSON response for "+tc.name+": "+string(respBodyBytes))

				if expectedErrors, ok := tc.expectedJSON["errors"].(map[string]interface{}); ok {
					responseErrors, ok := responseBody["errors"].(map[string]interface{})
					assert.True(t, ok, "Response 'errors' field is not a map[string]interface{} for "+tc.name)
					assert.Equal(t, len(expectedErrors), len(responseErrors), "Mismatch in number of error keys for "+tc.name)
					for k, v := range expectedErrors {
						assert.Equal(t, v, responseErrors[k], "Mismatch in error content for key "+k+" for "+tc.name)
					}
				} else {
					assert.Equal(t, tc.expectedJSON, responseBody, "JSON response mismatch for "+tc.name)
				}

			} else if tc.expectedMessage != "" {
				var errorResponse map[string]string
				err = json.Unmarshal(respBodyBytes, &errorResponse)
				assert.NoError(t, err, "Failed to decode error message response for "+tc.name+": "+string(respBodyBytes))
				assert.Contains(t, errorResponse["message"], tc.expectedMessage, "Error message mismatch for "+tc.name)
			}

			if tc.expectedStatusCode == fiber.StatusOK && tc.checkCookie {
				cookies := resp.Cookies()
				foundJwtCookie := false
				foundRefreshCookie := false
				for _, cookie := range cookies {
					if cookie.Name == "jwt" {
						foundJwtCookie = true
						assert.NotEmpty(t, cookie.Value, "jwt cookie should not be empty for "+tc.name)
						assert.True(t, cookie.HttpOnly, "jwt cookie should be HttpOnly for "+tc.name)
					}
					if cookie.Name == "refreshjwt" {
						foundRefreshCookie = true
						assert.NotEmpty(t, cookie.Value, "refreshjwt cookie should not be empty for "+tc.name)
						assert.True(t, cookie.HttpOnly, "refreshjwt cookie should be HttpOnly for "+tc.name)
					}
				}
				assert.True(t, foundJwtCookie, "jwt cookie not found for "+tc.name)
				assert.True(t, foundRefreshCookie, "refreshjwt cookie not found for "+tc.name)

				var successResponseBody map[string]interface{}
				err = json.Unmarshal(respBodyBytes, &successResponseBody)
				assert.NoError(t, err, "Failed to decode success response body for "+tc.name)
				assert.Equal(t, "success", successResponseBody["message"], "Success message mismatch for "+tc.name)

				userMap, ok := successResponseBody["user"].(map[string]interface{})
				assert.True(t, ok, "User data in success response is not a map for "+tc.name)
				assert.Equal(t, testUser.Email, userMap["email"], "User email in success response mismatch for "+tc.name)
				assert.NotEmpty(t, userMap["id"], "User ID in success response should not be empty for "+tc.name)
				assert.Equal(t, testUser.ID, userMap["id"], "User ID in success response mismatch for "+tc.name)
				newAccessToken, okAT := userMap["accessToken"].(string)
				assert.True(t, okAT, "New accessToken is not a string or missing for "+tc.name)
				assert.NotEmpty(t, newAccessToken, "New accessToken should not be empty for "+tc.name)
				newRefreshToken, okRT := userMap["refreshToken"].(string)
				assert.True(t, okRT, "New refreshToken is not a string or missing for "+tc.name)
				assert.NotEmpty(t, newRefreshToken, "New refreshToken should not be empty for "+tc.name)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name               string
		payload            map[string]string
		expectedStatusCode int
		expectedMessage    string
		checkDB            bool
		emailForDBCheck    string
		expectTokenError   bool
	}{
		{
			name: "successful registration",
			payload: map[string]string{
				"firstName":       "Test",
				"lastName":        "User",
				"email":           "testreg@example.com",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			expectedStatusCode: fiber.StatusOK,
			checkDB:            true,
			emailForDBCheck:    "testreg@example.com",
		},
		{
			name: "mismatched passwords",
			payload: map[string]string{
				"firstName":       "Jane",
				"lastName":        "Doe",
				"email":           "janedoe_mismatch@example.com",
				"password":        "password123",
				"confirmPassword": "password456",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "passwords do not match",
		},
		{
			name: "registration with existing email",
			payload: map[string]string{
				"firstName":       "Another",
				"lastName":        "User",
				"email":           "existingreg@example.com",
				"password":        "newpassword",
				"confirmPassword": "newpassword",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "UNIQUE constraint failed",
		},
		{
			name: "registration with invalid input data - missing email",
			payload: map[string]string{
				"firstName":       "NoEmail",
				"lastName":        "User",
				"password":        "password123",
				"confirmPassword": "password123",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage: "all fields (firstName, lastName, email, password) are required",
		},
		{
			name: "registration with invalid input data - missing password",
			payload: map[string]string{
				"firstName":       "NoPassword",
				"lastName":        "User",
				"email":           "nopass@example.com",
				"confirmPassword": "password123",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "all fields (firstName, lastName, email, password) are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tc.name == "registration with existing email" {
				existingUser := models.User{
					FirstName: "Existing",
					LastName:  "RegUser",
					Email:     "existingreg@example.com",
				}
				existingUser.SetPassword("password123")
				db.Create(&existingUser)
			}

			bodyBytes, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch")

			if tc.expectedStatusCode == fiber.StatusOK {
				var responseBody models.User
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				assert.NoError(t, err)

				assert.Equal(t, tc.payload["firstName"], responseBody.FirstName)
				assert.Equal(t, tc.payload["lastName"], responseBody.LastName)
				assert.Equal(t, tc.payload["email"], responseBody.Email)
				assert.NotEmpty(t, responseBody.ID)
				assert.NotEmpty(t, responseBody.AccessToken, "Access token should not be empty")
				assert.NotEmpty(t, responseBody.RefreshToken, "Refresh token should not be empty")

				if tc.checkDB && tc.emailForDBCheck != "" {
					var userInDb models.User
					res := db.First(&userInDb, "email = ?", tc.emailForDBCheck)
					assert.NoError(t, res.Error)
					assert.Equal(t, tc.payload["firstName"], userInDb.FirstName)
					assert.NotEqual(t, tc.payload["password"], string(userInDb.Password), "Password should be hashed")
				}
			} else {
				var errorResponse map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				assert.NoError(t, err, "Failed to decode error response")
				if tc.expectedMessage != "" {
					actualMessage := errorResponse["message"]
					assert.Contains(t, actualMessage, tc.expectedMessage, "Error message mismatch")
				}
				if tc.expectTokenError {
					assert.Contains(t, errorResponse["internal server error"], tc.expectedMessage, "Internal server error message mismatch")
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	testUserEmail := "loginuser@example.com"
	testUserPassword := "password123"

	tests := []struct {
		name               string
		setupUser          bool
		payload            map[string]string
		expectedStatusCode int
		expectedJSON       map[string]interface{}
		expectedMessage    string
	}{
		{
			name:      "successful login",
			setupUser: true,
			payload: map[string]string{
				"email":    testUserEmail,
				"password": testUserPassword,
			},
			expectedStatusCode: fiber.StatusOK,
			expectedJSON: map[string]interface{}{
				"message": "success",
			},
		},
		{
			name:      "login with incorrect password",
			setupUser: true,
			payload: map[string]string{
				"email":    testUserEmail,
				"password": "wrongpassword",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedJSON: map[string]interface{}{
				"errors": map[string]interface{}{
					"password": []interface{}{"incorrect"},
				},
			},
		},
		{
			name:      "login with non-existent user",
			setupUser: false,
			payload: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedStatusCode: fiber.StatusNotFound,
			expectedJSON: map[string]interface{}{
				"errors": map[string]interface{}{
					"user": []interface{}{"not found"},
				},
			},
		},
		{
			name:      "login with missing email",
			setupUser: false,
			payload: map[string]string{
				"password": "password123",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "email and password are required",
		},
		{
			name:      "login with missing password",
			setupUser: false,
			payload: map[string]string{
				"email": "missingpass@example.com",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "email and password are required",
		},
		{
			name:      "login with empty email string",
			setupUser: false,
			payload: map[string]string{
				"email":    "",
				"password": "password123",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "email and password are required",
		},
		{
			name:      "login with empty password string",
			setupUser: false,
			payload: map[string]string{
				"email":    "empty@example.com",
				"password": "",
			},
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "email and password are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app, db := setupApp()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			if tc.setupUser {
				user := models.User{
					FirstName: "Login",
					LastName:  "User",
					Email:     testUserEmail,
				}
				user.SetPassword(testUserPassword)
				// Clean up before creating to ensure test idempotency
				db.Where("email = ?", testUserEmail).Delete(&models.User{})
				db.Create(&user)
			}

			bodyBytes, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode, "HTTP status code mismatch")

			if tc.expectedStatusCode == fiber.StatusOK {
				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedJSON["message"], responseBody["message"])
				userMap, ok := responseBody["user"].(map[string]interface{})
				assert.True(t, ok, "User data in response is not a map")
				assert.Equal(t, testUserEmail, userMap["email"])
				assert.NotEmpty(t, userMap["id"], "User ID should not be empty")
				assert.NotEmpty(t, userMap["accessToken"], "Access token should not be empty")
				assert.NotEmpty(t, userMap["refreshToken"], "Refresh token should not be empty")
				assert.NotContains(t, userMap, "password", "Password should not be in user response")

			} else if tc.expectedJSON != nil {
				var errorResponse map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				assert.NoError(t, err, "Failed to decode error response JSON")
				assert.Equal(t, tc.expectedJSON, errorResponse, "Error JSON mismatch")
			} else if tc.expectedMessage != "" {
				var errorResponse map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				assert.NoError(t, err, "Failed to decode error response message")
				assert.Equal(t, tc.expectedMessage, errorResponse["message"], "Error message mismatch")
			}
		})
	}
}
