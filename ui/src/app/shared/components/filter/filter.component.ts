import { ChangeDetectorRef, Component, ElementRef, HostListener, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
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
export class FilterComponent implements OnInit, OnDestroy {
  @Input()
  availableFilters!: FilterOption[];

  activeFilters$: Observable<FilterOption[]>;

  selectedFilter!: FilterOption | null;

  exact: boolean;

  searchText: string;

  lastSearchText: string;

  loadedSearchText: string;

  optionsMenuType = '';

  showMenu = false;

  currentFilter: FilterOption | null = null;

  filterText: string = '';

  showOptions = true;

  selectedFilterParams = [];

  isMultiTags = false;

  @ViewChild('filtersDisplay') filtersDisplay!: ElementRef;

  @ViewChild('filterList') filterList!: ElementRef;

  @ViewChild('filterMenu') filterMenu!: ElementRef;

  private observer: MutationObserver | null = null;

  constructor(
    private filter: FilterService,
    private elRef: ElementRef,
    private cdr: ChangeDetectorRef,
    ) {
    this.exact = false;
    this.availableFilters = [];
    this.activeFilters$ = filter.getFilters().pipe(map((filters) => filters));
  }
  @HostListener('document:keydown.enter', ['$event'])
  addFilterEnter(event: KeyboardEvent) {
    if (this.filterText && !this.currentFilter?.options) {
      this.onAddFilterButton(this.filterText);
    }
  }


  @HostListener('document:click', ['$event'])
  handleOutsideClick(event: Event) {
    const target = event.target as HTMLElement;
    const parentId = (target?.parentNode as HTMLElement)?.id;

    if (!this.elRef.nativeElement.contains(event.target) && parentId !== 'filterMenu' && !target.closest('#remove-button')) {
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
    this.availableFilters = this.availableFilters.filter(filter => filter.name !== 'Name' && filter.name !== 'MultiTags');
    this.searchText = this.filter.searchName || '';
    if (this.filter.searchName) {
      this.searchText = this.filter.searchName;
      this.loadedSearchText = this.searchText;
    } else {
      this.searchText = '';
    }
  }

  ngOnDestroy() {
    if (this.observer) {
      this.observer.disconnect();
      this.observer = null;
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
    this.selectedFilterParams = [];
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
    this.currentFilter = Object.assign({}, option);
    const icon = document.querySelector('.icon');
    icon.classList.toggle('flipped');

    if (this.currentFilter.type === FilterTypes.MultiSelect) {
      this.selectedFilterParams = this.filter.getMultiAppliedParams(this.currentFilter.name);
      this.optionsMenuType = 'multiselectsync';
    } else if (this.currentFilter.type === FilterTypes.MultiSelectAsync) {
      this.selectedFilterParams = this.filter.getMultiAppliedParams(this.currentFilter.name);
      this.optionsMenuType = 'multiselectasync';
    } else if (this.currentFilter.type === FilterTypes.Tags) {
      this.optionsMenuType = 'tags';
      this.isMultiTags = false;
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
  toggleMultiTags() {
    this.isMultiTags = !this.isMultiTags;
    if (this.isMultiTags) {
      this.currentFilter.name = 'MultiTags';
      this.selectedFilterParams = this.filter.getMultiAppliedParams(this.currentFilter.name);
    } else {
      this.currentFilter.name = 'Tags';
      this.selectedFilterParams = [];
    }
  }
  handleClickTags(event: any, param: any) {
    if (event.target.type !== 'checkbox') {
      if (this.isMultiTags) {
        if (!this.selectedFilterParams.find(f => f === param)) {
          this.filter.findAndRemove(this.selectedFilterParams, this.currentFilter.name);
          this.selectedFilterParams.push(param);
          this.filter.addFilter({ ...this.currentFilter, param: this.selectedFilterParams });
        }
        this.currentFilter = null;
        this.selectedFilterParams = [];
        this.filterText = '';
      } else {
        if (this.filter.findFilter(param, this.currentFilter.name) === -1) {
          this.filter.addFilter({ ...this.currentFilter, param: param });
        }
        this.currentFilter = null;
        this.filterText = '';
      }
    }
  }
  onCheckboxTags(event: Event, param: any) {
    const checkbox = event.target as HTMLInputElement;
    const isChecked = checkbox.checked;
    const oldParamList = this.selectedFilterParams;

    if (this.isMultiTags) {
      if (isChecked) {
        this.selectedFilterParams.push(param);
      } else {
        this.selectedFilterParams = this.selectedFilterParams.filter((f) => f !== param);
      }
      this.filter.findAndRemove(oldParamList, this.currentFilter.name);
      if (this.selectedFilterParams.length > 0) {
        this.filter.addFilter({ ...this.currentFilter, param: this.selectedFilterParams });
      }
    } else {
      if (isChecked) {
        this.filter.addFilter({ ...this.currentFilter, param: param });
      } else {
        this.filter.findAndRemove(param, this.currentFilter.name);
      }
    }
  }

  onAddFilterButton(param: any) {
    const isTags = this.currentFilter.name === 'Tags';
    if (!isTags || (isTags  && this.isMultiTags)) {
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
  hasActiveFilters(): Observable<boolean> {
    return this.activeFilters$.pipe(
      map(filters => filters && (filters.filter((f) => f.name !== 'Name')).length > 0),
    );
  }
}
