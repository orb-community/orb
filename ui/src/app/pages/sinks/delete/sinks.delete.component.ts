import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-sinks-delete-component',
  templateUrl: './sinks.delete.component.html',
  styleUrls: ['./sinks.delete.component.scss'],
})

export class SinksDeleteComponent {
  @Input() sink = {
    name: '',
    id: '',
  };

  strings = STRINGS.sink;

  sinkName: string = '';

  constructor(
    protected dialogRef: NbDialogRef<SinksDeleteComponent>,
  ) {
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    return this.sinkName.toLowerCase() === this.sink.name.toLowerCase();
  }
}
