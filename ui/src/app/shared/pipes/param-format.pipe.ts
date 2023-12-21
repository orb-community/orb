import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'paramformatter',
})
export class ParamFormatterPipe implements PipeTransform {
    transform(value: string | string[]): string {
      if (typeof value === 'string') {
        return value;
      } else if (Array.isArray(value)) {
        return value.join(', ');
      }
      return value;
    }
  }
