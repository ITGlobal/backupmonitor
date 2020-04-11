import { Injectable } from '@angular/core';
import Timespan from 'readable-timespan';

@Injectable({
  providedIn: 'root'
})
export class PrettyTimeService {

  constructor() { }

  public formatDuration(ms: number): string {
    const ts = new Timespan();
    const s = ts.parse(ms);

    return s;
  }

  public formatRelative(time: Date, now?: Date): string {
    if (!now) {
      now = new Date(Date.now());
    }

    const age = now.getTime() - time.getTime()

    const s = `${this.formatDuration(age)} ago`;
    return s;
  }
}
