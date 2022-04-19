import { Component, Input } from '@angular/core';
import { Tags } from 'app/common/interfaces/orb/tag';

@Component({
  selector: 'ngx-tag-display',
  templateUrl: './tag-display.component.html',
  styleUrls: ['./tag-display.component.scss'],
})
export class TagDisplayComponent {

  @Input()
  tags: Tags;

  constructor() {
    this.tags = {};
  }

}
