import {Component, Input} from '@angular/core';
import {NbDialogRef} from '@nebular/theme';
import {SinksService} from 'app/common/services/sinks/sinks.service';

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

  sinkName: string = '';

  constructor(
      protected dialogRef: NbDialogRef<SinksDeleteComponent>,
      protected sinkService: SinksService,
  ) {
  }

  onDelete() {
    this.sinkService.deleteSink(this.sink.id);
  }

  onClose() {
    this.dialogRef.close(true);
  }

  isEnabled(): boolean {
    return this.sinkName.toLowerCase() === this.sink.name.toLowerCase();
  }
}
