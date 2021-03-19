import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router} from '@angular/router';

import {AuthService} from './auth.service';

@Injectable()
export class RoleGuard implements CanActivate {

  constructor(private authService: AuthService, private router: Router) { }

  // ADMIN ROLES:
  // flag_vip, can_create, can_disable, review_whitelist, is_admin
  canActivate(route: ActivatedRouteSnapshot) {
    const expectedRole = route.data.expectedRole;
    if (this.authService.isAuth()) {
      switch (expectedRole) {
        case 'is_admin': {
         return this.authService.isAdmin();
        }
        case 'can_create': {
          return this.authService.canCreate();
        }
        case 'can_disable': {
          return this.authService.canDisable();
        }
        case 'can_manipulate_users': {
          return this.authService.canDisable() || this.authService.canCreate();
        }
      }
      return true;
    }

    this.router.navigate(['/']);
    return false;
  }

}
