import {Component, Input} from '@angular/core';
import {NbDialogRef} from '@nebular/theme';
import {STRINGS} from 'assets/text/strings';

const strings = STRINGS.sink.details;

@Component({
  selector: 'ngx-sinks-details-component',
  templateUrl: './sinks.details.component.html',
  styleUrls: ['./sinks.details.component.scss'],
})
export class SinksDetailsComponent {
  header: string = strings.header;
  close: string = strings.close;
  name: string = strings.name;
  description: string = strings.description;
  backend: string = strings.backend;
  remote_host: string = strings.remote_host;
  ts_created: string = strings.ts_created;

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
    console.error('sink edit unavailable');
  }

  onClose() {
    this.dialogRef.close();
  }
}
