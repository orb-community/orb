import { Pipe, PipeTransform } from '@angular/core';

// Convert from seconds to milliseconds
@Pipe({ name: 'tagsArrayPipe' })
export class TagsArrayPipe implements PipeTransform {
  transform(tags: { [propName: string]: string }): Array<{ key: string, value: string }> {
    return Object.keys(tags).map(tagKey => ({ key: `${tagKey}`, value: `${tags[tagKey]}` }));
  }
}
