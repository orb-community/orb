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

  availableTypes: any;

  handlersNames: any;

  metrics: any;

  resources: any;

  handlers: any;

  constructor(
      protected orb: OrbService,
      protected route: ActivatedRoute,
      ) {
    this.availableTypes = [];
    this.handlersNames = [];
    this.metrics = {};
    this.resources = {};
    this.handlers = {};
  }

  ngOnInit(): void {
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.metricsSub = this.orb.getAgentMetricsView(this.agentID).subscribe(value => {
      this.metrics = value;
      this.handlers = value.handlers;
      this.resources = value.resources;
      this.availableTypes = value.types;
      this.handlersNames = value.handler_names;
    }, error => {
      console.error(error);
    });
    this.orb.refreshNow();
  }

  isCardinal(value) {
    return typeof value === 'number';
  }

  ngOnDestroy(): void {
    this.metricsSub?.unsubscribe();
  }
}
