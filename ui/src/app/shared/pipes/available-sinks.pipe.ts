import { Pipe, PipeTransform } from '@angular/core';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Pipe({
  name: 'unSelectedSinks',
})
export class UnSelectedSinksPipe implements PipeTransform {

  transform(sinks: Sink[], selectedSinks: Sink[]): Sink[] {
    return sinks.filter(sink => {
      return !selectedSinks.find(sel => sel.id === sink.id) && sink;
    });
  }
}
