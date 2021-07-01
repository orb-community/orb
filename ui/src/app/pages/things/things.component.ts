import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';

import { Thing, PageFilters, TableConfig, TablePage } from 'app/common/interfaces/mainflux.interface';
import { ThingsService } from 'app/common/services/things/things.service';
import { FsService } from 'app/common/services/fs/fs.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';
import { ThingsAddComponent } from './add/things.add.component';

const defSearchBarMs: number = 100;

@Component({
  selector: 'ngx-things-component',
  templateUrl: './things.component.html',
  styleUrls: ['./things.component.scss'],
})
export class ThingsComponent implements OnInit {
  tableConfig: TableConfig = {
    colNames: ['', '', '', 'Name', 'Type', 'ID', 'Key'],
    keys: ['edit', 'delete', 'details', 'name', 'type', 'id', 'key'],
  };
  page: TablePage = {};
  pageFilters: PageFilters = {};

  searchTime = 0;
  columnChar = '|';

  constructor(
    private router: Router,
    private dialogService: NbDialogService,
    private thingsService: ThingsService,
    private fsService: FsService,
    private notificationsService: NotificationsService,
  ) { }

  ngOnInit() {
    this.getThings();
  }

  getThings(name?: string): void {
    this.pageFilters.name = name;
    this.thingsService.getThings(this.pageFilters).subscribe(
      (resp: any) => {
        this.page = {
          offset: resp.offset,
          limit: resp.limit,
          total: resp.total,
          rows: resp.things,
        };

        // Check if there is a type defined in the metadata
        this.page.rows.forEach( (thing: Thing) => {
          thing.type = thing.metadata ? thing.metadata.type : '';
        });
      },
    );
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      this.pageFilters.offset = this.page.offset - this.page.limit;
    }
    if (dir === 'next') {
      this.pageFilters.offset = this.page.offset + this.page.limit;
    }
    this.getThings();
  }

  onChangeLimit(lim: number) {
    this.pageFilters.limit = lim;
    this.getThings();
  }

  openAddModal() {
    this.dialogService.open(ThingsAddComponent, { context: { action: 'Create' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getThings();
        }
      },
    );
  }

  openEditModal(row: any) {
    this.dialogService.open(ThingsAddComponent, { context: { formData: row, action: 'Edit' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.getThings();
        }
      },
    );
  }

  openDeleteModal(row: any) {
    this.dialogService.open(ConfirmationComponent, { context: { type: 'Thing' } }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.thingsService.deleteThing(row.id).subscribe(
            resp => {
              this.page.rows = this.page.rows.filter((t: Thing) => t.id !== row.id);
              this.notificationsService.success('Thing successfully deleted', '');
            },
          );
        }
      },
    );
  }

  onOpenDetails(row: any) {
    if (row.id) {
      this.router.navigate([`${this.router.routerState.snapshot.url}/details/${row.id}`]);
    }
  }

  searchThing(input) {
    const t = new Date().getTime();
    if ((t - this.searchTime) > defSearchBarMs) {
      this.getThings(input);
      this.searchTime = t;
    }
  }

  onClickSave() {
    this.fsService.exportToCsv('mfx_things.csv', this.page.rows);
  }

  onFileSelected(files: FileList) {
    if (files && files.length > 0) {
      const file: File = files.item(0);
      const reader: FileReader = new FileReader();
      reader.readAsText(file);
      reader.onload = () => {
        const csv: string = reader.result as string;
        const lines = csv.split('\n');
        const things: Thing[] = [];

        lines.forEach( line => {
          const col = line.split(this.columnChar);
          const name = col[0];
          if (name !== '' && name !== '<empty string>') {
            let metadata = {};
            if (col[1] !== undefined) {
              try {
                metadata = JSON.parse(col[1]);
              } catch (e) {
                this.notificationsService.warn('Wrong metadata format', '');
              }
            }
            const thing = {
              name: name,
              metadata: metadata,
            };
            things.push(thing);
          }
        });

        this.thingsService.addThings(things).subscribe(
          resp => {
            this.getThings();
          },
        );
      };
    }
  }
}
