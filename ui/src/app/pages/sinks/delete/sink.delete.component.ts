import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-sink-delete-component',
  templateUrl: './sink.delete.component.html',
  styleUrls: ['./sink.delete.component.scss'],
})

export class SinkDeleteComponent {
  @Input() sink;

  strings = STRINGS.sink;

  userInput: string = '';

  constructor(
    protected dialogRef: NbDialogRef<SinkDeleteComponent>,
  ) {
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    return this.userInput === this.sink.name;
  }
}
