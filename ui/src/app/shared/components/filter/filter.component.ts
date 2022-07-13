import { Component, Input } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import {
  FilterOption,
  FilterTypes,
} from 'app/common/interfaces/orb/filter-option';
import { FilterService } from 'app/common/services/filter.service';
import { Observable } from 'rxjs';

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

  constructor(private filter: FilterService) {
    this.availableFilters = [];
    this.activeFilters$ = filter.getFilters();
  }

  addFilter(): void {
    if (!this.selectedFilter) return;

    this.filter.addFilter({ ...this.selectedFilter, param: this.filterParam });

    this.selectedFilter = null;
    this.filterParam = null;
  }

  removeFilter(index: number) {
    this.filter.removeFilter(index);
  }

  filterChanged(event: MatSelectChange) {
    this.selectedFilter = { ...event.source.value };
  }
}
