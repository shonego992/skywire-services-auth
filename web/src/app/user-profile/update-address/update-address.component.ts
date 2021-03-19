import { Component, OnInit } from '@angular/core';
import { FormGroup, Validators, FormBuilder } from '@angular/forms';
import { UserService } from '../../services/user.service';

@Component({
  selector: 'app-update-address',
  templateUrl: './update-address.component.html',
  styleUrls: ['./update-address.component.scss']
})
export class UpdateAddressComponent implements OnInit {
  updateAddress: FormGroup;

  constructor(private userService: UserService, private formBuilder: FormBuilder) {}


  ngOnInit() {
    this.updateAddress = this.formBuilder.group({
      skycoinAddress: ['', { validators: [Validators.required] }]
    });
  }

  onSubmit() {
    const val = this.updateAddress.value;
    if (val.skycoinAddress) {
      this.userService.updateAddress(val.skycoinAddress);
    }
  }

}
