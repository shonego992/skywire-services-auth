import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {UserService} from '../../services/user.service';
import {Subscription} from 'rxjs/Rx';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  loginForm: FormGroup;
  tokenNeeded = false;
  private tokenSub: Subscription;

  constructor(private userService: UserService) { }

  ngOnInit() {
    this.loginForm = new FormGroup({
      username: new FormControl('', {validators: [Validators.required, this.emailValidator]}),
      password: new FormControl('', {validators: [Validators.required, Validators.minLength(8)]}),
      token: new FormControl()
    });
  }

  onSubmit(): void {
    const val = this.loginForm.value;
        if (val.username && val.password) {
            this.userService.login(val.username.trim(), val.password, val.token);
            this.tokenSub = this.userService.tokenNeeded.subscribe((value: any) => {
              this.tokenNeeded = value;
            });
        }
  }

  public emailValidator(control: FormControl) {
    const re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{1,}))$/;
    const valid =  re.test((String(control.value || '').trim()).toLowerCase());
    return valid ? null : { 'error': true };
  }
}
