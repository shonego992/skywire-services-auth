import {Injectable} from '@angular/core';
import {CanActivate, Router} from '@angular/router';

import {AuthService} from './auth.service';

@Injectable()
export class LoginGuard implements CanActivate {

  constructor(private router: Router, private authService: AuthService) { }

  canActivate() {
    if (this.authService.isAuth()) {
      if (this.authService.isAdmin()) {
        this.router.navigate(['/admins-list']);
        return false;
      }
      this.router.navigate(['/account-info']);
      return false;
    }
    return true;
  }

}
