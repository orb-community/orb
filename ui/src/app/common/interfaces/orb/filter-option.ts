import { Observable } from 'rxjs';

export enum FilterTypes {
  Input, // string input
  AutoComplete,
  Select, // allows select one option
  MultiSelect, // allows select multi options
  Checkbox, // on|off option
}

export interface FilterOption {
  name: string;
  prop: string;
  filter: (item: any, prop: any, value: any) => any;
  type: FilterTypes;
  param?: any;
  options?: string[];
  autoSuggestion?: Observable<string[]>;
}

export function filterExact(item: any, prop: any, value: any): boolean {
  return item[prop] === value;
}

export function filterSubstr(item: any, prop: any, value: any) {
  return item[prop].includes(value);
}

export function filterTags(item: any, prop: any, value: any) {
  const values = Object.entries(item[prop]).map(
    (entry) => `${entry[0]}:${entry[1]}`,
  );
  return values.reduce((acc, val) => {
    acc = acc || val.includes(value.replace(' ', ''));
    return acc;
  }, false);
}

export function filterMultiSelect(item: any, prop: any, values: any) {
  return values.reduce((prev, cur) => {
    return item[prop] === cur || prev;
  }, false);
}
