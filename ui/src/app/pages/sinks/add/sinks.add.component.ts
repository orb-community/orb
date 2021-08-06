import {Component, Input} from '@angular/core';
import {NbDialogRef} from '@nebular/theme';

import {NotificationsService} from 'app/common/services/notifications/notifications.service';
import {SinksService} from 'app/common/services/sinks/sinks.service';

@Component({
  selector: 'ngx-sinks-add-component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
    description: '',
    tags: '',
    backend: '',
    config: {
      remote_host: '',
      username: '',
    },
    metadata: {},
  };
  @Input() action: string = '';

  constructor(
      protected dialogRef: NbDialogRef<SinksAddComponent>,
      private sinksService: SinksService,
      private notificationsService: NotificationsService,
  ) {
  }

  cancel() {
    this.dialogRef.close(false);
  }

  submit() {
    if (this.formData.tags !== '') {
      try {
        this.formData.tags = JSON.parse(this.formData.tags);
      } catch (e) {
        this.notificationsService.error('Wrong metadata format', '');
        return;
      }
    }

    this.formData.backend && (this.formData.metadata['backend'] = this.formData.backend);
    if (this.action === 'Add') {
      this.sinksService.addSink(this.formData).subscribe(
          resp => {
            this.notificationsService.success('Sink successfully created', '');
            this.dialogRef.close(true);
          },
      );
    }
    if (this.action === 'Edit') {
      this.sinksService.editSink(this.formData).subscribe(
          resp => {
            this.notificationsService.success('Sink successfully edited', '');
            this.dialogRef.close(true);
          },
      );
    }
  }
}
