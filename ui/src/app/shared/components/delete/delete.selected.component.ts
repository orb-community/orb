import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-delete-selected-component',
  templateUrl: './delete.selected.component.html',
  styleUrls: ['./delete.selected.component.scss'],
})

export class DeleteSelectedComponent {
  strings = STRINGS.agents;
  @Input() selected: any[] = [];
  @Input() elementName: String;

  validationInput: Number;

  constructor(
    protected dialogRef: NbDialogRef<DeleteSelectedComponent>,
  ) {
  }

  onDelete() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }

  isEnabled(): boolean {
    return this.validationInput === this.selected.length;
  }
}