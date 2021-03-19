import {Component, OnInit, ViewChild} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {MatDialog, MatPaginator} from '@angular/material';
import {Router} from '@angular/router';
import {UserDataService} from '../services/userData.service';

@Component({
  selector: 'app-user-miner-overview',
  templateUrl: './user-miner-overview.component.html',
  styleUrls: ['./user-miner-overview.component.scss']
})
export class UserMinerOverviewComponent implements OnInit {

//   displayedColumns = ['username', 'actions']; //TODO add maybe 'status' 'access rights' or 'last_sign_in'
//   dataSource: ExampleDataSource | null;
//   index: number;
//   id: number;
//
//   constructor(public httpClient: HttpClient,
//               public dialog: MatDialog,
//               public dataService: UserDataService,
//               public router: Router) {}
//
//   @ViewChild(MatPaginator) paginator: MatPaginator;
//   @ViewChild(MatSort) sort: MatSort;
//   @ViewChild('filter') filter: ElementRef;
//
  ngOnInit() {
    // this.loadData();
  }
//
//   refresh() {
//     this.loadData();
//   }
//
//   addNew(issue: User) {
//     const dialogRef = this.dialog.open(AddDialogComponent, {
//       data: {issue: issue }
//     });
//
//     dialogRef.afterClosed().subscribe(result => {
//       if (result === 1) {
//         // After dialog is closed we're doing frontend updates
//         // For add we're just pushing a new row inside DataService
//         this.dataService.dataChange.value.push(this.dataService.getDialogData());
//         this.refreshTable();
//       }
//     });
//   }
//
//   startEdit(i: number,state: string, url: string) {
//     // index row is used just for debugging proposes and can be removed
//     this.index = i;
//     this.router.navigate(['edit-admin']);
//   }
//
//   deleteAdmin(i: number, id: number, username: string) {
//     this.index = i;
//     this.id = id;
//     const dialogRef = this.dialog.open(DeleteDialogComponent, {
//       data: {id: id, state: '', url: ''}
//     });
//
//     dialogRef.afterClosed().subscribe(result => {
//       if (result === 1) {
//         this.dataService.deleteUser(username);
//         const foundIndex = this.dataService.dataChange.value.findIndex(x => x.id === this.id);
//         // for delete we use splice in order to remove single object from DataService
//         this.dataService.dataChange.value.splice(foundIndex, 1);
//         this.refreshTable();
//       }
//     });
//   }
//
//
//   // If you don't need a filter or a pagination this can be simplified, you just use code from else block
//   private refreshTable() {
//     // if there's a paginator active we're using it for refresh
//     if (this.dataSource._paginator.hasNextPage()) {
//       this.dataSource._paginator.nextPage();
//       this.dataSource._paginator.previousPage();
//       // in case we're on last page this if will tick
//     } else if (this.dataSource._paginator.hasPreviousPage()) {
//       this.dataSource._paginator.previousPage();
//       this.dataSource._paginator.nextPage();
//       // in all other cases including active filter we do it like this
//     } else {
//       this.dataSource.filter = '';
//       this.dataSource.filter = this.filter.nativeElement.value;
//     }
//   }
//
//   public loadData() {
//     this.dataService = new UserDataService(this.httpClient);
//     this.dataSource = new ExampleDataSource(this.dataService, this.paginator, this.sort);
//     Observable.fromEvent(this.filter.nativeElement, 'keyup')
//       .debounceTime(150)
//       .distinctUntilChanged()
//       .subscribe(() => {
//         if (!this.dataSource) {
//           return;
//         }
//         this.dataSource.filter = this.filter.nativeElement.value;
//       });
//   }
// }
//
//
// export class ExampleDataSource extends DataSource<User> {
//   _filterChange = new BehaviorSubject('');
//
//   get filter(): string {
//     return this._filterChange.value;
//   }
//
//   set filter(filter: string) {
//     this._filterChange.next(filter);
//   }
//
//   filteredData: User[] = [];
//   renderedData: User[] = [];
//
//   constructor(public _exampleDatabase: UserDataService,
//               public _paginator: MatPaginator,
//               public _sort: MatSort) {
//     super();
//     // Reset to the first page when the user changes the filter.
//     this._filterChange.subscribe(() => this._paginator.pageIndex = 0);
//   }
//
//   /** Connect function called by the table to retrieve one stream containing the data to render. */
//   connect(): Observable<User[]> {
//     // Listen for any changes in the base data, sorting, filtering, or pagination
//     const displayDataChanges = [
//       this._exampleDatabase.dataChange,
//       this._sort.sortChange,
//       this._filterChange,
//       this._paginator.page
//     ];
//
//     this._exampleDatabase.getAllAdmins();
//
//     return Observable.merge(...displayDataChanges).map(() => {
//       // Filter data
//       this.filteredData = this._exampleDatabase.data.slice().filter((user: User) => {
//         const searchStr = (user.id + user.username).toLowerCase();
//         return searchStr.indexOf(this.filter.toLowerCase()) !== -1;
//       });
//
//       // Sort filtered data
//       const sortedData = this.sortData(this.filteredData.slice());
//
//       // Grab the page's slice of the filtered sorted data.
//       const startIndex = this._paginator.pageIndex * this._paginator.pageSize;
//       this.renderedData = sortedData.splice(startIndex, this._paginator.pageSize);
//       return this.renderedData;
//     });
//   }
//   disconnect() {
//   }
//
//
//
//   /** Returns a sorted copy of the database data. */
//   sortData(data: User[]): User[] {
//     if (!this._sort.active || this._sort.direction === '') {
//       return data;
//     }
//
//     return data.sort((a, b) => {
//       let propertyA: number | string = '';
//       let propertyB: number | string = '';
//
//       switch (this._sort.active) {
//         case 'id': [propertyA, propertyB] = [a.id, b.id]; break;
//         case 'username': [propertyA, propertyB] = [a.username, b.username]; break;
//         case 'status': [propertyA, propertyB] = [a.status, b.status]; break;
//       }
//
//       const valueA = isNaN(+propertyA) ? propertyA : +propertyA;
//       const valueB = isNaN(+propertyB) ? propertyB : +propertyB;
//
//       return (valueA < valueB ? -1 : 1) * (this._sort.direction === 'asc' ? 1 : -1);
//     });
//   }
// }
}
