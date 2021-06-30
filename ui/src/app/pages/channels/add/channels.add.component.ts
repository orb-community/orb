import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { ChannelsService } from 'app/common/services/channels/channels.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-channels-add-component',
  templateUrl: './channels.add.component.html',
  styleUrls: ['./channels.add.component.scss'],
})
export class ChannelsAddComponent {
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
    protected dialogRef: NbDialogRef<ChannelsAddComponent>,
    private channelsService: ChannelsService,
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
      this.channelsService.addChannel(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Channel successfully created', '');
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.channelsService.editChannel(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Channel successfully edited', '');
          this.dialogRef.close(true);
        },
      );
    }
  }
}
