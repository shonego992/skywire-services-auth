import {Injectable} from '@angular/core';

import {environment} from '../../environments/environment';
import {HttpService} from './http.service';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from './shared.service';
import {AuthService} from './auth.service';
import {User} from '../models/user.model';
import {AuthData} from '../models/auth-data.model';
import {Router} from '@angular/router';
import {Observable, Subject} from 'rxjs/Rx';
import {HttpClient, HttpHeaders} from '@angular/common/http';

@Injectable()
export class UserService {

  constructor(private httpService: HttpService,  private sharedService: SharedService, private authService: AuthService,
              private router: Router, private httpClient: HttpClient) {}
  public tokenNeeded: Subject<any> = new Subject<any>();

  public updateAddress (newAddress: string) {
    const data = {
      address: newAddress
    };
    this.httpService.patchToUrl<any>(environment.service + ApiRoutes.USER.Address, data).subscribe(
        (res: any) => {
          this.sharedService.showSuccess('Address successfully updated');
        },
        (err: any) => {
          this.sharedService.showError('Can\'t update user\'s address', err.split(': ')[1]);
          // TODO  customize messages
        }
      );
    }

  public updatePassword(oldPass: string, newPass: string) {
    const data = {
      oldPassword: oldPass,
      password: newPass
    };
    this.httpService.patchToUrl<any>(environment.service + ApiRoutes.USER.Password, data).subscribe(
      (res: any) => {
        this.sharedService.showSuccess('Password successfully updated');
        this.sharedService.sleep(1000);
        this.authService.logout();
        this.router.navigate(['/success-message']);
      },
      (err: any) => {
        this.sharedService.showError('Can\'t update user\'s password', err.split(': ')[1]);
        // TODO  customize messages
      }
    );
  }

  public resetPass(email: string, onComplete: any) {
    let success = false;
    return this.httpService.getFromUrl(environment.service + ApiRoutes.USER.ForgotPassword + email).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Email with instructions was sent successfully');
        success = true;
      },
      (err: any) => {
        this.sharedService.showError('Can\'t reset new user', err.split(': ')[1]);
      },
      () => onComplete(success)
    );
  }

  public changePassword (email, resetToken, password) {
    const data = {
      email: email,
      token: resetToken,
      password: password
    };
    return this.httpService.postToUrl(environment.service + ApiRoutes.USER.ResetPassword, data).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Password was updated successfully');
        this.authService.logout();
      },
      (err: any) => {
        this.sharedService.showError('Can\'t reset user\'s password', err.split(': ')[1]);
      }
    );
  }

  public verifyProfile (token) {
    return this.httpService.getFromUrl(environment.service + ApiRoutes.USER.VerifyProfile + token).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Profile was validated successfully');
        this.authService.setUserVefified(true);
        if (this.authService.isAuth()) {
          this.sharedService.sleep(1000);
          if (this.authService.isAdmin()) {
            this.router.navigate(['/admins-list']);
            return;
          }
          this.router.navigate(['/account-info']);
          return;
        }
        this.router.navigate(['/']);
      },
      (err: any) => {
        this.sharedService.showError('Can\'t verify user\'s profile', err.split(': ')[1]);
      }
    );
  }

  public registerUser(authData: AuthData, onComplete: any) {
    let success = false;
    this.httpService.postToUrl<AuthData>(environment.service + ApiRoutes.USER.Users, authData).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Registration successful. Please check your email to confirm profile');
        this.sharedService.sleep(1000);
        success = true;
      },
      (err: any) => {
        this.sharedService.showError('Can\'t sign up new user', err.split(': ')[1]);

        // TODO  customize messages
        // switch (err) {
        //   case "user service: provided email is already taken by another user": {this.showMessage('Can\'t register new user', err.split(': ')[1]); break;}
        //   case "user service: provided email is already taken by another user": {this.showMessage('Can\'t register new user', err.split(': ')[1]); break;}
        // }
      }, () => onComplete(success)
    );
  }

  public login (username: string, password: string, token?: string) {
    if (!token) {
      token = '';
    }
    const data = {
      username: username,
      password: password
    };
    this.httpClient
      .post(environment.service + '/auth/login', data, {
        headers: new HttpHeaders().set('2fa', token),
      })
      .subscribe((res: any) => {
        if (res && res.token && res.expire) {
          this.authService.saveToken(res);
          this.authService.authSuccessfully();
          this.authService.refreshUserData();
        } else {
          this.sharedService.showError('Unexpected error on sign in', res);
        }
      },
      (err: any) => {
        var errString = err.split(': ')[1];
        if (errString) {
          if (errString.indexOf('2FA') !== -1) {
            this.tokenNeeded.next(true);
            return;
          }
          this.sharedService.showError('Can\'t sign in', errString);
        }
      }
    );
   }

  public resendVerificationLink(mail: string) {
    this.httpService.getFromUrl(environment.service + ApiRoutes.USER.ResendToken + mail).subscribe((user: User) => {
        this.sharedService.showSuccess('Verification link resent');
      },
      (err: any) => {
        this.sharedService.showError('Can\'t resend mail', err.split(': ')[1]);
      }
    );
  }
}
