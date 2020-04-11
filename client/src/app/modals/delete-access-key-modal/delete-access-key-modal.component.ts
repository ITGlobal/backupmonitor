import { Component, OnInit, Input } from '@angular/core';
import { IAccessKey, ApiService, IProject } from 'src/app/api.service';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-delete-access-key-modal',
  templateUrl: './delete-access-key-modal.component.html',
  styleUrls: ['./delete-access-key-modal.component.scss']
})
export class DeleteAccessKeyModalComponent {
  constructor(private modal: NgbActiveModal, private api: ApiService) { }

  @Input() public project: IProject;
  @Input() public accessKey: IAccessKey;

  isBusy: boolean;
  error?: string;

  ok() {
    if (this.isBusy) {
      return;
    }

    this.isBusy = true;
    this.error = undefined;

    this.api.deleteProjectAccessKey(this.project.id, this.accessKey.id)
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
