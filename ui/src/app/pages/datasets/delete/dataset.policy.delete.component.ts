import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-dataset-policydelete-component',
  templateUrl: './dataset.policy.delete.component.html',
  styleUrls: ['./dataset.policy.delete.component.scss'],
})

export class DatasetDeleteComponent {
  @Input() sink;

  strings = STRINGS.sink;

  userInput: string = '';

  constructor(
    protected dialogRef: NbDialogRef<DatasetDeleteComponent>,
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
