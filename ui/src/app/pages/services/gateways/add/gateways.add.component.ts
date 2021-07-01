import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { GatewaysService } from 'app/common/services/gateways/gateways.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-gateways-add-component',
  templateUrl: './gateways.add.component.html',
  styleUrls: ['./gateways.add.component.scss'],
})
export class GatewaysAddComponent {
  @Input() formData = {
    name: '',
    externalID: '',
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<GatewaysAddComponent>,
    private gatewaysService: GatewaysService,
    private notificationsService: NotificationsService,
  ) { }

  cancel() {
    this.dialogRef.close(false);
  }

  submit() {
    if (this.formData.name === '' || this.formData.name.length > 32) {
      this.notificationsService.warn(
        'Name is required and must be maximum 32 characters long.', '');
        return false;
    }

    if (this.formData.externalID === '' || this.formData.externalID.length < 8) {
      this.notificationsService.warn(
        'External ID is required and must be at least 8 characters long.', '');
      return false;
    }

    if (this.action === 'Create') {
      this.gatewaysService.addGateway(this.formData).subscribe(
        resp => {
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.gatewaysService.editGateway(this.formData).subscribe(
        resp => {
          this.dialogRef.close(true);
        },
      );
    }
  }
}
