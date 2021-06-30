import { Pipe, PipeTransform } from '@angular/core';

// Convert from seconds to milliseconds
@Pipe({name: 'toMillisecs'})
export class ToMillisecsPipe implements PipeTransform {
  transform(time: number): number {
    return time * 1000;
  }
}
