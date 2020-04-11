import { Component, Input } from '@angular/core';
import { IProject } from 'src/app/api.service';
import { PrettyTimeService } from 'src/app/pretty-time.service';
import { IconDefinition } from '@fortawesome/fontawesome-svg-core';
import {
  faCheckCircle,
  faExclamationTriangle,
  faExclamationCircle,
  faQuestionCircle
} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-project-list-item',
  templateUrl: './project-list-item.component.html',
  styleUrls: ['./project-list-item.component.scss']
})
export class ProjectListItemComponent {
  constructor(private time: PrettyTimeService) {
  }

  @Input() project: IProject;

  get backupStatusText(): string {
    let s = '';

    switch (this.project.backupStatus) {
      case 'ok':
        s = 'OK';
        break;
      case 'none':
        return 'No backup exists';

      case 'outdated':
        s = 'Out of date';
        break;

      default:
        return this.project.backupStatus;
    }

    if (this.project.lastBackup) {
      s += `, backup created ${this.time.formatRelative(this.project.lastBackup.time)}`
    }

    return s;
  }

  get backupStatusIcon(): IconDefinition {
    switch (this.project.backupStatus) {
      case 'ok':
        return faCheckCircle;

      case 'none':
        return faExclamationTriangle;

      case 'outdated':
        return faExclamationCircle;

      default:
        return faQuestionCircle;
    }
  }

  get description(): string {
    let description = '';
    if (!this.project.isActive) {
      description += 'Project is disabled .';
    }

    if (!this.project.notifications) {
      description += 'Notifications are disabled. ';
    }

    description += `Backups are expected to be taken every ${this.time.formatDuration(this.project.backupFrequency * 1000)}. `;
    return description;
  }

  get className(): string {
    if (!this.project.isActive) {
      return 'list-group-item-light';
    }

    switch (this.project.backupStatus) {
      case 'ok':
        return 'list-group-item-success';
      case 'none':
        return 'list-group-item-warning';
      case 'outdated':
        return 'list-group-item-warning';
    }

    return '';
  }
}
