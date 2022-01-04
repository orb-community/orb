import { Component, OnInit } from '@angular/core';

import { UsersService } from 'app/common/services/users/users.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
})
export class ProfileComponent implements OnInit {
  user: any = {};

  newPassword: string = '';
  confirmPassword: string = '';
  oldPassword: string = '';
  ngxAdminMinPasswordSize = 6;

  constructor(
    private usersService: UsersService,
    private notificationsService: NotificationsService,
  ) { }

  ngOnInit(): void {
    this.usersService.getProfile().subscribe(
      resp => {
        this.user.email = resp.email ? resp.email : '';

        if (resp.metadata !== undefined) {
          this.user.fullName = resp.metadata.fullName ? resp.metadata.fullName : '';
          this.user.company = resp.metadata.company ? resp.metadata.company : '';
        }
      },
    );
  }

  onClickSaveInfos(event): void {
    const userReq = {
      metadata: {
        fullName: this.user.fullName,
        company: this.user.company,
      },
    };

    this.usersService.editUser(userReq).subscribe(
      resp => {
        this.notificationsService.success('User successfully edited', '');
      },
    );
  }

  onClickSavePassword(event): void {
    if (this.newPassword.length < this.ngxAdminMinPasswordSize) {
      this.notificationsService.warn('Password must be at least 6 characters long.', '');
      return;
    }

    if (this.newPassword === this.confirmPassword) {
      const passReq = {
        password: this.newPassword,
        old_password: this.oldPassword,
      };

      this.usersService.changeUserPassword(passReq).subscribe(
        resp => {
          this.notificationsService.success('Password successfully changed', '');
        },
      );
    } else {
      this.notificationsService.warn('New password and Confirmation password do not match.', '');
    }
  }
}
