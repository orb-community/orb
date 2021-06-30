import { Component, Input, OnInit, OnChanges, OnDestroy } from '@angular/core';
import { Channel, Thing, MainfluxMsg, Message, MsgFilters, Dataset,
  TableConfig, TablePage, ReaderUrl } from 'app/common/interfaces/mainflux.interface';
import { IntervalService } from 'app/common/services/interval/interval.service';
import { MessagesService } from 'app/common/services/messages/messages.service';
import { ChannelsService } from 'app/common/services/channels/channels.service';
import { environment } from 'environments/environment';
import { MessageValuePipe } from 'app/shared/pipes/message-value.pipe';

@Component({
  selector: 'ngx-message-monitor',
  templateUrl: './message-monitor.component.html',
  styleUrls: ['./message-monitor.component.scss'],
})
export class MessageMonitorComponent implements OnInit, OnChanges, OnDestroy {
  messages: MainfluxMsg[] = [];
  chanID = '';

  mode: string = 'json';
  modes: string[] = ['json', 'table', 'chart'];
  httpAdaptType: string = 'float';
  httpAdaptVal: any;
  httpAdaptTypes: string[] = ['float', 'bool', 'string', 'data'];

  msgDatasets: Dataset[] = [];

  filters: MsgFilters = {
    offset: 0,
    limit: 20,
    publisher: '',
    subtopic: '',
    name: '',
    from: 0,
    to: 0,
  };

  readerUrl: ReaderUrl = {
    prefix: environment.readerPrefix,
    suffix: environment.readerSuffix,
  };

  publishers: Thing[] = [];

  tableConfig: TableConfig = {
    colNames: ['Name', 'Value', 'Time', 'Subtopic', 'Channel', 'Publisher', 'Protocol'],
    keys: ['name', 'value', 'time', 'subtopic', 'channel', 'publisher', 'protocol'],
  };
  messagesPage: TablePage = {};

  @Input() channels: Channel[] = [];
  @Input() thingKey: string;
  constructor(
    private intervalService: IntervalService,
    private messagesService: MessagesService,
    private channelsService: ChannelsService,
    private messageValuePipe: MessageValuePipe,
  ) {}

  ngOnInit() {
    this.intervalService.set(this, this.getChannelMessages);
  }

  getRangeDate(event) {
    if (event.start && event.end) {
      this.filters = {
        from: new Date(event.start).getTime() / 1000,
        to: new Date(event.end).getTime() / 1000,
      };
    }

    this.getChannelMessages();
  }

  ngOnChanges() {
    if (this.channels === undefined) {
        return;
    }

    if (this.channels.length > 0 && this.channels[0].id && this.thingKey !== '') {
      this.chanID = this.channels[0].id;
      this.getChannelMessages();
    }
  }

  getChannelMessages() {
    if (this.chanID === '' || this.thingKey === '') {
      return;
    }

    switch (this.httpAdaptType) {
      case 'string':
        this.filters.vs = this.httpAdaptVal;
        break;
      case 'data':
        this.filters.vd = this.httpAdaptVal;
        break;
      case 'bool':
        this.filters.vb = this.httpAdaptVal;
        break;
      case 'float':
        this.filters.v = this.httpAdaptVal;
        break;
    }

    this.messagesPage.rows = [];
    this.messagesService.getMessages(this.chanID, this.thingKey, this.filters, this.readerUrl).subscribe(
      (resp: any) => {
        if (resp.messages) {
          this.messagesPage = {
            offset: resp.offset,
            limit: resp.limit,
            total: resp.total,
            rows: resp.messages.map((msg: MainfluxMsg) => {
              msg.value = this.messageValuePipe.transform(msg);
              return msg;
            }),
          };
          this.msgDatasets = [{
            label: `Channel: ${this.chanID}`,
            messages: <Message[]>this.messagesPage.rows,
          }];
        }
      },
    );

    this.channelsService.connectedThings(this.chanID).subscribe(
      (resp: any) => {
        if (resp.things) {
          this.publishers = resp.things;
        }
      },
    );
  }

  onChangeLimit(lim: number) {
    this.filters.limit = lim;
    this.getChannelMessages();
  }

  onChangePage(dir: any) {
    if (dir === 'prev') {
      this.filters.offset = this.messagesPage.offset - this.messagesPage.limit;
    }
    if (dir === 'next') {
      this.filters.offset = this.messagesPage.offset + this.messagesPage.limit;
    }
    this.getChannelMessages();
  }

  ngOnDestroy(): void {
    this.intervalService.clear();
  }
}
