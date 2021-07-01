import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { UsersService } from 'app/common/services/users/users.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-users-add-component',
  templateUrl: './users.add.component.html',
  styleUrls: ['./users.add.component.scss'],
})
export class UsersAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
    password: '',
    metadata: {},
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<UsersAddComponent>,
    private usersService: UsersService,
    private notificationsService: NotificationsService,
  ) { }

  cancel() {
    this.dialogRef.close(false);
  }

  submit() {
    if (this.editorMetadata !== '') {
      try {
        this.formData.metadata = JSON.parse(this.editorMetadata) || {};
      } catch (e) {
        this.notificationsService.error('Wrong metadata format', '');
        return;
      }
    }

    if (this.action === 'Create') {
      this.usersService.addUser(this.formData).subscribe(
        resp => {
          this.notificationsService.success('User successfully created', '');
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.usersService.editUser(this.formData).subscribe(
        resp => {
          this.notificationsService.success('User successfully edited', '');
          this.dialogRef.close(true);
        },
      );
    }
  }
}
