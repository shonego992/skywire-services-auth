import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import { User } from '../models/user.model';

@Injectable()
export class UserDataService {
  private readonly API_URL_ADMINS = environment.service + ApiRoutes.ADMIN.AdminList;
  private readonly API_URL_USERS = environment.service + ApiRoutes.ADMIN.Users;
  private readonly API_URL_ADMIN = environment.service + ApiRoutes.ADMIN.Admin;

  dataChange: BehaviorSubject<User[]> = new BehaviorSubject<User[]>([]);
  // Temporarily stores data from dialogs
  dialogData: any;

  constructor (private httpClient: HttpClient) {}

  get data(): User[] {
    return this.dataChange.value;
  }

  getDialogData() {
    return this.dialogData;
  }

  /** CRUD METHODS */
  getAllUsers(): void {
    this.httpClient.get<User[]>(this.API_URL_USERS).subscribe(data => {
        this.dataChange.next(data);
      },
      (error: HttpErrorResponse) => {
      console.log (error.name + ' ' + error.message);
      });
  }

  addUser (issue: User): void {
    this.dialogData = issue;
  }

  updateUser (issue: User): void {
    this.dialogData = issue;
  }

  deleteUser (username: string, callback: any): void {
    this.httpClient.delete<User>(this.API_URL_USERS + '/' + username).subscribe(
      data => {
        callback();
      },
      (error: HttpErrorResponse) => {
        console.log(error.name + ' ' + error.message);
      });
  }

  activateUser(username: string, callback: any): void {
    this.httpClient.get<User>(this.API_URL_USERS + '/' + username + '/activate').subscribe(
      data => {
        callback();
      },
      (error: HttpErrorResponse) => {
        console.log(error.name + ' ' + error.message);
      });
  }

  getAllAdmins(): void {
    this.httpClient.get<User[]>(this.API_URL_ADMINS).subscribe(data => {
      this.dataChange.next(data);
    },
      (error: HttpErrorResponse) => {
        console.log(error.name + ' ' + error.message);
      });
  }

  addAdmin(newAdmin: User, callback: any): void {
    this.httpClient.post<User[]>(this.API_URL_ADMIN, newAdmin).subscribe(response => {
      this.dataChange.next(response);
      callback();
    },
      (error: HttpErrorResponse) => {
        console.log(error.name + ' ' + error.message);
    });
  }

  updateAdmin(issue: User): void {
    this.dialogData = issue;
  }

}