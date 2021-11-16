import { Pipe, PipeTransform } from '@angular/core';
import { KeyValue } from '@angular/common';

/**
 * Tag List tag[] => Chip Pipe
 * @class TaglistChipPipe
 */
@Pipe({name: 'tagchiplist'})
export class TaglistChipPipe implements PipeTransform {
  /**
   * <{key: string, value: string}>[] => <'{ `${key}`: `${value}` }'>[]
   * @param tagList {<{key: string, value: string}>[]}
   * @return {<string>[]}
   */
  transform(tagList: KeyValue<string, string>[]): string[] {
    return tagList?.map(tag => tag.key !== '' ? `${tag.key}: ${tag.value || ''}` : '');
  }
}
