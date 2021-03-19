import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {SignupComponent} from './auth/signup/signup.component';
import {LoginComponent} from './auth/login/login.component';
import {UserProfileComponent} from './user-profile/user-profile.component';
import {AuthGuard} from './services/auth.guard';
import {RoleGuard} from './services/role.guard';
import {AdminProfileComponent} from './admin-profile/admin-profile.component';
import {ForgotPasswordComponent} from './auth/forgot-password/forgot-password.component';
import {ResponseResetComponent} from './auth/response-reset/response-reset.component';
import {VerifyUserComponent} from './auth/verify-user/verify-user.component';
import {AccountInfoComponent} from './account-info/account-info.component';
import {ChangingPasswordComponent} from './user-profile/changing-password/changing-password.component';
import {UpdateAddressComponent} from './user-profile/update-address/update-address.component';
import {AdminsListComponent} from './admin-profile/admins-list/admins-list.component';
import {UsersListComponent} from './admin-profile/users-list/users-list.component';
import {EditAdminComponent} from './shared/dialogs/edit-admin/edit-admin.component';
import {EditUserComponent} from './shared/dialogs/edit-user/edit-user.component';
import { SuccessComponent } from './user-profile/success/success.component';
import {LoginGuard} from './services/login.guard';
import {PageNotFoundComponent} from './page-not-found/page-not-found.component';
import {SetupOtpComponent} from './setup-otp/setup-otp.component';
import { LayoutComponent } from './layout/layout.component';

const routes: Routes = [
  {path: '', component: LoginComponent, canActivate: [LoginGuard]},
  {path: 'signup', component: SignupComponent},
  {path: 'reset-password', component: ForgotPasswordComponent},
  {path: 'change-password', component: ResponseResetComponent },
  {path: 'verify-profile', component: VerifyUserComponent },
  {
    path: '',
    component: LayoutComponent,
    children: [
      {path: 'success-message', component: SuccessComponent},
      {path: 'user', component: UserProfileComponent, canActivate: [AuthGuard]},
      {path: 'admin', component: AdminProfileComponent, canActivate: [RoleGuard], data: {expectedRole: 'is_admin'}},
      {path: 'account-info', component: AccountInfoComponent, canActivate: [AuthGuard]},
      {path: 'changing-password', component: ChangingPasswordComponent, canActivate: [AuthGuard]},
      {path: 'update-address', component: UpdateAddressComponent, canActivate: [AuthGuard]},
      {path: 'admins-list', component: AdminsListComponent, canActivate: [RoleGuard], data: { expectedRole: 'can_manipulate_users'}},
      {path: 'users-list', component: UsersListComponent, canActivate: [RoleGuard], data: { expectedRole: 'can_manipulate_users' }},
      {path: 'edit-admin', component: EditAdminComponent, canActivate: [RoleGuard], data: { expectedRole: 'can_create'}},
      {path: 'edit-user', component: EditUserComponent, canActivate: [RoleGuard], data: { expectedRole: 'can_create'}},
      {path: 'setup-otp', component: SetupOtpComponent, canActivate: [AuthGuard]},
    ],
  },
  {path: '**', component: PageNotFoundComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  providers: [AuthGuard, RoleGuard, LoginGuard]
})
export class AppRoutingModule { }
