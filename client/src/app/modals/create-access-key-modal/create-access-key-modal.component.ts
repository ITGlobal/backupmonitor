import { Component, Input } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { ApiService, IProject } from 'src/app/api.service';

@Component({
  selector: 'app-create-access-key-modal',
  templateUrl: './create-access-key-modal.component.html',
  styleUrls: ['./create-access-key-modal.component.scss']
})
export class CreateAccessKeyModalComponent {
  constructor(private modal: NgbActiveModal, private api: ApiService) { }

  @Input() public project: IProject;

  label: string;
  isLabelValid: boolean;

  isBusy: boolean;
  error?: string;

  ok() {
    if (this.isBusy) {
      return;
    }

    this.isLabelValid = true;

    if (!this.label) {
      this.isLabelValid = false;
      return;
    }

    this.isBusy = true;
    this.error = undefined;

    this.api.createProjectAccessKey(this.project.id, this.label)
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
