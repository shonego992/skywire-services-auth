import {Component, OnInit, Inject} from '@angular/core';
import { FormControl, FormGroup, Validators, FormBuilder} from '@angular/forms';
import {User} from '../../models/user.model';
import { UserService } from '../../services/user.service';
import { DOCUMENT } from '@angular/common';
import { ActivatedRoute, Router } from '../../../../node_modules/@angular/router';
import { PasswordValidator } from '../../shared/validators/password.validator';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent implements OnInit {
  passwordMatchGroup: FormGroup;
  signupForm: FormGroup;
  public user: any;
  private redirectURL: string;

  ngOnInit() {
    this.activeRoute.queryParams.subscribe(params => {
      this.redirectURL = params['redirectURL'];
    });

    this.passwordMatchGroup = this.formBuilder.group({
      password: ['', [Validators.required, Validators.minLength(8)]],
      confirmPassword: ['', [Validators.required, Validators.minLength(8)]]
    }, {
        validator: PasswordValidator.validate.bind(this)
      });

    this.signupForm = this.formBuilder.group({
      username: ['', Validators.required],
      passwordMatchGroup: this.passwordMatchGroup
    });

  }

  constructor(public activeRoute: ActivatedRoute,
    @Inject(DOCUMENT) private document: any,
    private router: Router,
    private formBuilder: FormBuilder,
    private userService: UserService) { }

  resetForm(signupForm: FormGroup) {
    if (signupForm != null){
      signupForm.reset();
      this.user = {
        username: '',
        password: ''
      };
    }
  }

  onSubmit() {
    let _self = this;
    this.userService.registerUser({
      username: this.signupForm.value.username,
      password: this.passwordMatchGroup.value.password,
   }, function(success) {
      if (!success) return
      if (_self.redirectURL) {
        _self.document.location.href = _self.redirectURL;
      } else {
        _self.router.navigate(['/']);
      }
    });
  }
}
