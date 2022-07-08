import { Component, OnInit } from '@angular/core';
import { OrbService } from 'app/common/services/orb.service';
import { Observable } from 'rxjs';
import { shareReplay } from 'rxjs/operators';

@Component({
  selector: 'ngx-poll-control',
  templateUrl: './poll-control.component.html',
  styleUrls: ['./poll-control.component.scss'],
})
export class PollControlComponent implements OnInit {
  lastUpdate$: Observable<number>;

  constructor(private orb: OrbService) {
    this.lastUpdate$ = this.orb.lastPollUpdate$
      .asObservable()
      .pipe(shareReplay());
  }

  ngOnInit(): void {}

  forceRefresh() {
    this.orb.refreshNow();
  }
}
