import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { LoraService } from 'app/common/services/lora/lora.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-lora-add-component',
  templateUrl: './lora.add.component.html',
  styleUrls: ['./lora.add.component.scss'],
})
export class LoraAddComponent {
  @Input() formData = {
    name: '',
    appID: '',
    devEUI: '',
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<LoraAddComponent>,
    private loraService: LoraService,
    private notificationsService: NotificationsService,
  ) { }

  cancel() {
    this.dialogRef.close(false);
  }

  submit() {
    if (this.formData.devEUI === '' || this.formData.appID === '') {
      this.notificationsService.warn('Application ID and Device EUI are required', '');
      return false;
    }

    if (this.action === 'Create') {
      this.loraService.addDevice(this.formData).subscribe(
        resp => {
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.loraService.editDevice(this.formData).subscribe(
        resp => {
          this.dialogRef.close(true);
        },
      );
    }
  }
}
