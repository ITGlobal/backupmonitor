import { Component, OnInit, forwardRef, Input } from '@angular/core';
import { NG_VALUE_ACCESSOR, ControlValueAccessor } from '@angular/forms';
import { INotificationParams, ApiService } from 'src/app/api.service';
import {
  IAddNotificationTargetModalResult,
  NotificationTargetType,
  AddNotificationTargetModalComponent
} from 'src/app/modals/add-notification-target-modal/add-notification-target-modal.component';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

interface INotificationTarget {
  type: NotificationTargetType;
  value: string;
}

@Component({
  selector: 'app-notification-targets-editor',
  templateUrl: './notification-targets-editor.component.html',
  styleUrls: ['./notification-targets-editor.component.scss'],
  providers: [{
    provide: NG_VALUE_ACCESSOR,
    useExisting: forwardRef(() => NotificationTargetsEditorComponent),
    multi: true
  }]
})
export class NotificationTargetsEditorComponent implements ControlValueAccessor {
  constructor(private modalService: NgbModal, private api: ApiService) { }

  @Input()
  readonly: boolean;

  @Input()
  disabled: boolean;

  private _value: INotificationParams;
  private onChange(_: any) { }

  get value(): INotificationParams {
    return this._value;
  }

  @Input()
  set value(val: INotificationParams) {
    this._value = val;
    this.onChange(this._value);
  }

  writeValue(value: any) {
    this.value = value;
  }

  registerOnChange(fn: any) {
    this.onChange = fn;
  }

  registerOnTouched() { }

  // enabled
  get enabled(): boolean {
    return this._value.enabled;
  }

  set enabled(value: boolean) {
    this._value.enabled = value;
    this.onChange(this._value);
  }

  // has any targets
  hasAnyTargets() {
    return this.slack.length > 0 ||
      this.telegram.length > 0 ||
      this.webhook.length > 0;
  }

  // slack
  get slack(): string[] {
    return this._value.slack;
  }

  set slack(value: string[]) {
    this._value.slack = value;
    this.onChange(this._value);
  }

  // telegram
  get telegram(): string[] {
    return this._value.telegram;
  }

  set telegram(value: string[]) {
    this._value.telegram = value;
    this.onChange(this._value);
  }

  // webhook
  get webhook(): string[] {
    return this._value.webhook;
  }

  set webhook(value: string[]) {
    this._value.webhook = value;
    this.onChange(this._value);
  }

  listTargets(): INotificationTarget[] {
    const array: INotificationTarget[] = [];

    this.slack.forEach((value) => {
      array.push({ type: 'slack', value });
    });
    this.telegram.forEach((value) => {
      array.push({ type: 'telegram', value });
    });
    this.webhook.forEach((value) => {
      array.push({ type: 'webhook', value });
    });

    return array;
  }

  addTarget() {
    const modalRef = this.modalService.open(AddNotificationTargetModalComponent);

    modalRef.result.then(({ type, value }: IAddNotificationTargetModalResult) => {
      switch (type) {
        case 'slack':
          if (this.slack.indexOf(value) < 0) {
            this.slack.push(value);
          }
          break;
        case 'telegram':
          if (this.telegram.indexOf(value) < 0) {
            this.telegram.push(value);
          }
          break;
        case 'webhook':
          if (this.webhook.indexOf(value) < 0) {
            this.webhook.push(value);
          }
          break;
      }
    })
      .catch(() => { });;
  }

  removeTarget(type: NotificationTargetType, value: string) {
    switch (type) {
      case 'slack':
        this.slack.splice(this.slack.indexOf(value));
        break;
      case 'telegram':
        this.telegram.splice(this.telegram.indexOf(value));
        break;
      case 'webhook':
        this.webhook.splice(this.webhook.indexOf(value));
        break;
    }
  }

  testTarget(type: NotificationTargetType, value: string) {
    switch (type) {
      case 'slack':
        this.api.testSlackNotification(value).subscribe();
        break;
      case 'telegram':
        this.api.testTelegramNotification(value).subscribe();
        break;
      case 'webhook':
        this.api.testWebhookNotification(value).subscribe();
        break;
    }
  }


}
