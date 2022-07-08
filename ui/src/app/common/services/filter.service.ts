import { Injectable } from '@angular/core';
<<<<<<< HEAD
import { NavigationEnd, Router } from '@angular/router';
import { FilterOption } from 'app/common/interfaces/orb/filter-option';
import { combineLatest, Observable, ReplaySubject } from 'rxjs';
import { filter, map, shareReplay, tap } from 'rxjs/operators';
=======
import { FilterOption } from 'app/common/interfaces/orb/filter-option';
import { Observable, Subject } from 'rxjs';
>>>>>>> added filters to all list pages

@Injectable({
  providedIn: 'root',
})
export class FilterService {
  private _filters: FilterOption[];

<<<<<<< HEAD
  private filters: ReplaySubject<FilterOption[]>;

  private activeRoute$: Observable<string>;

  constructor(private router: Router) {
    this.activeRoute$ = this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map((event: NavigationEnd) => event.urlAfterRedirects),
      tap((route) => this.loadFromRoute(route)),
      shareReplay(),
    );
    this.filters = new ReplaySubject();
    this._filters = [];

    this.activeRoute$.subscribe();
  }

  private loadFromRoute(route: string) {
    const storedFilters = window.sessionStorage.getItem(route) || '[]';

    this.resetFilters(JSON.parse(storedFilters));
  }

  private saveToRoute(route: string) {
    const filtersToStore = this._filters.map((_filter) => ({
      name: _filter.name,
      prop: _filter.prop,
      param: _filter.param,
    }));
    window.sessionStorage.setItem(route, JSON.stringify(filtersToStore));
  }

  getFilters() {
    return this.filters;
=======
  private filters: Subject<FilterOption[]>;

  private filters$: Observable<FilterOption[]>;

  constructor() {
    this.filters = new Subject();
    this.filters$ = this.filters.asObservable();
    this._filters = [];
    this.cleanFilters();
  }

  getFilters() {
    return this.filters$;
>>>>>>> added filters to all list pages
  }

  resetFilters(filters: FilterOption[]) {
    this._filters = filters;
    this.filters.next(this._filters);
  }

  cleanFilters() {
<<<<<<< HEAD
    this.commitFilter([]);
  }

  private commitFilter(filters: FilterOption[]) {
    this._filters = filters;
    this.filters.next(this._filters);
    this.saveToRoute(this.router.url);
  }

  addFilter(filterToAdd: FilterOption) {
    this.commitFilter([...this._filters, filterToAdd]);
=======
    this._filters = [];
    this.filters.next([]);
  }

  addFilter(filter: FilterOption) {
    this._filters.push(filter);
    this.filters.next(this._filters);
>>>>>>> added filters to all list pages
  }

  removeFilter(index: number) {
    if (index >= 0 && index < this._filters.length) {
<<<<<<< HEAD
      const copy = [...this._filters];
      copy.splice(index, 1);
      this.commitFilter(copy);
    }
  }

  // make a decorator out of this?
  createFilteredList() {
    return (
      itemsList: Observable<any[]>,
      filtersList: Observable<FilterOption[]>,
      filterOptions: FilterOption[],
    ) => {
      return combineLatest([itemsList, filtersList]).pipe(
        map(([agents, _filters]) => {
          let filtered = agents;
          _filters.forEach((_filter) => {
            filtered = filtered.filter((value) => {
              const paramValue = _filter.param;
              const filterDef = filterOptions.find(
                (_item) => _item.name === _filter.name,
              );
              const filterFn = filterDef.filter;
              const propName = filterDef.prop;
              const result =
                !!filterFn && filterFn(value, propName, paramValue);
              return result;
            });
          });

          return filtered;
        }),
      );
    };
  }
=======
      this._filters.splice(index, 1);
      this.filters.next(this._filters);
    }
  }
>>>>>>> added filters to all list pages
}
