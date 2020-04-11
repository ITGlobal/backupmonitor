import { Component, OnInit, Input } from '@angular/core';
import { IBackup, ApiService, IProject } from 'src/app/api.service';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-delete-backup-modal',
  templateUrl: './delete-backup-modal.component.html',
  styleUrls: ['./delete-backup-modal.component.scss']
})
export class DeleteBackupModalComponent {
  constructor(private modal: NgbActiveModal, private api: ApiService) { }

  @Input() public project: IProject;
  @Input() public backup: IBackup;

  isBusy: boolean;
  error?: string;

  ok() {
    if (this.isBusy) {
      return;
    }

    this.isBusy = true;
    this.error = undefined;

    this.api.deleteProjectBackup(this.backup.id)
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
