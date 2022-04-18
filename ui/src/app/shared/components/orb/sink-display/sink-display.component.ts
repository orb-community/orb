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

  constructor() { }

  ngOnInit(): void {
  }

  onRemoveSink(sinkID) {
    this.sinksChange.emit(this.sinks.filter(sink => sink.id !== sinkID));
  }

}
