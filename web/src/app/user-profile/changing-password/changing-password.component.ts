import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormGroup, Validators, FormBuilder } from '@angular/forms';
import { UserService } from '../../services/user.service';
import { PasswordValidator } from '../../shared/validators/password.validator';

@Component({
  selector: 'app-changing-password',
  templateUrl: './changing-password.component.html',
  styleUrls: ['./changing-password.component.scss']
})
export class ChangingPasswordComponent implements OnInit {
  passwordChange: FormGroup;
  passwordFormGroup: FormGroup;

  constructor(private userService: UserService, private formBuilder: FormBuilder, private router: Router) {}

  ngOnInit() {
    this.passwordChange = this.formBuilder.group({
      password: ['', [Validators.required, Validators.minLength(8)]],
      confirmPassword: ['', [Validators.required, Validators.minLength(8)]]
    }, {
        validator: PasswordValidator.validate.bind(this)
      });

    this.passwordFormGroup = this.formBuilder.group({
      oldPassword: ['', Validators.required],
      passwordChange: this.passwordChange
    });
  }

  onSubmit() {
    const val = this.passwordFormGroup.value;
    const val2 = this.passwordChange.value
    if (val.oldPassword && val2.password) {
      this.userService.updatePassword(val.oldPassword, val2.password);
    } 
  }

}
