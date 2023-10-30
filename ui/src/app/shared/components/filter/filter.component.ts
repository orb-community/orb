import { Component, ElementRef, HostListener, Input, OnInit, ViewChild } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import {
  FilterOption,
  FilterTypes,
  filterString,
} from 'app/common/interfaces/orb/filter-option';
import { FilterService } from 'app/common/services/filter.service';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
  selector: 'ngx-filter',
  templateUrl: './filter.component.html',
  styleUrls: ['./filter.component.scss'],
})
export class FilterComponent implements OnInit {
  @Input()
  availableFilters!: FilterOption[];

  activeFilters$: Observable<FilterOption[]>;

  selectedFilter!: FilterOption | null;

  exact: boolean;

  searchText: string;

  lastSearchText: string;

  loadedSearchText: string;

  filterType = '';

  FiltersTypes = {
    selectAsync: 'selectAsync',
    selectSync: 'selectSync',
    input: 'input',
  };

  showMenu = false;

  currentFilter: FilterOption | null = null;

  filterText: string = '';

  showOptions = true;

  selectedFilterParams = [];

  @ViewChild('filterMenu') filterMenu: ElementRef;

  constructor(
    private filter: FilterService,
    private elRef: ElementRef,
    ) {
    this.exact = false;
    this.availableFilters = [];
    this.activeFilters$ = filter.getFilters().pipe(map((filters) => filters));
  }

  @HostListener('document:click', ['$event'])
  handleOutsideClick(event: Event) {
    const target = event.target as HTMLElement;
    const parentId = (target.parentNode as HTMLElement).id;

    if (!this.elRef.nativeElement.contains(event.target) && parentId !== 'filterMenu' && parentId !== 'remove-button') {
      if (this.currentFilter) {
        if (this.showOptions) {
          this.showOptions = false;
        } else {
          this.currentFilter = null;
          this.selectedFilterParams = [];
          this.filterText = '';
        }
      }
      if (this.showMenu) {
        const icon = document.querySelector('.icon');
        icon.classList.toggle('flipped');
      }
      this.showMenu = false;
    }
  }

  ngOnInit() {
    this.availableFilters = this.availableFilters.filter(filter => filter.name !== 'Name');
    this.searchText = this.filter.searchName || '';
    if (this.filter.searchName) {
      this.searchText = this.filter.searchName;
      this.loadedSearchText = this.searchText;
    } else {
      this.searchText = '';
    }
  }
  onSearchTextChange() {
    if (this.loadedSearchText) {
      this.filter.findAndRemove(this.loadedSearchText, 'Name');
      this.loadedSearchText = undefined;
    }
    if (this.lastSearchText !== '') {
      this.filter.findAndRemove(this.lastSearchText, 'Name');
    }
    if (this.searchText !== '') {
      const filterOptions: FilterOption = {
        name: 'Name',
        prop: 'name',
        param: this.searchText,
        type: FilterTypes.Input,
        filter: filterString,
      };
      this.filter.addFilter(filterOptions);
    }
    this.lastSearchText = this.searchText;
  }

  removeFilter(index: number, filter: any) {
    if (filter.param === this.selectedFilterParams && filter.name === this.currentFilter?.name) {
      this.selectedFilterParams = [];
    }
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

  onFilterClick() {
    this.showMenu = !this.showMenu;
    this.currentFilter = null;
    this.filterText = '';
    this.selectedFilterParams = [];
    const icon = document.querySelector('.icon');
    icon.classList.toggle('flipped');
    this.showOptions = true;
  }
  handleFilterOptionClick(option: any) {
    this.showMenu = false;
    this.currentFilter = option;
    const icon = document.querySelector('.icon');
    icon.classList.toggle('flipped');

    const filterName = this.currentFilter.filter.name;
    if (filterName === 'filterMultiSelect') {
      this.selectedFilterParams = this.filter.getMultiAppliedParams(this.currentFilter.name);

      if (this.currentFilter.autoSuggestion) {
        this.filterType = this.FiltersTypes.selectAsync;
      } else if (this.currentFilter.options) {
        this.filterType = this.FiltersTypes.selectSync;
      }
    } else {
      this.filterType = this.FiltersTypes.input;
    }
  }

  handleClickMultiSelect(event: any, param: any) {
    if (event.target.type !== 'checkbox') {
      if (!this.selectedFilterParams.find(f => f === param)) {
        this.filter.findAndRemove(this.selectedFilterParams, this.currentFilter.name);
        this.selectedFilterParams.push(param);
        this.filter.addFilter({ ...this.currentFilter, param: this.selectedFilterParams });
      }
      this.currentFilter = null;
      this.selectedFilterParams = [];
      this.filterText = '';
    }
  }

  onCheckboxMultiSelect(event: Event, param: any) {
    const checkbox = event.target as HTMLInputElement;
    const isChecked = checkbox.checked;
    const oldParamList = this.selectedFilterParams;

    if (isChecked) {
      this.selectedFilterParams.push(param);
    } else {
      this.selectedFilterParams = this.selectedFilterParams.filter((f) => f !== param);
    }

    this.filter.findAndRemove(oldParamList, this.currentFilter.name);
    if (this.selectedFilterParams.length > 0) {
      this.filter.addFilter({ ...this.currentFilter, param: this.selectedFilterParams });
    }
  }

  handleClickInput(event: any, param: any) {
    if (event.target.type !== 'checkbox') {
      if (this.filter.findFilter(param, this.currentFilter.name) === -1) {
        this.filter.addFilter({ ...this.currentFilter, param: param });
      }
      this.currentFilter = null;
      this.filterText = '';
    }
  }
  onCheckboxInput(event: Event, param: any) {
    const checkbox = event.target as HTMLInputElement;
    const isChecked = checkbox.checked;

    if (isChecked) {
      this.filter.addFilter({ ...this.currentFilter, param: param });
    } else {
      this.filter.findAndRemove(param, this.currentFilter.name);
    }
  }
  onAddFilterButton(param: any) {
    if (this.filterType === this.FiltersTypes.selectAsync) {
      if (this.selectedFilterParams.length > 0) {
        this.filter.findAndRemove(this.selectedFilterParams, this.currentFilter.name);
        this.selectedFilterParams.push(param);
        this.filter.addFilter({ ...this.currentFilter, param: this.selectedFilterParams });
      } else {
        this.selectedFilterParams.push(param);
        this.filter.addFilter({ ...this.currentFilter, param: this.selectedFilterParams });
      }
    } else {
      this.filter.addFilter({ ...this.currentFilter, param: param });
    }
    this.filterText = '';
  }
  isOptionSelected(option: any) {
    return this.activeFilters$.pipe(
      map(filters => {
        return filters.some(filter => {
          const isArray = Array.isArray(filter.param);
          if (isArray) {
            return filter.param.includes(option) && this.currentFilter.name === filter.name;
          }
          return filter.param === option && this.currentFilter.name === filter.name;
        });
      }),
    );
  }

}
