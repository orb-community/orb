import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-sink-control',
  templateUrl: './sink-control.component.html',
  styleUrls: ['./sink-control.component.scss'],
})
export class SinkControlComponent implements OnInit {
  @Input()
  selectedSinks: Sink[];

  @Output()
  selectedSinksChange: EventEmitter<Sink[]>;

  @Input()
  availableSinks: Sink[];

  selectedSink: Sink;

  constructor() {
    this.selectedSinks = [];
    this.selectedSinksChange = new EventEmitter<Sink[]>();
    this.availableSinks = [];
  }

  ngOnInit(): void {
  }

  onAddSink(sink) {
    this.selectedSinks.push(sink);
    this.selectedSinksChange.emit(this.selectedSinks);
  }

}
