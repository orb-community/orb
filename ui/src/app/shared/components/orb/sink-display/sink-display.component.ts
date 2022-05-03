import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-sink-display',
  templateUrl: './sink-display.component.html',
  styleUrls: ['./sink-display.component.scss'],
})
export class SinkDisplayComponent implements OnInit {
  @Input()
  sinks: Sink[];

  @Output()
  sinksChange: EventEmitter<Sink[]>;

  constructor() {
    this.sinks = [];
    this.sinksChange = new EventEmitter<Sink[]>();
  }

  ngOnInit(): void {
  }

  onRemoveSink(sink) {
    this.sinks.splice(this.sinks.indexOf(sink), 1);
    this.sinksChange.emit(this.sinks);
  }

}
