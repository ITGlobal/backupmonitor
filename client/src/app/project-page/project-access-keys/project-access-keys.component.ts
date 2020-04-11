import { Component, Input, Output, EventEmitter } from '@angular/core';
import { IProject, IAccessKey } from 'src/app/api.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { ViewAccessKeyModalComponent } from 'src/app/modals/view-access-key-modal/view-access-key-modal.component';
import { DeleteAccessKeyModalComponent } from 'src/app/modals/delete-access-key-modal/delete-access-key-modal.component';
import { CreateAccessKeyModalComponent } from 'src/app/modals/create-access-key-modal/create-access-key-modal.component';

@Component({
  selector: 'app-project-access-keys',
  templateUrl: './project-access-keys.component.html',
  styleUrls: ['./project-access-keys.component.scss']
})
export class ProjectAccessKeysComponent {
  constructor(private modalService: NgbModal) { }

  @Input() project?: IProject;
  @Input() accessKeys: IAccessKey[];

  @Output() refreshRequested = new EventEmitter<void>();

  refresh() {
    this.refreshRequested.emit();
  }

  showKey(accessKey: IAccessKey) {
    const modalRef = this.modalService.open(ViewAccessKeyModalComponent);
    const instance = modalRef.componentInstance as ViewAccessKeyModalComponent;
    instance.accessKey = accessKey;
  }

  createKey() {
    const modalRef = this.modalService.open(CreateAccessKeyModalComponent);
    const instance = modalRef.componentInstance as CreateAccessKeyModalComponent;
    instance.project = this.project!;

    modalRef.result.then((created: boolean) => {
      if (created) {
        this.refreshRequested.emit();
      }
    })
      .catch(() => { });
  }

  deleteKey(accessKey: IAccessKey) {
    const modalRef = this.modalService.open(DeleteAccessKeyModalComponent);
    const instance = modalRef.componentInstance as DeleteAccessKeyModalComponent;
    instance.project = this.project!;
    instance.accessKey = accessKey;

    modalRef.result.then((deleted: boolean) => {
      if (deleted) {
        this.refresh();
      }
    })
      .catch(() => { });;
  }
}
