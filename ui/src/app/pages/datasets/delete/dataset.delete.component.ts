import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';

@Component({
  selector: 'ngx-dataset-delete-component',
  templateUrl: './dataset.delete.component.html',
  styleUrls: ['./dataset.delete.component.scss'],
})

export class DatasetDeleteComponent {
  @Input() name: string;

  validationInput: string = '';

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
    return this.validationInput.toLowerCase() === this.name.toLowerCase();
  }
}
