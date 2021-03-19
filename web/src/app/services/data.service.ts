import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Issue } from '../models/issue';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {environment} from '../../environments/environment';
import {ApiRoutes} from '../shared/routes';
import {WhitelistApplication} from '../models/application-model';
import {SharedService} from './shared.service';

@Injectable()
export class DataService {
  private readonly API_URL = environment.service + ApiRoutes.ADMIN.Whitelists;

  dataChange: BehaviorSubject<WhitelistApplication[]> = new BehaviorSubject<WhitelistApplication[]>([]);
  // Temporarily stores data from dialogs
  dialogData: any;

  constructor (private httpClient: HttpClient, private sharedService: SharedService) {}

  get data(): WhitelistApplication[] {
    return this.dataChange.value;
  }

  getDialogData() {
    return this.dialogData;
  }

  /** CRUD METHODS */
  getAllIssues(): void {
    this.httpClient.get<WhitelistApplication[]>(this.API_URL).subscribe(data => {
      for (let item of data) {
        item.currentStatus = this.sharedService.mapStatus(item.currentStatus);
      }

        this.dataChange.next(data);
      },
      (error: HttpErrorResponse) => {
      console.log (error.name + ' ' + error.message);
      });
  }

  addIssue (issue: Issue): void {
    this.dialogData = issue;
  }

  updateIssue (issue: Issue): void {
    this.dialogData = issue;
  }

  deleteIssue (id: number): void {
    console.log(id);
  }
}
/* REAL LIFE CRUD Methods with ToasterService uses Material Toasts for displaying messages:
    // ADD, POST METHOD
    addItem(kanbanItem: KanbanItem): void {
    this.httpClient.post(this.API_URL, kanbanItem).subscribe(data => {
      this.dialogData = kanbanItem;
      this.toasterService.showToaster('Successfully added', 3000);
      },
      (err: HttpErrorResponse) => {
      this.toasterService.showToaster('Error occurred. Details: ' + err.name + ' ' + err.message, 8000);
    });
   }
    // UPDATE, PUT METHOD
     updateItem(kanbanItem: KanbanItem): void {
    this.httpClient.put(this.API_URL + kanbanItem.id, kanbanItem).subscribe(data => {
        this.dialogData = kanbanItem;
        this.toasterService.showToaster('Successfully edited', 3000);
      },
      (err: HttpErrorResponse) => {
        this.toasterService.showToaster('Error occurred. Details: ' + err.name + ' ' + err.message, 8000);
      }
    );
  }
  // DELETE METHOD
  deleteItem(id: number): void {
    this.httpClient.delete(this.API_URL + id).subscribe(data => {
      console.log(data['']);
        this.toasterService.showToaster('Successfully deleted', 3000);
      },
      (err: HttpErrorResponse) => {
        this.toasterService.showToaster('Error occurred. Details: ' + err.name + ' ' + err.message, 8000);
      }
    );
  }
*/
