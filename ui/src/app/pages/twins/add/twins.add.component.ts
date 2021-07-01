import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { TwinsService } from 'app/common/services/twins/twins.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-twins-add-component',
  templateUrl: './twins.add.component.html',
  styleUrls: ['./twins.add.component.scss'],
})
export class TwinsAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<TwinsAddComponent>,
    private twinsService: TwinsService,
    private notificationsService: NotificationsService,
  ) { }

  cancel() {
    this.dialogRef.close(false);
  }

  submit() {
    if (this.action === 'Create') {
      this.twinsService.addTwin(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Twin successfully created', '');
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.twinsService.editTwin(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Twin successfully edited', '');
          this.dialogRef.close(true);
        },
      );
    }
  }
}
