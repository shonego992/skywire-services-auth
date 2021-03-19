import { FormGroup } from '@angular/forms';

export class PasswordValidator {

   static validate(passwordFormGroup: FormGroup) {
       let password = passwordFormGroup.controls.password.value;
       let confirmPassword = passwordFormGroup.controls.confirmPassword.value;

       if (confirmPassword.length <= 0) {
           return null;
       }

       if (confirmPassword !== password) {
           return {
               doesMatchPassword: true
           };
       }

       return null;

   }
}
