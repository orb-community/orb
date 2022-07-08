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
  filter: (item: any, value: any) => any;
  type: FilterTypes;
  param?: any;
  options?: string[];
  autoSuggestion?: Observable<string[]>;
}
