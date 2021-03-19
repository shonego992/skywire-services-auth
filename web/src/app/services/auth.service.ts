import {Injectable, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {Subject} from 'rxjs/Subject';
import 'rxjs/add/operator/map';
import {environment} from '../../environments/environment';
import {HttpService} from './http.service';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from './shared.service';
import * as moment from 'moment';
import {User} from '../models/user.model';
import decode from 'jwt-decode';
import {AdminClaims} from '../models/admin.claims';

const TOKEN_KEY = 'token_key';
const TOKEN_EXPIRE = 'token_expire';

@Injectable()
export class AuthService implements OnInit {

  private userVerified = false;
  private user: User;
  private adminClaims: AdminClaims;

  authChange = new Subject<any>();
  userInfo = new Subject<User>();

  ngOnInit() {
  }

  public refreshUserData(): void {
      this.httpService.getFromUrl(environment.service + '/auth/info').subscribe((user: User) => {
      this.setUser(user);
    },
    (err: any) => {
      this.sharedService.showError('Can\'t get user details', err.split(': ')[1]);
    }
    );
  }

  public getAdminClaims(): AdminClaims {
    return this.adminClaims;
  }

  public getUser(): User {
    return this.user;
  }

  public isAdmin(): boolean {
    this.checkClaims();
    return this.canDisable() || this.canCreate();
  }

  public canCreate(): boolean {
    this.checkClaims();
    return this.adminClaims ? this.adminClaims.can_create: false;
  }

  public canDisable(): boolean {
    this.checkClaims();
    return this.adminClaims ? this.adminClaims.can_disable: false;
  }

  private checkClaims(): void {
    if (!this.adminClaims) {
      this.saveTokenClaims(localStorage.getItem(TOKEN_KEY));
    }
  }

  // TODO: need to set user as an observable and follow it where it is needed
  public setUser(user: User) {
    this.user = user;
    this.userVerified = user.status === 1;
    this.saveTokenClaims(localStorage.getItem(TOKEN_KEY));
    this.userInfo.next(this.user);
  }

  public isUserVerified() {
    return this.userVerified;
  }

  public setUserVefified(value: boolean) {
   this.userVerified = value;
  }

  constructor (private router: Router, private httpService: HttpService, private sharedService: SharedService) {}

  public getToken (): string {
    return localStorage.getItem(TOKEN_KEY);
  }

  // Save token data into local storage
  public saveToken (data: any) {
    const expiresAt = moment().add(data.expire, 'second');
    localStorage.removeItem(TOKEN_KEY);
    localStorage.setItem(TOKEN_KEY, data.token);

    localStorage.removeItem(TOKEN_EXPIRE);
    localStorage.setItem(TOKEN_EXPIRE, JSON.stringify(expiresAt.valueOf()));
    this.saveTokenClaims(data.token);
  }

  private saveTokenClaims(token: any): void {
    const tokenPayload: AdminClaims = decode(token);
    this.adminClaims = tokenPayload;
    this.authChange.next({isAuth: true, claims: this.adminClaims});
  }

  private getTokenExpirationDate (): any {
    const expiration =  localStorage.getItem(TOKEN_EXPIRE);
    const expiresAt = JSON.parse(expiration);
    return moment(expiresAt);
  }

  private isTokenExpired (): boolean {
    const isExpired = moment().isBefore(this.getTokenExpirationDate());
    return isExpired;
  }

  // TODO: check this method and is auth.. for now the same, but include rolse inside?
  public isLoggedin () {
    return localStorage.getItem(TOKEN_KEY) && !this.isTokenExpired();
  }

  public logout (): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(TOKEN_EXPIRE);
    localStorage.clear();

    this.adminClaims = new AdminClaims;
    this.authChange.next({isAuth: false, claims: this.adminClaims});
    this.router.navigate(['/']);
  }

  private refreshToken () {
    return this.httpService.getFromUrl(environment.service + ApiRoutes.AUTH.Refresh + localStorage.getItem(TOKEN_KEY))
      .map((response: any) => {
        this.saveToken(response);
      });
  }

  getAccount () {
    return this.httpService.getFromUrl(environment.service + '/account');
  }

  getUserFromUrl () {
    return this.httpService.getFromUrl(environment.service + '/auth/info');
    // let payload: any = decode(localStorage.getItem(TOKEN_KEY));
    // console.log(payload);
  }

  public isAuth (): boolean {
    const isAuth: boolean = localStorage.getItem(TOKEN_KEY) && !this.isTokenExpired();
    return isAuth;
  }


  public authSuccessfully () {
    this.authChange.next({isAuth: true, claims: this.adminClaims});
    if (this.isAdmin()) {
      this.router.navigate(['/admins-list']);
      return;
    }
    this.router.navigate(['/account-info']);
  }
}
