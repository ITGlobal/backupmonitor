import { Component, Input } from '@angular/core';
import { IProject } from 'src/app/api.service';
import { IconDefinition } from '@fortawesome/fontawesome-svg-core';
import {
  faExclamationTriangle,
  faExclamationCircle,
  faCheckCircle,
  faQuestion,
  faGlobe
} from '@fortawesome/free-solid-svg-icons';
import { faSlack, faTelegram } from '@fortawesome/free-brands-svg-icons';
import { PrettyTimeService } from 'src/app/pretty-time.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { DeleteProjectModalComponent } from 'src/app/modals/delete-project-modal/delete-project-modal.component';
import { Router } from '@angular/router';


interface INotificationTarget {
  icon: IconDefinition;
  type: string;
  value: string;
}

@Component({
  selector: 'app-project-summary',
  templateUrl: './project-summary.component.html',
  styleUrls: ['./project-summary.component.scss']
})
export class ProjectSummaryComponent {
  constructor(private time: PrettyTimeService, private router: Router, private modalService: NgbModal) {
  }

  @Input() project?: IProject;

  getBackupStatusClass(): string {
    switch (this.project?.backupStatus) {
      case 'ok':
        return 'text-success';

      case 'none':
        return 'text-warning';

      case 'outdated':
        return 'text-danger';

      default:
        return '';
    }
  }

  getBackupStatusText(): string {
    switch (this.project?.backupStatus) {
      case 'ok':
        return 'Backup is up to date';

      case 'none':
        return 'No backup has been taken';

      case 'outdated':
        return 'Backup is out of date';

      default:
        return this.project?.backupStatus || '';
    }
  }

  getLastBackupStatusText(): string {
    if (!this.project?.lastBackup) {
      return 'No backups are available';
    }

    const str = `Last backup has been taken ${this.time.formatRelative(this.project.lastBackup.time)}`;
    return str;
  }

  getBackupPeriodText(): string {
    if (!this.project) {
      return '';
    }

    const str = `Backups are expected to be taken every ${this.time.formatDuration(this.project.backupFrequency * 1000)}`
    return str;
  }

  getBackupRetentionText(): string {
    if (!this.project) {
      return '';
    }

    const str = `Will keep at least ${this.project.backupRetention} last backups`;
    return str;
  }

  getBackupStatusIcon(): IconDefinition | null {
    switch (this.project?.backupStatus) {
      case 'ok':
        return faCheckCircle;

      case 'none':
        return faExclamationTriangle;

      case 'outdated':
        return faExclamationCircle;

      default:
        return faQuestion;
    }
  }

  getNotificationTargets(): INotificationTarget[] {
    const result: INotificationTarget[] = [];

    const project = this.project;
    if (project) {
      if (project.notifications.slack) {
        project.notifications.slack.forEach((x) => {
          result.push({
            icon: faSlack,
            type: 'Slack user/channel',
            value: x
          })
        })
      }

      if (project.notifications.telegram) {
        project.notifications.telegram.forEach((x) => {
          result.push({
            icon: faTelegram,
            type: 'Telegram group',
            value: x
          })
        })
      }

      if (project.notifications.webhook) {
        project.notifications.webhook.forEach((x) => {
          result.push({
            icon: faGlobe,
            type: 'Webhook',
            value: x
          })
        })
      }
    }

    return result;
  }

  hasNotificationTargets(): boolean {
    const project = this.project;
    if (!project) {
      return false;
    }

    return (project.notifications.slack?.length > 0) ||
      (project.notifications.telegram?.length > 0) ||
      (project.notifications.webhook?.length > 0);
  }

  deleteProject() {
    const modalRef = this.modalService.open(DeleteProjectModalComponent);
    const instance = modalRef.componentInstance as DeleteProjectModalComponent;
    instance.project = this.project!;

    modalRef.result.then((deleted: boolean) => {
      if (deleted) {
        this.router.navigate(['/projects']);
      }
    })
      .catch(() => { });;
  }
}
