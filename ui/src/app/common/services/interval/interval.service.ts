import { Injectable } from '@angular/core';

type Callback = () => void;
const defaultInterval: number = 1000 * 20; // ms * sec

@Injectable()
export class IntervalService {
  ids: number[] = [];
  callbacks: Callback[] = [];

  constructor() {}

  set(context: any, callback: Callback, interval?: number) {
    if (this.callbacks.indexOf(callback) > -1) {
      return;
    }
    this.callbacks.push(callback);

    interval = interval || defaultInterval;
    const id = window.setInterval(callback.bind(context), interval);
    this.ids.push(id);
  }

  clear() {
    this.ids.forEach(id => window.clearInterval(id));
  }
}
