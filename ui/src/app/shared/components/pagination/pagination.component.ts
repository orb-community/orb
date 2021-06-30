import { Component, OnChanges, Input, Output, EventEmitter } from '@angular/core';

import { TablePage } from 'app/common/interfaces/mainflux.interface';

@Component({
  selector: 'ngx-pagination-component',
  templateUrl: './pagination.component.html',
  styleUrls: ['./pagination.component.scss'],
})
export class PaginationComponent implements OnChanges {
  currentPage = 0;
  totalPages = 0;
  limits = [5, 10, 20, 50, 100];

  @Input() page: TablePage = {};
  @Output() changeLimitEvent: EventEmitter<any> = new EventEmitter();
  @Output() changePageEvent: EventEmitter<any> = new EventEmitter();
  constructor(
  ) { }

  ngOnChanges() {
    if (this.page !== undefined) {
      // Ceil offset by limit ratio
      const pageNum = (this.page.offset + 1) / this.page.limit;
      this.currentPage = Math.ceil(pageNum);
      // Calculate the number of pages
      this.totalPages = this.page.total / this.page.limit;
    }
  }

  onChangeLimit(lim: number) {
    this.changeLimitEvent.emit(lim);
  }

  onChangePage(dir: any) {
    if (dir === 'prev' && this.currentPage > 1) {
      this.changePageEvent.emit(dir);
    }
    if (dir === 'next' && this.totalPages > this.currentPage) {
      this.changePageEvent.emit(dir);
    }
  }
}
