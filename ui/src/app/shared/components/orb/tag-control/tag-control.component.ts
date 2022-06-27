import { AfterViewInit, Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { STRINGS } from '../../../../../assets/text/strings';
import { Tags } from 'app/common/interfaces/orb/tag';

@Component({
  selector: 'ngx-tag-control',
  templateUrl: './tag-control.component.html',
  styleUrls: ['./tag-control.component.scss'],
})
export class TagControlComponent implements AfterViewInit {

  strings = { ...STRINGS.tags };

  key: string;

  value: string;

  @ViewChild('keyInput')
  keyInput: ElementRef;

  @Input()
  focusAfterViewInit: boolean;

  @Input()
  tags: Tags;

  @Input()
  required: boolean;

  @Output()
  tagsChange: EventEmitter<Tags>;

  constructor() {
    this.required = true;
    this.tags = {};
    this.tagsChange = new EventEmitter<Tags>();
    this.key = '';
    this.value = '';
  }

  ngAfterViewInit() {
    if (this.focusAfterViewInit) this.focus();
  }

  public focus() {
    this.keyInput.nativeElement.focus();
  }

  onAddTag() {
    this.tags[this.key] = this.value;

    this.tagsChange.emit(this.tags);

    this.key = '';
    this.value = '';

    this.keyInput.nativeElement.focus();
  }

  onRemoveTag(tag: any) {
    delete this.tags[tag];
    this.tagsChange.emit(this.tags);
  }
}
