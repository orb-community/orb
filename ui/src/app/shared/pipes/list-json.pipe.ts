import { Pipe, PipeTransform } from '@angular/core';

@Pipe({name: 'jsonlist'})
export class JsonListPipe implements PipeTransform {
    
  transform(object: any): string {
    if (!object || typeof object !== 'object') {
      return '';
    }
    const entries = Object.entries(object);
    const formattedList = entries.map(([key, value]) => `${key}: ${value}`).join(', ');
    return formattedList;
  }
}