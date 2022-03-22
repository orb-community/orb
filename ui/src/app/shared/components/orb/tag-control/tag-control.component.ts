import { Component, ElementRef, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core';
import { STRINGS } from '../../../../../assets/text/strings';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'ngx-tag-control',
  templateUrl: './tag-control.component.html',
  styleUrls: ['./tag-control.component.scss'],
})
export class TagControlComponent implements OnInit {

  strings = { ...STRINGS.tags };

  fg: FormGroup;

  @ViewChild('keyInput')
  keyInput: ElementRef;

  @Input()
  selectedTags: { [propName: string]: string };

  @Output()
  selectedTagsChange: EventEmitter<{ [propName: string]: string }>;

  constructor() {
  }

  ngOnInit(): void {
  }

  onAddTag() {
    const { key, value } = this.fg.controls;

    this.selectedTags[key.value] = value.value;

    this.selectedTagsChange.emit(this.selectedTags);

    key.reset('');
    value.reset('');

    this.keyInput.nativeElement.focus();
  }

  onRemoveTag(tag: any) {
    delete this.selectedTags[tag];
    this.selectedTagsChange.emit(this.selectedTags);
  }
}
