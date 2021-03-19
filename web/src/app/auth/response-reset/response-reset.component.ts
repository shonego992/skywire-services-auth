import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';

import {PasswordValidator} from '../../shared/validators/password.validator';
import {FormBuilder, FormControl, FormGroup, Validators} from '@angular/forms';
import {UserService} from '../../services/user.service';

@Component({
    selector: 'app-response-reset',
    templateUrl: './response-reset.component.html',
    styleUrls: ['./response-reset.component.scss']
})
export class ResponseResetComponent implements OnInit {
    responseReset: FormGroup

    public form = {
        email: null,
        password: null,
        confirmPassword: null,
        resetToken: null
    };

    constructor(private activeRoute: ActivatedRoute, private userService: UserService, private formBuilder: FormBuilder) {
    }

    ngOnInit() {
        let frm = this.form;
        this.activeRoute.queryParams
            .subscribe(params => {
                frm.email = params['email'];
                frm.resetToken = params['token'];
            });

        this.responseReset = this.formBuilder.group({
            password: new FormControl('', { validators: [Validators.required, Validators.minLength(8)] }),
            confirmPassword: new FormControl('', { validators: [Validators.required] })
        }, {
                validator: PasswordValidator.validate.bind(this)
            });
    }

    onSubmit() {
        const val = this.responseReset.value
        if (val && val.password) {
            this.userService.changePassword(this.form.email, this.form.resetToken, val.password)
        }
    }

}
