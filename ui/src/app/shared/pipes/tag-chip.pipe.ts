import { Pipe, PipeTransform } from '@angular/core';

/**
 * Tag => Chip Pipe
 * @class TagchiplistPipe
 */
@Pipe({name: 'tagchip'})
export class TagChipPipe implements PipeTransform {
  /**
   * <{key: string, value: string}> => <'{ `${key}`: `${value}` }'>
   * @param tagList <{key: string, value: string}>
   * @return {<string>}
   */
  transform({key, value}): string {
    return `${key}: ${value}`;
  }
}
