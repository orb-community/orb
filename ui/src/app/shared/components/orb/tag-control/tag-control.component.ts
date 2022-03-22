import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { STRINGS } from '../../../../../assets/text/strings';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

export interface Tag {
  [key: string]: string;
}

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
  selectedTags: Tag;

  @Output()
  selectedTagsChange: EventEmitter<Tag>;

  constructor(
    formBuilder: FormBuilder,
  ) {
    this.selectedTags = {};
    this.selectedTagsChange = new EventEmitter<Tag>();
    this.fg = formBuilder.group({
      key: ['', [Validators.required]],
      value: [''],
    });
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
