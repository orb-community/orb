import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { IntervalService } from 'app/common/services/interval/interval.service';
import { TwinsService } from 'app/common/services/twins/twins.service';
import { Twin, TableConfig, TablePage, PageFilters } from 'app/common/interfaces/mainflux.interface';

@Component({
  selector: 'ngx-twins-states-component',
  templateUrl: './twins.states.component.html',
  styleUrls: ['./twins.states.component.scss'],
})
export class TwinsStatesComponent implements OnInit, OnDestroy {
  lowerLimit = 0;
  upperLimit = 10;
  interval = 5 * 1000;
  intervalID: number;

  twin: Twin = {};

  twinID: string = '';
  tableConfig: TableConfig = {
    colNames: ['ID', 'Definition', 'Created', 'Payload'],
    keys: ['id', 'definition', 'created', 'payload'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  constructor(
    private route: ActivatedRoute,
    private twinsService: TwinsService,
    private intervalService: IntervalService,
  ) { }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    this.getTwin(id);
  }

  getTwin(id: string) {
    this.twinsService.getTwin(id).subscribe(
      (resp: Twin) => {
        this.twin = resp;
        this.getStates();
        this.intervalService.set(this, this.getStates, this.interval);
      },
    );
  }

  getStates() {
    this.twinsService.listStates(this.twin.id, this.pageFilters.offset, this.pageFilters.limit).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.states,
        };
      },
    );
  }

  lower($event) {
    const val = +$event.srcElement.value;
    if (Number.isInteger(val)) {
      this.lowerLimit = val;

      this.pageFilters.offset = val - 1;
      this.pageFilters.offset = Math.max(0, this.pageFilters.offset);
      this.pageFilters.limit = this.upperLimit - this.pageFilters.offset;
      this.pageFilters.limit = Math.max(0, this.pageFilters.limit);

      this.getStates();
    }
  }

  upper($event) {
    const val = +$event.srcElement.value;

    if (Number.isInteger(val)) {
      this.upperLimit = val;

      this.pageFilters.limit = val - this.pageFilters.offset;
      this.pageFilters.limit = Math.max(0, this.pageFilters.limit);
      this.pageFilters.offset = this.lowerLimit - 1;
      this.pageFilters.offset = Math.max(0, this.pageFilters.offset);

      this.getStates();
    }
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      this.pageFilters.offset = this.page.offset - this.page.limit;
    }
    if (dir === 'next') {
      this.pageFilters.offset = this.page.offset + this.page.limit;
    }
    this.getStates();
  }

  onChangeLimit(lim: number) {
    this.pageFilters.limit = lim;
    this.getStates();
  }

  ngOnDestroy() {
    this.intervalService.clear();
  }
}
