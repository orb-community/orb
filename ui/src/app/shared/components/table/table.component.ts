import { Component, Input, Output, EventEmitter } from '@angular/core';

import { TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';

@Component({
  selector: 'ngx-table-component',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
})
export class TableComponent {
  isHover: boolean;

  isObject(val: any): boolean {
    return typeof val === 'object';
  }

  @Input() config: TableConfig = {};
  @Input() page: TablePage = {};
  @Output() addEvent: EventEmitter<any> = new EventEmitter();
  @Output() checkEvent: EventEmitter<any> = new EventEmitter();
  @Output() delEvent: EventEmitter<any> = new EventEmitter();
  @Output() detailsEvent: EventEmitter<any> = new EventEmitter();
  @Output() editEvent: EventEmitter<any> = new EventEmitter();
  @Output() hoverEvent: EventEmitter<any> = new EventEmitter();

  constructor() {
  }

  onAdd() {
    this.addEvent.emit();
  }

  onDetails(row: any) {
    this.detailsEvent.emit(row);
  }

  onEdit(row: any) {
    this.editEvent.emit(row);
  }

  onDelete(row: any) {
    this.delEvent.emit(row);
  }

  onToggleCheckbox(row: any) {
    this.checkEvent.emit(row);
  }

  onMouseEnter(evt: any, row: any) {
    row.isHover = true;
  }

  onMouseLeave(evt: any, row: any) {
    row.isHover = false;
  }
}
