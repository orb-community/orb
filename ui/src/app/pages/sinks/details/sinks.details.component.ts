import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { Sink } from 'app/common/interfaces/sink.interface';
import { SinksService } from 'app/common/services/sinks/sinks.service';

@Component({
  selector: 'ngx-sinks-details-component',
  templateUrl: './sinks.details.component.html',
  styleUrls: ['./sinks.details.component.scss'],
})
export class SinksDetailsComponent implements OnInit {

  sink: Sink = {};

  constructor(
    private route: ActivatedRoute,
    private sinkService: SinksService,
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');

    this.sinkService.getSinkById(id).subscribe(
      (resp: any) => {
        this.sink = resp;
      },
    );
  }
}
