import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

const strings = STRINGS.sink;

@Component({
  selector: 'ngx-sinks-details-component',
  templateUrl: './sinks.details.component.html',
  styleUrls: ['./sinks.details.component.scss'],
})
export class SinksDetailsComponent {
  header = strings.details.header;
  close = strings.details.close;
  name = strings.propNames.name;
  description = strings.propNames.description;
  backend = strings.propNames.backend;
  remote_host = strings.propNames.config_remote_host;
  ts_created = strings.propNames.ts_created;

  @Input() sink = {
    id: '',
    name: '',
    description: '',
    backend: '',
    config: {
      remote_host: '',
      username: '',
    },
    ts_created: '',
  };

  constructor(
      protected dialogRef: NbDialogRef<SinksDetailsComponent>,
  ) {
  }


  onOpenEdit(row: any) {
    // TODO implement router call
  }

  onClose() {
    this.dialogRef.close();
  }
}
