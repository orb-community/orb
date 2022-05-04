import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
             selector: 'ngx-sink-display',
             templateUrl: './sink-display.component.html',
             styleUrls: ['./sink-display.component.scss'],
           })
export class SinkDisplayComponent implements OnInit {
  @Input()
  selectedSinks: Sink[];

  @Output()
  selectedSinksChange: EventEmitter<Sink[]>;

  constructor() {
    this.selectedSinks = [];
    this.selectedSinksChange = new EventEmitter<Sink[]>();
  }

  ngOnInit() {
  }

  onRemoveSink(sink) {
    this.selectedSinks.splice(this.selectedSinks.indexOf(sink), 1);
    this.selectedSinksChange.emit(this.selectedSinks);
  }

}
