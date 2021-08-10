import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { STRINGS } from 'assets/text/strings';
import { sinkTypesList } from 'app/pages/sinks/sinks.component';

@Component({
  selector: 'ngx-sinks-add-component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent {

  strings = STRINGS;

  sinkForm = {
    name: '',
    description: '',
    backend: sinkTypesList.prometheus,
    config: {
      host_name: '',
      username: '',
      password: '',
    },
    tags: {},
  };
  
  tagsForm = { key: '', value: '' };

  sink: Sink;

  sinkTypesList = sinkTypesList;

  isEdit: boolean;

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.sink = this.router.getCurrentNavigation().extras.state?.sink as Sink || null;
    this.isEdit = !!this.sink;
    if (!this.isEdit) {
      this.sink = this.emptySink();
    }
  }

  goBack() {
    this.sink &&
    this.router.navigate(['../../sinks'], { relativeTo: this.route });
  }

  cancel() {
    this.goBack();
  }

  onSubmit() {
    const { name, description, backend, tags, config } = this.sinkForm;

    // Create Sink State
    if (!this.isEdit) {
      this.sinksService.addSink(this.sink).subscribe(
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
    // Edit Sink State
    if (this.isEdit) {
      this.sinksService.editSink(this.sink).subscribe(
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

  addTag() {
    const { key, value } = this.tagsForm;

    // TODO check if all keys need a value or there could
    // be something like 'key': '' <empty string>
    // TODO check if tags are overridable without warning
    if (!this.sink.tags[key]) {
      this.sink.tags = { ...this.sink.tags, ...{ [key]: value } };
      // TODO add create tag confirmation snack/toaster
    } else {
      // TODO add tag already exists snack/toaster
    }
  }

  removeTag(tag) {
    delete this.sink.tags[tag];
    // TODO add delete tag confirmation snack/toaster
  }

  /**
   * Utility: returns empty-initialized sink object
   * @return <Sink> <Sink>{[propName::Sink]: '' | {[propName::SinkConfig|any]: any}};
   */
  emptySink()
    :
    Sink {
    return <Sink>{
      name: '',
      description: '',
      // TODO initialize backend as ''
      backend: sinkTypesList.prometheus,
      config: { remote_host: '', username: '', password: '' },
      tags: {},
    };
  }
}
