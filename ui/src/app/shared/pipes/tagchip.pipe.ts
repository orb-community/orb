import { Pipe, PipeTransform } from '@angular/core';

/**
 * Tag => Chip Pipe
 * @class TagChipPipe
 */
@Pipe({name: 'tagchip'})
export class TagChipPipe implements PipeTransform {
  /**
   * <{key: string, value: string}>[] => <'{ `${key}`: `${value}` }'>[]
   * @param tagList {<{key: string, value: string}>[]}
   * @return {<string>[]}
   */
  transform(tagList: {key: string, value: string}[]): string[] {
    return tagList.map(tag => tag.key !== '' ? `${tag.key} : ${tag.value || ''}` : '');
  }
}
