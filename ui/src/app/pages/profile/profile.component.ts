import { Component, OnInit } from '@angular/core';

import { UsersService } from 'app/common/services/users/users.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
})
export class ProfileComponent implements OnInit {
  picture: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  company: string;
  department: string;
  occupation: string;
  location: string;

  newPassword: string = '';
  confirmPassword: string = '';
  oldPassword: string = '';
  ngxAdminMinPasswordSize = 6;

  constructor(
    private usersService: UsersService,
    private notificationsService: NotificationsService,
  ) { }

  ngOnInit(): void {
    this.picture = this.usersService.getUserPicture();

    this.usersService.getProfile().subscribe(
      resp => {
        this.email = resp.email ? resp.email : '';

        if (resp.metadata !== undefined) {
          this.firstName = resp.metadata.firstName ? resp.metadata.firstName : '';
          this.lastName = resp.metadata.lastName ? resp.metadata.lastName : '';
          this.phone = resp.metadata.phone ? resp.metadata.phone : '';
          this.company = resp.metadata.company ? resp.metadata.company : '';
          this.department = resp.metadata.department ? resp.metadata.department : '';
          this.occupation = resp.metadata.occupation ? resp.metadata.occupation : '';
          this.location = resp.metadata.location ? resp.metadata.location : '';
        }
      },
    );
  }

  onClickSaveInfos(event): void {
    const userReq = {
      metadata: {
        firstName: this.firstName,
        lastName: this.lastName,
        phone: this.phone,
        department: this.department,
        occupation: this.occupation,
        location: this.location,
        company: this.company,
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
