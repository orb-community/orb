import { Component, Input, OnInit } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';

@Component({
  selector: 'ngx-combined-tag',
  templateUrl: './combined-tag.component.html',
  styleUrls: ['./combined-tag.component.scss'],
})
export class CombinedTagComponent implements OnInit {
  @Input()
  agent: Agent;

  constructor() { }

  ngOnInit(): void {
  }

}
