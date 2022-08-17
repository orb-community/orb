import {Component, OnDestroy, OnInit} from '@angular/core';
import {OrbService} from 'app/common/services/orb.service';
import {ActivatedRoute} from '@angular/router';
import {Subscription} from 'rxjs';

@Component({
  selector: 'ngx-metrics',
  templateUrl: './metrics.component.html',
  styleUrls: ['./metrics.component.scss'],
})
export class MetricsComponent implements OnInit, OnDestroy {
  agentID: string;

  metricsSub: Subscription;

  metrics: any;

  constructor(
      protected orb: OrbService,
      protected route: ActivatedRoute,
      ) {
    this.metrics = {};
  }

  ngOnInit(): void {
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.metricsSub = this.orb.getAgentMetricsView(this.agentID).subscribe(value => {
      this.metrics = value;
    }, error => {
      console.error(error);
    });
  }

  ngOnDestroy(): void {
    this.metricsSub?.unsubscribe();
  }

}
