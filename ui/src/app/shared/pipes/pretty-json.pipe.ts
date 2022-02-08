import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'prettyJson',
})
export class PrettyJsonPipe implements PipeTransform {

  transform(value: any, args: any[]): any {
    try {
      /**
       * check and try to parse value if it's not an object
       * if it fails to parse which means it is an invalid JSON
       */
      return this.applyColors(
        typeof value === 'object' ? value : JSON.parse(value),
        args[0],
        args[1],
      );
    } catch (e) {
      return this.applyColors({ error: 'Invalid JSON' }, args[0], args[1]);
    }
  }

  applyColors(obj: any, showNumebrLine: boolean = false, padding: number = 4) {
    // line number start from 1
    let line = 1;

    if (typeof obj !== 'string') {
      obj = JSON.stringify(obj, undefined, 3);
    }

    /**
     * Converts special charaters like &, <, > to equivalent HTML code of it
     */
    obj = obj.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    /* taken from https://stackoverflow.com/a/7220510 */

    /**
     * wraps every datatype, key for e.g
     * numbers from json object to something like
     * <span class="number" > 234 </span>
     * this is why needed custom themeClass which we created in _global.css
     * @return final bunch of span tags after all conversion
     */
    obj = obj.replace(
      /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
      (match: any) => {
        // class to be applied inside pre tag
        let themeClass = 'number';
        if (/^"/.test(match)) {
          if (/:$/.test(match)) {
            themeClass = 'key';
          } else {
            themeClass = 'string';
          }
        } else if (/true|false/.test(match)) {
          themeClass = 'boolean';
        } else if (/null/.test(match)) {
          themeClass = 'null';
        }
        return '<span class="' + themeClass + '">' + match + '</span>';
      },
    );

    /**
     * Regex for the start of the line, insert a number-line themeClass tag before each line
     */
    return showNumebrLine
      ? obj.replace(
        /^/gm,
        () =>
          `<span class="number-line pl-3 select-none" >${ String(line++).padEnd(padding) }</span>`,
      )
      : obj;
  }
}
