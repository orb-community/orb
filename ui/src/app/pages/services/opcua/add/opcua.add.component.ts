import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

import { OpcuaService } from 'app/common/services/opcua/opcua.service';
import { OpcuaTableRow } from 'app/common/interfaces/opcua.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-opcua-add-component',
  templateUrl: './opcua.add.component.html',
  styleUrls: ['./opcua.add.component.scss'],
})
export class OpcuaAddComponent {
  @Input() formData: OpcuaTableRow = {
    name: '',
    serverURI: '',
    nodeID: '',
  };
  @Input() action: string = '';
  constructor(
    protected dialogRef: NbDialogRef<OpcuaAddComponent>,
    private opcuaService: OpcuaService,
    private notificationsService: NotificationsService,
  ) { }

  cancel() {
    this.dialogRef.close(false);
  }

  submit() {
    if (this.formData.serverURI === '' || this.formData.nodeID === '') {
      this.notificationsService.warn('Server URI and Node ID are required', '');
      return false;
    }

    const nodes = [{
      name: this.formData.name,
      serverURI: this.formData.serverURI,
      nodeID: this.formData.nodeID,
    }];

    if (this.action === 'Create') {
      this.opcuaService.addNodes(this.formData.serverURI, nodes).subscribe(
        resp => {
          this.dialogRef.close(true);
        },
      );
    }
    if (this.action === 'Edit') {
      this.opcuaService.editNode(this.formData).subscribe(
        resp => {
          this.dialogRef.close(true);
        },
      );
    }
  }
}
