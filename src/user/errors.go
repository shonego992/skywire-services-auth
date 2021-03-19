package user

import "errors"

var errPasswordTooShort = errors.New("user service: provided password is too short (min 8 chars) to be used")
var errEmailExists = errors.New("user service: provided email is already taken by another user")
var errAdminAlreadyExists = errors.New("user service: admin with provided email already exists")
var errEmailNotValid = errors.New("user service: provided email has no valid format")
var ErrMissingMandatoryFields = errors.New("user service: missing some mandatory fields")
var errUnableToSave = errors.New("user service: unable to persist provided data")
var errUnableToRead = errors.New("user service: unable to query persisted data")
var errIncorrectLengthToken = errors.New("user service: provided confirmation token has wrong length")
var errAlreadyConfirmed = errors.New("user service: user already confirmed")
var ErrNotConfirmed = errors.New("user service: user is not confirmed, please verify email first") //TODO close this out
var ErrCannotFindUser = errors.New("user service: cannot find user by email")
var ErrCannotFindOtp = errors.New("user service: cannot find 2FA for user")
var ErrCannotGenerateKey = errors.New("user service: cannot generate otp key for user")
var ErrWrongCode = errors.New("user service: wrong token")
var ErrCannotDisable2FA = errors.New("user service: cannot disable two factor authentication")
var Err2FADisabled = errors.New("user service: 2fa for user is disabled")

var errWrongTokenAndMailCombination = errors.New("user service: wrong combination of email and token for password reset")
var errTokenExpired = errors.New("user service: user action token expired")
var errPasswordDoesNotMatch = errors.New("user service: old and new password are not the same")
var errTechnicalError = errors.New("user service: technical error occured")
var errCannotSaveUserAgent = errors.New("user service: cannot create user agent")
var errUnableToProcessRequest = errors.New("user controller: unable to process fields from the request")
var errForbidden = errors.New("user controller: forbidden from completing the request")
var errCannotLoadUsers = errors.New("user controller: cannot load users")
