import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FlexLayoutModule} from '@angular/flex-layout';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HTTP_INTERCEPTORS, HttpClient, HttpClientModule} from '@angular/common/http';
import {ToastrModule} from 'ngx-toastr';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';

import {AppComponent} from './app.component';
import {MaterialModule} from './material.module';
import {SignupComponent} from './auth/signup/signup.component';
import {LoginComponent} from './auth/login/login.component';
import {AppRoutingModule} from './app-routing.module';
import {HeaderComponent} from './navigation/header/header.component';
import {SidenavListComponent} from './navigation/sidenav-list/sidenav-list.component';
import {AuthService} from './services/auth.service';
import {UserService} from './services/user.service';
import {DataService} from './services/data.service';
import {UserDataService} from './services/userData.service';
import {AuthInterceptor} from './shared/interceptors/auth.interceptor';
import {UserProfileComponent} from './user-profile/user-profile.component';
import {ChangingPasswordComponent} from './user-profile/changing-password/changing-password.component';
import {UpdateAddressComponent} from './user-profile/update-address/update-address.component';
import {AdminProfileComponent} from './admin-profile/admin-profile.component';
import {HttpService} from './services/http.service';
import {AccountInfoComponent} from './account-info/account-info.component';
import {AdminsListComponent} from './admin-profile/admins-list/admins-list.component';
import {UsersListComponent} from './admin-profile/users-list/users-list.component';
import {CreateNewAdminComponent} from './admin-profile/create-new-admin/create-new-admin.component';
import {ForgotPasswordComponent} from './auth/forgot-password/forgot-password.component';
import {ResponseResetComponent} from './auth/response-reset/response-reset.component';
import {AddDialogComponent} from './shared/dialogs/add/add.dialog.component';
import {DeleteDialogComponent} from './shared/dialogs/delete/delete.dialog.component';
import {VerifyUserComponent} from './auth/verify-user/verify-user.component';
import {TranslateLoader, TranslateModule} from '@ngx-translate/core';
import {ErrorInterceptor} from './shared/interceptors/error.interceptor';
import {TranslateHttpLoader} from '@ngx-translate/http-loader';
import {SharedService} from './services/shared.service';
import {NgxUploaderModule} from 'ngx-uploader';
import {SuccessComponent} from './user-profile/success/success.component';
import {EditAdminComponent} from './shared/dialogs/edit-admin/edit-admin.component';
import {EditUserComponent} from './shared/dialogs/edit-user/edit-user.component';
import {EmailNormalization} from './shared/validators/email.normalization';
// TODO: consider to drop components below this line
import {AlertComponent} from './shared/dialogs/alert/alert.component';
import {NgxImageGalleryModule} from 'ngx-image-gallery';
import {UserMinerOverviewComponent} from './user-miner-overview/user-miner-overview.component';
import {DisableComponent} from './shared/dialogs/disable/disable.component';
import {PageNotFoundComponent} from './page-not-found/page-not-found.component';
import {SetupOtpComponent} from './setup-otp/setup-otp.component';
import {LayoutComponent} from './layout/layout.component';

@NgModule({
  declarations: [
    AppComponent,
    SignupComponent,
    LoginComponent,
    HeaderComponent,
    SidenavListComponent,
    UserProfileComponent,
    ChangingPasswordComponent,
    UpdateAddressComponent,
    AdminProfileComponent,
    AccountInfoComponent,
    AdminsListComponent,
    UsersListComponent,
    CreateNewAdminComponent,
    ForgotPasswordComponent,
    ResponseResetComponent,
    VerifyUserComponent,
    AddDialogComponent,
    DeleteDialogComponent,
    AlertComponent,
    EditAdminComponent,
    EditUserComponent,
    SuccessComponent,
    UserMinerOverviewComponent,
    DisableComponent,
    PageNotFoundComponent,
    SetupOtpComponent,
    LayoutComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    MaterialModule,
    AppRoutingModule,
    FlexLayoutModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    NgxUploaderModule,
    NgxImageGalleryModule,
    ToastrModule.forRoot({
      timeOut: 10000,
      positionClass: 'toast-top-center',
      preventDuplicates: true,
    }),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: (createTranslateLoader),
        deps: [HttpClient]
      }
    })
  ],
  entryComponents: [
    AddDialogComponent,
    DeleteDialogComponent,
    DisableComponent
  ],
  providers: [
    AuthService,
    UserService,
    HttpService,
    DataService,
    UserDataService,
    SharedService,
    EmailNormalization,
    { provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true
    },
    { provide: HTTP_INTERCEPTORS,
      useClass: ErrorInterceptor,
      multi: true
    },
    { provide: MAT_DIALOG_DATA,
      useValue: {}
    },
    { provide: MatDialogRef,
      useValue: {}
    },

  ],
  bootstrap: [AppComponent]
})
export class AppModule { }

export function createTranslateLoader(http: HttpClient) {
  return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}
