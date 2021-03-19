import { Component, OnInit, Inject } from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';

@Component({
  selector: 'app-disable',
  templateUrl: './disable.component.html',
  styleUrls: ['./disable.component.scss']
})
export class DisableComponent implements OnInit {
  constructor(public dialogRef: MatDialogRef<DisableComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any) { }

    onNotClick(): void {
      this.dialogRef.close();
    }

    confirmDisable(): void {
      return;
    }
    ngOnInit() {

    }
}
