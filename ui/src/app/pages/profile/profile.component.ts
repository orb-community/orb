import { Component, OnInit } from '@angular/core';

import { UsersService } from 'app/common/services/users/users.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { User } from 'app/common/interfaces/mainflux.interface';
import { OrbService, pollIntervalKey } from 'app/common/services/orb.service';


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

  showPassword = false;
  showPassword2 = false;
  showPassword3 = false;

  availableTimers = [15, 30, 60]
  selectedTimer: Number;

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
    private orb: OrbService,
  ) { 
    this.oldPasswordInput = '';
    this.newPasswordInput = '';
    this.confirmPasswordInput = '';
    this.selectedTimer = this.getPollInterval();
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
      error => {
        this.isRequesting = false;
      }
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

  switch (name) {
    case 'name':
      editMode.profileName = !editMode.profileName;
      if (!editMode.profileName) {
        this.userFullName = this.user.fullName;
      }
      this.editMode.password = false;
      this.editMode.work = false;
      break;
    case 'work':
      editMode.work = !editMode.work;
      if (!editMode.work) {
        this.userCompany = this.user.company;
      }
      this.editMode.password = false;
      this.editMode.profileName = false;
      break;
    case 'password':
      editMode.password = !editMode.password;
      if (!editMode.password) {
        this.oldPasswordInput = '';
        this.newPasswordInput = '';
        this.confirmPasswordInput = '';
      }
      this.editMode.profileName = false;
      this.editMode.work = false;
      break;
    case '':
      editMode.profileName = false;
      editMode.work = false;
      editMode.password = false;
      break;
  }
}
  setPollInterval(timer) {
    const pollKeyString = (timer * 1000).toString();
    localStorage.setItem(pollIntervalKey, pollKeyString);
    this.orb.pollInterval = timer * 1000;
  }
  getPollInterval() {
    const value = Number(localStorage.getItem(pollIntervalKey));
    return value / 1000;
  }
}
