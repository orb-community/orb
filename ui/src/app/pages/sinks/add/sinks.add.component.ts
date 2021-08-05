import { Component, Input } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/sink.interface';
import { STRINGS } from 'assets/text/strings';

@Component({
  selector: 'ngx-sinks-add-component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent {

  strings = STRINGS.sink;

  action = '';

  @Input()
  formData = {
    name: '',
    description: '',
    tags: '',
    backend: '',
    config: {
      remote_host: '',
      username: '',
    },
    metadata: {},
  };

  sink: Sink;
  isEdit: boolean;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.sink = this.router.getCurrentNavigation().extras.state?.sink || null;
    this.isEdit = !!this.sink;
    this.action = this.isEdit ? 'edit' : 'add';
  }

  goBack() {
    this.router.navigate(['../../sinks'], {relativeTo: this.route});
  }

  cancel() {
    this.goBack();
  }

  submit() {
    const action = this.router.routerState.snapshot.url.split('/').pop().toLowerCase();

    if (this.formData.tags !== '') {
      try {
        this.formData.tags = JSON.parse(this.formData.tags);
      } catch (e) {
        this.notificationsService.error('Wrong metadata format', '');
        return;
      }
    }

    this.formData.backend && (this.formData.metadata['backend'] = this.formData.backend);
    if (action === 'add') {
      this.sinksService.addSink(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Sink successfully created', '');
        },
        error => {
          this.notificationsService.error('Sink creation failed', error);
        },
        () => {
          this.goBack();
        },
      );
    }
    if (action === 'edit') {
      this.sinksService.editSink(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Sink successfully updated', '');
        },
        error => {
          this.notificationsService.error('Sink update failed', error);
        },
        () => {
          this.goBack();
        },
      );
    }
  }
}
