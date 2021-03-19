import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {UserService} from '../../services/user.service';

@Component({
    selector: 'app-verify-user',
    templateUrl: './verify-user.component.html',
    styleUrls: ['./verify-user.component.scss']
})
export class VerifyUserComponent implements OnInit {
    constructor(private activeRoute: ActivatedRoute, private userService: UserService) {}

    ngOnInit() {
        // TODO init spinner or something
        let token = '';
        this.activeRoute.queryParams
            .subscribe(params => {
                token = params['token'];
            });
        this.userService.verifyProfile(token);
    }
}
