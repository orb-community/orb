import { Pipe, PipeTransform } from '@angular/core';
import { KeyValue } from '@angular/common';

/**
 * Tag => Chip Pipe
 * @class TagListChipPipe
 */
@Pipe({name: 'taglistchip'})
export class TagListChipPipe implements PipeTransform {
  /**
   * <{key: string, value: string}>[] => <'{ `${key}`: `${value}` }'>[]
   * @param tagList {<{key: string, value: string}>[]}
   * @return {<string>[]}
   */
  transform(tagList: {key: string, value: string}[]): string[] {
    return tagList?.map(tag => tag.key !== '' ? `${tag.key}: ${tag.value || ''}` : '');
  }
}

/**
 * Tag => Chip Pipe
 * @class TagChipPipe
 */
@Pipe({name: 'tagchip'})
export class TagChipPipe implements PipeTransform {
  /**
   * <Map<string,string>>[] => <'{ `${key}`: `${value}` }'>
   * @param tagList {<{key: string, value: string}>[]}
   * @return {<string>[]}
   */
  transform(tag: KeyValue<string, string>): string {
    return `${tag.key}: ${tag.value}`;
  }
}

/**
 * Tag => Chip Pipe
 * @class TagKeyPipe
 */
@Pipe({name: 'tagkey'})
export class TagKeyPipe implements PipeTransform {
  /**
   * <'Map<string, string>'> -> `${[0]}`
   * @param tag {<{key: string, value: string}>[]}
   * @return {<string>[]}
   */
  transform(tag: KeyValue<string, string>): string {
    return `${tag.key}: ${tag.value}`;
  }
}

/**
 * Tag => Chip Color Pipe
 * @class TagChipColorPipe
 */
 @Pipe({name: 'tagchipcolor'})
 export class TagChipColorPipe implements PipeTransform {
   /**
    * <`hsl(${h}deg, 90%, 65%)`>
    * @param tagList {<{key: string, value: string}>[]}
    * @return {<string>}
    */

   transform(tag: string): string {
    if (tag !== '') {
      const h = Math.abs(
                  `${tag}}`
                  .split('')
                  .map(v => v.charCodeAt(0))
                  .reduce((a, v) => a + ((a << 7) + (a << 3)) ^ v) % 360);
      return `hsl(${h}, 90%, 65%)`;
    } else return 'transparent';
  }
 }
