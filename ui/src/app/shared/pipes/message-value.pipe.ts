import { Pipe, PipeTransform } from '@angular/core';

import { MainfluxMsg } from 'app/common/interfaces/mainflux.interface';

@Pipe({name: 'messageValue'})
export class MessageValuePipe implements PipeTransform {
  transform(msg: MainfluxMsg): any {
    if (typeof(msg.value) !== 'undefined') return msg.value;
    if (typeof(msg.bool_value) !== 'undefined') return msg.bool_value;
    if (typeof(msg.string_value) !== 'undefined') return msg.string_value;
    if (typeof(msg.data_value) !== 'undefined') return msg.data_value;
  }
}
