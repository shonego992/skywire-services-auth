import {Component, OnInit, Inject} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {UserService} from '../../services/user.service';
import { DOCUMENT } from '@angular/common';
import { ActivatedRoute, Router } from '../../../../node_modules/@angular/router';
import { EmailNormalization } from 'src/app/shared/validators/email.normalization';

@Component({
  selector: 'app-forgot-password',
  templateUrl: './forgot-password.component.html',
  styleUrls: ['./forgot-password.component.scss']
})
export class ForgotPasswordComponent implements OnInit {
  resetPass: FormGroup;
  private redirectURL: string;

  constructor(public activeRoute: ActivatedRoute, 
    @Inject(DOCUMENT) private document: any,
    private router: Router,
    private userService: UserService,
    private emailNormalization: EmailNormalization) { }

    ngOnInit() {
      this.activeRoute.queryParams.subscribe(params => {
        this.redirectURL = params['redirectURL'];
      });
      this.resetPass = new FormGroup({
        email: new FormControl('', {validators: [Validators.required, Validators.email]})
      });
    }

  onSubmit() {
    let _self = this;
    const val = this.resetPass.value;
    let stringValue = val.email.toLowerCase().replace(' ', '');
    let normalizedEmail = this.emailNormalization.emailNormalization(stringValue);
    if (normalizedEmail) {
      this.userService.resetPass(normalizedEmail, function (success) {
        if (!success) return
        if (_self.redirectURL) {
          _self.document.location.href = _self.redirectURL;
        } else {
          this.authService.logout();
        }
      });
    }
  }

}
