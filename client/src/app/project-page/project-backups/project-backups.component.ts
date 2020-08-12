import { Component, Input, Output, EventEmitter } from '@angular/core';
import { IProject, IBackup } from 'src/app/api.service';
import { IconDefinition } from '@fortawesome/fontawesome-svg-core';
import { faStar as fasStar } from '@fortawesome/free-solid-svg-icons';
import { faStar as farStar } from '@fortawesome/free-regular-svg-icons';
import { PrettyTimeService } from 'src/app/pretty-time.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { DeleteBackupModalComponent } from 'src/app/modals/delete-backup-modal/delete-backup-modal.component';
import * as pretty from 'pretty-bytes';

@Component({
  selector: 'app-project-backups',
  templateUrl: './project-backups.component.html',
  styleUrls: ['./project-backups.component.scss']
})
export class ProjectBackupsComponent {
  constructor(private modalService: NgbModal, private time: PrettyTimeService) {
  }

  @Input() project?: IProject;
  @Input() backups: IBackup[];

  @Output() refreshRequested = new EventEmitter<void>();

  refresh() {
    this.refreshRequested.emit();
  }

  getBackupIcon(b: IBackup): IconDefinition {
    if (b.type === 'last') {
      return fasStar;
    }

    return farStar;
  }

  getBackupAge(backup: IBackup): string {
    return this.time.formatRelative(backup.time);
  }

  getBackupDownloadUrl(backup: IBackup): string {
    return `${window.location.protocol}//${window.location.host}/api/backup/${backup.id}`;
  }

  getBackupSize(backup: IBackup): string {
    if (!backup.length || backup.length < 0) {
      return '';
    }

    const size = pretty(backup.length);
    return size;
  }

  deleteBackup(backup: IBackup) {
    const modalRef = this.modalService.open(DeleteBackupModalComponent);
    const instance = modalRef.componentInstance as DeleteBackupModalComponent;
    instance.project = this.project!;
    instance.backup = backup;

    modalRef.result.then((deleted: boolean) => {
      if (deleted) {
        this.refresh();
      }
    })
      .catch(() => { });;
  }
}
