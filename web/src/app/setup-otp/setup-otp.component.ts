import { Component, OnInit } from '@angular/core';
import {HttpService} from '../services/http.service';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from '../services/shared.service';
import {AuthService} from '../services/auth.service';
import {Subscription} from 'rxjs/Rx';
import {HttpClient} from '@angular/common/http';

@Component({
  selector: 'app-setup-otp',
  templateUrl: './setup-otp.component.html',
  styleUrls: ['./setup-otp.component.scss']
})
export class SetupOtpComponent implements OnInit {

  public imageToShow: any;
  public image: any;
  public token: string;
  public tokenSet: boolean = false;
  public password: string;
  public user: any;
  private authSubscription: Subscription;
  public code: string;

  constructor(private httpService: HttpService, private sharedService: SharedService, private authService: AuthService, private httpClient: HttpClient) { }

  ngOnInit() {
    if (this.authService.getUser()) {
      this.user = this.authService.getUser();
      this.checkToken();
    }
    this.authSubscription = this.authService.userInfo.subscribe(user => {
      this.user = user;
      this.checkToken();
    });
  }

  private clearForm(): void {
    this.token = '';
    this.password = '';
  }

  copyToClipboard(element: any) {
    element.select();
    document.execCommand('copy');
  }

  private checkToken(): void {
    if (!this.user.useOtp) {
      this.httpService.getFromUrl(environment.service + ApiRoutes.USER.SetupOTP).subscribe(resp =>{
        this.dataURItoBlob(resp.Image);
        this.createImageFromBlob(this.image);
        this.code = resp.Code;
      }, error => {
        if (error.indexOf('2fa already set') > 0) {

        } else {
          console.log(error);
        }
      });
    } else {
      this.tokenSet = true;
    }
  }

  dataURItoBlob(dataURI) {
    const byteString = window.atob(dataURI);
    const arrayBuffer = new ArrayBuffer(byteString.length);
    const int8Array = new Uint8Array(arrayBuffer);
    for (let i = 0; i < byteString.length; i++) {
      int8Array[i] = byteString.charCodeAt(i);
    }
    const blob = new Blob([int8Array], { type: 'image/jpeg' });
    this.image = blob;
  }

  private createImageFromBlob(image: Blob) {
    let reader = new FileReader();
    reader.addEventListener('load', () => {
      this.imageToShow = reader.result;
    }, false);

    if (image) {
      reader.readAsDataURL(image);
    }
  }

  private saveOtp() {
    return this.httpService.getFromUrl(environment.service + ApiRoutes.USER.ConfirmOTP + this.token).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Otp set for user');
        this.authService.refreshUserData();
        this.clearForm();
      },
      (err: any) => {
        this.sharedService.showError('Can\'t set otp for user user', err.split(': ')[1]);
        this.token = '';
      }
    );
  }

  private disableOtp() {
    const data = {
      token: this.token.toString(),
      password: this.password
    };
    return this.httpService.postToUrl(environment.service + ApiRoutes.USER.DisableOTP, data).subscribe(
      (data: any) => {
        this.sharedService.showSuccess('Disabled two factor authentication for user');
        this.authService.refreshUserData();
        this.tokenSet = false;
        this.clearForm();
      },
      (err: any) => {
        this.sharedService.showError('Can\'t disable two factor authentication for user user', '');
      }
    );
  }

}
