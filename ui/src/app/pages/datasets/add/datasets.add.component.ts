import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { UserGroupsService } from 'app/common/services/users/groups.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-dataset-explorer-add.component',
  templateUrl: './datasets.add.component.html',
  styleUrls: ['./datasets.add.component.scss'],
})
export class DatasetsAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
    description: '',
    metadata: {},
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<DatasetsAddComponent>,
    private userGroupsService: UserGroupsService,
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
      this.userGroupsService.addGroup(this.formData).subscribe(
        resp => {
          this.notificationsService.success('User Group successfully created', '');
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.userGroupsService.editGroup(this.formData).subscribe(
        resp => {
          this.notificationsService.success('User Group successfully edited', '');
          this.dialogRef.close(true);
        },
      );
    }
  }
}
