import { Component, Input } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { ApiService, IProject } from 'src/app/api.service';

@Component({
  selector: 'app-delete-project-modal',
  templateUrl: './delete-project-modal.component.html',
  styleUrls: ['./delete-project-modal.component.scss']
})
export class DeleteProjectModalComponent {
  constructor(private modal: NgbActiveModal, private api: ApiService) { }

  @Input() public project: IProject;

  isBusy: boolean;
  confirmed: boolean;
  error?: string;

  ok() {
    if (this.isBusy) {
      return;
    }

    if (!this.confirmed) {
      this.error = 'you should confirm your action';
      return;
    }

    this.isBusy = true;
    this.error = undefined;

    this.api.deleteProject(this.project.id)
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
