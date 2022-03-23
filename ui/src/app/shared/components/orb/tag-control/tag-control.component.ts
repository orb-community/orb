import { AfterViewInit, Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { STRINGS } from '../../../../../assets/text/strings';

export interface Tag {
  [key: string]: string;
}

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
  selectedTags: Tag;

  @Output()
  selectedTagsChange: EventEmitter<Tag>;

  constructor() {
    this.selectedTags = {};
    this.selectedTagsChange = new EventEmitter<Tag>();
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
    this.selectedTags[this.key] = this.value;

    this.selectedTagsChange.emit(this.selectedTags);

    this.key = '';
    this.value = '';

    this.keyInput.nativeElement.focus();
  }

  onRemoveTag(tag: any) {
    delete this.selectedTags[tag];
    this.selectedTagsChange.emit(this.selectedTags);
  }
}
