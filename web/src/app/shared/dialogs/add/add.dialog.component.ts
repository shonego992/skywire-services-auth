import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {Component, Inject, OnInit} from '@angular/core';
import {FormControl, Validators} from '@angular/forms';
import { User } from 'src/app/models/user.model';
import { UserDataService } from '../../../services/userData.service';

import { DOCUMENT } from '@angular/common';


@Component({
  selector: 'app-add.dialog',
  templateUrl: './add.dialog.component.html',
  styleUrls: ['./add.component.scss']
})
export class AddDialogComponent implements OnInit {
  constructor(public dialogRef: MatDialogRef<AddDialogComponent>,
                @Inject(MAT_DIALOG_DATA) public data: User,
                @Inject(DOCUMENT) private document: any, 
                public dataService: UserDataService) { }

    formControl = new FormControl('', [
      Validators.required,
      Validators.email,
    ]);

    getErrorMessage() {
      return this.formControl.hasError('required') ? 'Required field' :
        this.formControl.hasError('email') ? 'Not a valid email' :
          '';
    }

    submit() {
    // emppty stuff
    }

    onNoClick(): void {
      this.dialogRef.close();
    }

    public confirmAdd(): void {
      var _self = this;
      this.dataService.addAdmin(this.data, () => { _self.document.location.reload() });
    }

    ngOnInit() {}
}
