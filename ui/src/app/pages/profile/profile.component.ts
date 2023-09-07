import { Component, OnInit } from '@angular/core';

import { UsersService } from 'app/common/services/users/users.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { User } from 'app/common/interfaces/mainflux.interface';
import { error } from 'console';

@Component({
  selector: 'ngx-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
})
export class ProfileComponent implements OnInit {
  user: any = {};
  userInfo: User;
  newPassword: string = '';
  confirmPassword: string = '';
  oldPassword: string = '';
  ngxAdminMinPasswordSize = 8;
  isLoading: boolean = true;

  userFullName: string;
  userCompany: string;

  oldPasswordInput: string;
  newPasswordInput: string;
  confirmPasswordInput: string;

  editMode = {
    work: false,
    profileName: false,
    password: false,
  }

  isPasswordValidSize: boolean;
  isPasswordValidMatch: boolean;
  isRequesting = false;


  constructor(
    private usersService: UsersService,
    private notificationsService: NotificationsService,
  ) { 
    this.oldPasswordInput = '';
    this.newPasswordInput = '';
    this.confirmPasswordInput = '';
  }

  ngOnInit(): void {
    this.retrieveUserInfo();
  }

  retrieveUserInfo(): void {
    this.isLoading = true;
    this.usersService.getProfile().subscribe(
      resp => {
        this.user.picture = this.usersService.getUserPicture();
        this.user.email = resp.email ? resp.email : '';

        if (resp.metadata !== undefined) {
          this.user.fullName = resp.metadata.fullName ? resp.metadata.fullName : '';
          this.user.company = resp.metadata.company ? resp.metadata.company : '';
          this.userFullName = this.user.fullName;
          this.userCompany =  this.user.company;

        }
        this.isLoading = false;
      },
    );
  }
  editUserDetails(fullName: string, company: string): void {
    this.isRequesting = true;
    const userReq = {
      metadata: {
        fullName: fullName,
        company: company,
      },
    };
  
    this.usersService.editUser(userReq).subscribe(
      resp => {
        this.notificationsService.success('User successfully edited', '');
        this.retrieveUserInfo();
        this.toggleEdit('');
        this.isRequesting = false;
      },
    );
  }
  
  canChangePassword(): boolean {
    this.isPasswordValidSize = this.newPasswordInput.length >= this.ngxAdminMinPasswordSize;
    this.isPasswordValidMatch = this.newPasswordInput === this.confirmPasswordInput;
    return this.isPasswordValidSize && this.isPasswordValidMatch;
  }

  changePassword(): void {
    this.isRequesting = true;
    const passReq = {
      password: this.newPasswordInput,
      old_password: this.oldPasswordInput,
    };

    this.usersService.changeUserPassword(passReq).subscribe(
      resp => {
        this.notificationsService.success('Password successfully changed', '');
        this.retrieveUserInfo();
        this.toggleEdit('');
        this.isRequesting = false;
        this.oldPasswordInput = '';
        this.newPasswordInput = '';
        this.confirmPasswordInput = '';
      },
      error => {
        this.isRequesting = false;
      }
    );
  }
  toggleEdit(name: string) {
    const { editMode } = this;
    editMode.profileName = name === 'name' ? !editMode.profileName : false;
    editMode.work = name === 'work' ? !editMode.work : false;
    editMode.password = name === 'password' ? !editMode.password : false;
    if (name === '') {
      editMode.profileName = false;
      editMode.work = false;
      editMode.password = false;
    }
  }
}
