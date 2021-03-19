import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '../../../../../node_modules/@angular/router';
import { HttpService } from '../../../services/http.service';
import { SharedService } from '../../../services/shared.service';
import { environment } from '../../../../environments/environment';
import { ApiRoutes } from '../../routes';
import { User, Right } from '../../../models/user.model';

@Component({
  selector: 'app-edit-admin',
  templateUrl: './edit-admin.component.html',
  styleUrls: ['./edit-admin.component.scss']
})
export class EditAdminComponent implements OnInit {

  public adminEmail: number;
  public admin: User;
  public rights: string[];
  public isDataAvailable: boolean = false;

  constructor(public activeRoute: ActivatedRoute, 
              public httpService: HttpService,
              public sharedService: SharedService) { }

  ngOnInit() {
    this.activeRoute.queryParams.subscribe(params => {
      this.adminEmail = params['email'];
      this.httpService.getFromUrl(environment.service + ApiRoutes.ADMIN.Users + '/' + this.adminEmail).subscribe(
        (data: User) => {
          this.admin = data;
          this.rights = this.userRights();
          this.isDataAvailable = true;
        },
        (err: any) => {
          this.sharedService.showError('Can\'t load admin data from server: ', err.split(': ')[1]);
        }
      );
    });
  }

  userRights() {
    let activeOnes =this.admin.rights.filter(function(right){
      return right.Value === true;
    })
    return activeOnes.map(function (right) {
      return right.Name
    })
  }

  compareWithFunc(a, b) {
    return a === b;
  }

  saveEdit() {
    for (let right of this.admin.rights) {
      right.Value = this.rights.indexOf(right.Name) !== -1;
    }
    this.httpService.postToUrl(environment.service + ApiRoutes.ADMIN.AdminList + '/' + this.adminEmail + '/rights', this.admin).subscribe(res => {
      this.sharedService.showSuccess('Access rights successfully updated. If the admin is logged in he needs to log out and log back in for them to have the effect.');
    },
      (err: any) => {
        this.sharedService.showError('Can\'t persist admin data on the server: ', err.split(': ')[1]);
      }
    );
  }

  onCancel() {
    window.close();
  }

}
