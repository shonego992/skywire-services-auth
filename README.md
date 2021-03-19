# skywire-user-system

Skywire User System API
=======================
This is the Skywire User System service.

**Version:** 1.0

Instructions for running the project:
1. run postgres_docker.sh
2. copy config-example into config.toml
3. ensure you have golang/dep installed on your PATH
4. run dep ensure
5. go run cmd/main.go

# Skywire User System API
This is a Skywire User System service.

## Version: 1.0



### /admin

#### POST
##### Summary:

Create a new Admin in the system

##### Description:

Collect provided Admin attributes from the body and create new Admin in the system

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| newAdmin | body | New Admin | Yes | [user.Model](#user.model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | [user.Model](#user.model) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /admin/admins

#### GET
##### Summary:

List all admins

##### Description:

Method for admins to get list of all admins

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [user.Model](#user.model) ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /admin/admins/

#### POST
##### Summary:

Update admin rights

##### Description:

Changes rights of a user according to the request

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| updated | body | Model with new rights | Yes | [user.Model](#user.model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [user.Model](#user.model) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /admin/users

#### GET
##### Summary:

List all users

##### Description:

Returns the list of all current users

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [user.Model](#user.model) ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /admin/users/

#### DELETE
##### Summary:

Removes user

##### Description:

Removes user for given username

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | query | Mail of user to be removed | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

#### GET
##### Summary:

Returns user

##### Description:

Returns user found by username provided

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| username | query | User's email | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [user.Model](#user.model) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |

### /auth/info

#### GET
##### Summary:

Retrieve signed in User's info

##### Description:

Information about currently signed in user is collected and returned as response.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [user.Model](#user.model) |
| 401 | Unauthorized | [api.ErrorResponse](#api.errorresponse) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users

#### POST
##### Summary:

Create a new User in the system

##### Description:

Collect provided User attributes from the body and create new User in the system

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| newUser | body | New User | Yes | [user.Model](#user.model) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | [user.Model](#user.model) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/confirmOTP

#### GET
##### Summary:

Confirm a setup for OTP code on your account

##### Description:

Send a otp code to the backend, and have it set up on your account, so it will be required for future actions

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/disableOtp

#### POST
##### Summary:

Disable OTP code for your account

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/forgotPassword

#### GET
##### Summary:

Start forgot password flow

##### Description:

If User has forgotten his password this endpoint enables him to reset password using link sent to provided email address.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| email | query | User's email | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/password

#### PATCH
##### Summary:

Update Users's password

##### Description:

Collect, validate and store User's new Skycoin address.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| newAddress | body | User's new Skycoin address | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [user.Model](#user.model) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/resendValidationToken

#### GET
##### Summary:

Re-send welcome mail containing verification token

##### Description:

In case User misplaces original welcome mail new one is sent to the registered email address

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| email | query | User's email address | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |

### /users/resetPassword

#### POST
##### Summary:

Process forgot password request

##### Description:

If User has forgotten his password and received password reset link this endpoint validates input and stores new password.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| forgotPasswordRequest | body | User's forgot password input | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/setupOTP

#### GET
##### Summary:

Request a setup for OTP code for your account

##### Description:

Request otp to be set up for your account

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/verify

#### GET
##### Summary:

Verify User's email address based on token provided via email

##### Description:

After User creates account first available action is to verify account's email address.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| token | query | User's token for email validation | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [user.Model](#user.model) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### Models


#### api.ErrorResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### user.Model

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| disabled | string |  | No |
| id | integer |  | No |
| password | string |  | No |
| rights | string |  | No |
| status | integer |  | No |
| useOtp | boolean |  | No |
| username | string |  | No |