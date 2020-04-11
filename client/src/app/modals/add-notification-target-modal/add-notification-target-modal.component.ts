import { Component } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';

export type NotificationTargetType = 'slack' | 'telegram' | 'webhook';

export interface IAddNotificationTargetModalResult {
  type: NotificationTargetType;
  value: string;
}

@Component({
  selector: 'app-add-notification-target-modal',
  templateUrl: './add-notification-target-modal.html',
  styleUrls: ['./add-notification-target-modal.scss']
})
export class AddNotificationTargetModalComponent {
  constructor(private modal: NgbActiveModal) { }

  type: NotificationTargetType = 'slack';
  value: string;
  valueError?: string;

  ok() {
    this.valueError = this.validate(this.type, this.value);

    if (this.valueError) {
      return;
    }

    this.modal.close({
      type: this.type, value: this.value
    });
  }

  private validate(type: NotificationTargetType, value: string): string | undefined {
    if (!value) {
      return 'Value is empty';
    }

    switch (type) {
      case 'slack':
        if (!value.match(/^(@|#)[a-zA-Z0-9-_]+$/)) {
          return 'This is neither a Slack username nor a Slack channel.\n' +
            'Username should start with "@", e.g. "@backup_admin"\n' +
            'Channel name should start with "#", e.g. "#general".';
        }
        break;
      case 'telegram':
        if (!value.match(/^-?[0-9]+$/)) {
          return 'This is not a valid Telegram group ID.\n' +
            'Group ID should be a number.\n';
        }
        break;
      case 'webhook':
        try {
          const url = new URL(value);
          const protocol = (url.protocol || '').toLowerCase();
          if (protocol !== 'http:' && protocol !== 'https:') {
            return 'This is not an HTTP(S) URL';
          }
        } catch {
          return 'This is not a valid URL';
        }
        break;
    }

    return undefined;
  }

  dismiss() {
    this.modal.dismiss();
  }
}
