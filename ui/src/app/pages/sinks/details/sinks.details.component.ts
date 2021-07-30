import {Component} from '@angular/core';

import {Sink} from 'app/common/interfaces/sink.interface';
import {NbDialogRef} from '@nebular/theme';

@Component({
  selector: 'ngx-sinks-details-component',
  templateUrl: './sinks.details.component.html',
  styleUrls: ['./sinks.details.component.scss'],
})
export class SinksDetailsComponent {

  sink: Sink;

  constructor(
      protected dialogRef: NbDialogRef<SinksDetailsComponent>,
  ) {
  }

  onOpenEdit(row: any) {
    // TODO implement router call
    console.error('sink edit unavailable');
  }
}
