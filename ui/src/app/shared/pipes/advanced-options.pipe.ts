import { Pipe, PipeTransform } from '@angular/core';
import { KeyValue } from '@angular/common';

/**
 * Custom Filter
 * @class AdvancedOptionsPipe
 */
@Pipe({name: 'advancedoptions'})
export class AdvancedOptionsPipe implements PipeTransform {
  /**
   * Filter Advanced Options true|false
   * @param items <any[]>
   * @param filter <boolean>
   * @return {<string>}
   */
  transform(items: KeyValue<string, any>[], filter: boolean): KeyValue<string, any>[] {
    if (!items) {
      return items;
    }

    return items.filter(item => {
      return (!!item.value?.props?.advanced && item.value.props.advanced === filter)
      || (filter === false && !item.value?.props?.advanced);
    });

  }
}
