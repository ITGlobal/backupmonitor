import { Component, OnInit } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { ApiService } from 'src/app/api.service';

@Component({
  selector: 'app-change-password-modal',
  templateUrl: './change-password-modal.component.html',
  styleUrls: ['./change-password-modal.component.scss']
})
export class ChangePasswordModalComponent {
  constructor(private modal: NgbActiveModal, private api: ApiService) {
  }

  oldPassword: string;
  oldPasswordError?: string;

  newPassword: string;
  newPasswordError?: string;

  newPasswordAgain: string;
  newPasswordAgainError?: string;

  isBusy: boolean;
  error?: string;

  ok() {
    if (this.isBusy) {
      return;
    }

    this.oldPasswordError = undefined;
    this.newPasswordError = undefined;
    this.newPasswordAgainError = undefined;

    this.error = undefined;

    if (!this.oldPassword) {
      this.oldPasswordError = 'This field is required';
      return;
    }

    if (!this.newPassword) {
      this.newPasswordError = 'This field is required';
      return;
    }

    if (!this.newPasswordAgain) {
      this.newPasswordAgainError = 'This field is required';
      return;
    }

    if (this.newPassword !== this.newPasswordAgain) {
      this.newPasswordAgainError = 'Passwords do not match';
      return;
    }

    this.isBusy = true;

    this.api.changePassword(this.oldPassword, this.newPassword)
      .subscribe(
        () => {
          this.modal.close(true);
          this.isBusy = false;
        },
        (e) => {
          this.isBusy = false;
          this.error = e;
        });
  }

  dismiss() {
    this.modal.dismiss();
  }
}
