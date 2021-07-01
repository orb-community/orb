import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { ThingsService } from 'app/common/services/things/things.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-things-add-component',
  templateUrl: './things.add.component.html',
  styleUrls: ['./things.add.component.scss'],
})
export class ThingsAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
    type: '',
    metadata: {
      type: '',
    },
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<ThingsAddComponent>,
    private thingsService: ThingsService,
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

    this.formData.type && (this.formData.metadata.type = this.formData.type);

    if (this.action === 'Create') {
      this.thingsService.addThing(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Thing successfully created', '');
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.thingsService.editThing(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Thing successfully edited', '');
          this.dialogRef.close(true);
        },
      );
    }
  }
}
