import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';

@Component({
  selector: 'ngx-sink-view',
  templateUrl: './sink.view.component.html',
  styleUrls: ['./sink.view.component.scss']
})
export class SinkViewComponent implements OnInit {
  strings = STRINGS;
  
  isLoading = false;

  sink: Sink;

  sinkSubscription: Subscription;

  editMode = {
    details: false,
    config: false,
  }

  constructor(private cdr: ChangeDetectorRef,
    private notifications: NotificationsService,
    private orb: OrbService,
    private sinks: SinksService,
    ) { }

  ngOnInit(): void {
  }

  fetchData() {

  }

  isEditMode() {
    return Object.values(this.editMode).reduce(
      (prev, cur) => prev || cur,
      false,
    );
  }

  canSave() {

  }

  discard() {

  }

  save() {

  }

  retrieveSink() {
    
  }

}
