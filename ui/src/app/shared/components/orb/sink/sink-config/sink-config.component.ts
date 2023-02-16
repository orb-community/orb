import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-sink-config',
  templateUrl: './sink-config.component.html',
  styleUrls: ['./sink-config.component.scss']
})
export class SinkConfigComponent implements OnInit {

  @Input()
  sink: Sink;

  @Input()
  editMode: boolean;

  @Output()
  editModeChange: EventEmitter<boolean>;

  constructor() { }

  ngOnInit(): void {
  }

}
