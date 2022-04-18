import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { forkJoin, Subscription } from 'rxjs';

@Component({
  selector: 'ngx-sink-display-list',
  templateUrl: './sink-display-list.component.html',
  styleUrls: ['./sink-display-list.component.scss'],
})
export class SinkDisplayListComponent implements OnInit, OnDestroy {
  @Input() sinkIDs: string[];

  sinks: Sink[];

  errors: any[];

  subscription: Subscription;

  constructor(protected sinkService: SinksService) {
  }

  ngOnInit() {
    if (!!this.sinkIDs) {
      this.subscription = forkJoin(this.sinkIDs
        .map(id => this.sinkService.getSinkById(id)))
        .subscribe(sinks => {
          // currently status 404 or any
          this.sinks = sinks.filter(sink => !sink['status']);
          this.errors = sinks.filter(sink => !!sink['status']);
        });
    }
  }

  ngOnDestroy() {
    this?.subscription?.unsubscribe();
  }
}
