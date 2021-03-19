import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import { FormGroup, Validators, FormBuilder, FormControl } from '@angular/forms';
import {HttpService} from '../services/http.service';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {SharedService} from '../services/shared.service';
import {AuthService} from '../services/auth.service';
import {User} from '../models/user.model';
import {UserService} from '../services/user.service';
import {MatDialog} from '@angular/material';
import {DeleteDialogComponent} from '../shared/dialogs/delete/delete.dialog.component';
import {Subscription} from 'rxjs/Rx';

@Component({
  selector: 'app-account-info',
  templateUrl: './account-info.component.html',
  styleUrls: ['./account-info.component.scss']
})
export class AccountInfoComponent implements OnInit {
  private user: User;
  private authSubscription: Subscription;
  updateAddress: FormGroup;
  displaySkycoinAddressForm = false;

  constructor(private userService: UserService, private authService: AuthService) { }

  ngOnInit() {
    if (this.authService.getUser()) {
      this.user = this.authService.getUser();
    }
    this.authSubscription = this.authService.userInfo.subscribe(user => {
      this.user = user;
    });
    this.updateAddress = new FormGroup({
      skycoinAddress: new FormControl()
    });
  }

  onSubmit() {
  // this.router.navigate(['update-address']);
  this.displaySkycoinAddressForm = true;
  }

  public resendMail(): void {
    this.userService.resendVerificationLink(this.user.username);
  }

  public getUser(): User {
    return this.user;
  }
}
