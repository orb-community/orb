import { Component, HostListener, Input } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import {
  FilterOption,
  FilterTypes,
  filterString,
} from 'app/common/interfaces/orb/filter-option';
import { FilterService } from 'app/common/services/filter.service';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';

@Component({
  selector: 'ngx-filter',
  templateUrl: './filter.component.html',
  styleUrls: ['./filter.component.scss'],
})
export class FilterComponent {
  @Input()
  availableFilters!: FilterOption[];

  activeFilters$: Observable<FilterOption[]>;

  type = FilterTypes;

  selectedFilter!: FilterOption | null;

  filterParam: any;

  exact: boolean;

  searchText: string;

  lastSearchText: string;

  constructor(private filter: FilterService) {
    this.exact = false;
    this.availableFilters = [];
    this.activeFilters$ = filter.getFilters().pipe(map((filters) => filters));
    this.searchText = '';
  }
  
  ngOnInit() {
    this.availableFilters = this.availableFilters.filter(filter => filter.name !== 'Name');
  }
  onSearchTextChange() {
    if (this.lastSearchText !== '') {
      this.filter.removeFilterByParam(this.lastSearchText);
    }
    if (this.searchText !== '') {
      const filterOptions: FilterOption = {
        name: 'Name',
        prop: 'name',
        param: this.searchText,
        type: FilterTypes.Input,
        filter: filterString,
      }
      this.filter.addFilter(filterOptions);
    }
    this.lastSearchText = this.searchText;
  }
  addFilter(): void {
    if (!this.selectedFilter) return;

    this.filter.addFilter({ ...this.selectedFilter, param: this.filterParam });

    this.selectedFilter = null;
    this.filterParam = null;
  }

  @HostListener('window:keydown.enter', ['$event'])
  handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter' && this.filterParam) {
      this.addFilter();
    }
  }

  removeFilter(index: number) {
    this.filter.removeFilter(index);
  }

  filterChanged(event: MatSelectChange) {
    this.selectedFilter = { ...event.source.value };
  }

  clearAllFilters() {
    this.filter.cleanFilters();
  }

  toggleExactMatch() {
    this.selectedFilter.exact = !this.selectedFilter.exact;
  }
}
