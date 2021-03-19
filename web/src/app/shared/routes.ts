// Shared variables to be used in application

export class ApiRoutes {

  public static USER = {
    'Address' : '/users/address',
    'Password' : '/users/password',
    'Users' : '/users',
    'ResetPassword': '/users/resetPassword',
    'ForgotPassword': '/users/forgotPassword?email=',
    'VerifyProfile': '/users/verify?token=',
    'Keys': '/users/keys',
    'ResendToken': '/users/resendValidationToken?email=',
    'Application': '/whitelist/application',
    'UpdateApplication': '/whitelist/updateApplication',
    'ApplicationNoImages': '/whitelist/updateApplicationNoImageChange',
    'SetupOTP': '/users/setupOTP',
    'ConfirmOTP': '/users/confirmOTP?otp=',
    'DisableOTP': '/users/disableOTP'
  };

  public static AUTH = {
    'Refresh': '/auth/refresh'
  };

  public static ADMIN = {
    'Whitelists': '/whitelist/whitelists',
    'Whitelist': '/whitelist/whitelist',
    'Admin': '/admin',
    'Users': '/admin/users',
    'AdminList': '/admin/admins'
  };
}
