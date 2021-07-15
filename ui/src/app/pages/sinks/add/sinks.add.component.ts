import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';

@Component({
  selector: 'ngx-sink-management-add.component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
    description: '',
    metadata: {},
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<SinksAddComponent>,
    private sinkService: SinksService,
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
      this.sinkService.addSink(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Sink Item successfully created', '');
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.sinkService.editSink(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Sink Item successfully edited', '');
          this.dialogRef.close(true);
        },
      );
    }
  }
}
