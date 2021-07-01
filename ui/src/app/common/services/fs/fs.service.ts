import { Injectable } from '@angular/core';

@Injectable()
export class FsService {

  constructor(
  ) { }

  exportToCsv(filename: string, rows: any) {
    if (!rows || !rows.length) {
      return;
    }
    const separator = ',';
    const keys = Object.keys(rows[0]);
    const bom = new Uint8Array([0xEF, 0xBB, 0xBF]); // UTF-8 BOM
    const csvContent =
      keys.join(separator) +
      '\n' +
      rows.map(row => {
        return keys.map(k => {
          let cell = row[k] === undefined ? '' : row[k];
          if (k === 'metadata') {
            try {
              cell = JSON.stringify(row[k]);
              cell =  cell.replace(/[\,]/g, '|');
            } catch (e) {
            }
          }
          return cell;
        }).join(separator);
      }).join('\n');

    const blob = new Blob([bom, csvContent], { type: 'text/csv;charset=utf-8;' });
    if (navigator.msSaveBlob) { // IE 10+
      navigator.msSaveBlob(blob, filename);
    } else {
      const link = document.createElement('a');
      if (link.download !== undefined) {
        // Browsers that support HTML5 download attribute
        const url = URL.createObjectURL(blob);
        link.setAttribute('href', url);
        link.setAttribute('download', filename);
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      }
    }
  }
}
