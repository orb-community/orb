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
  sinksList: Sink[];

  availableSinks: Sink[];

  selectedSink: Sink;

  constructor() {
    this.selectedSinks = [];
    this.selectedSinksChange = new EventEmitter<Sink[]>();
    this.sinksList = [];
    this.availableSinks = [];
  }

  ngOnInit(): void {
    this.updateAvailableSinks();
  }

  onAddSink(sink) {
    this.selectedSinks.push(sink);
    this.selectedSinksChange.emit(this.selectedSinks);
    this.selectedSink = null;
    this.updateAvailableSinks();
  }

  updateAvailableSinks() {
    this.availableSinks = this.sinksList.filter(sink => !this.selectedSinks.includes(sink));
  }

}
