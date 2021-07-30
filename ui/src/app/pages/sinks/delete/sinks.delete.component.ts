import {Component, Input} from '@angular/core';
import {NbDialogRef} from '@nebular/theme';
import {Sink} from 'app/common/interfaces/sink.interface';

@Component({
  selector: 'ngx-sinks-delete-component',
  templateUrl: './sinks.delete.component.html',
  styleUrls: ['./sinks.delete.component.scss'],
})

export class SinksDeleteComponent {
  sink: Sink;
  @Input() formData = {
    name: '',
  };
  @Input() sinkName: string = '';

  constructor(
      protected dialogRef: NbDialogRef<SinksDeleteComponent>,
  ) {
  }

  onDelete() {
    // TODO check this is the case --lowercase #probablynot
    // if (this.formData.name.toLowerCase() === this.sink.name.toLowerCase()) {
    // this.sinksService.deleteSink(this.sink.id);
    // }
  }

  onClose() {
    this.dialogRef.close(true);
  }
}
