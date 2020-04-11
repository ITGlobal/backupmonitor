import { Component, OnInit, Input } from '@angular/core';
import { IProject } from 'src/app/api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-project-integration-guide',
  templateUrl: './project-integration-guide.component.html',
  styleUrls: ['./project-integration-guide.component.scss']
})
export class ProjectintegrationGuideComponent {
  constructor(private router: Router) {
  }

  @Input() project?: IProject;

  getBackupEndpoint(): string {
    return `${window.location.protocol}//${window.location.host}/api/backup`;
  }

  getAccessKeysPageUrl(): string {
    return `/projects/${this.project?.id}?tab=keys`;
  }

  openAccessKeysPage() {
    this.router.navigate(['/projects', this.project?.id], {
      queryParams: {
        tab: 'keys'
      }
    })
    return false;
  }

  getShortCodeSample(): string {
    return `curl -X POST http://...`;
  }

  getGenericCodeSample(): string {
    return `curl -X POST "${this.getBackupEndpoint()}" -F "file=@my-awesome-backup.zip"`;
  }

  getCodeSampleWithHeaderAuth(): string {
    return `curl -X POST "${this.getBackupEndpoint()}" \\\n` +
      `     -H "Authorization: $YOUR_ACCESS_KEY" \\\n` +
      `     -F "file=@my-awesome-backup.zip"`;
  }

  getCodeSampleWithQueryAuth(): string {
    return `curl -X POST "${this.getBackupEndpoint()}?key=$YOUR_ACCESS_KEY" \\\n` +
      `     -F "file=@my-awesome-backup.zip"`;
  }

  get201HttpResponse(): string {
    return "" +
      `HTTP/1.1 201 Created\n` +
      `Location: ${this.getBackupEndpoint()}/0123456789abcdef\n` +
      `Content-Type: application/json\n` +
      `Content-Length: 134\n\n` +
      `{\n` +
      `   "id" : "0123456789abcdef",\n` +
      `   "filename" : "my-awesome-backup.zip",\n` +
      `   "type" : "latest",\n` +
      `   "time" : "2020-01-01T12:00:00Z"\n` +
      `}`;
  }

  get400HttpResponse(): string {
    return "" +
      `HTTP/1.1 400 Bad Request\n` +
      `Content-Type: application/json\n` +
      `Content-Length: 74\n\n` +
      `{\n` +
      `    "error" : "bad_request",\n` +
      `    "message" : "not a multipart form"\n` +
      `}`;
  }

  get403HttpResponse(): string {
    return "" +
      `HTTP/1.1 403 Forbidden\n` +
      `Content-Type: application/json\n` +
      `Content-Length: 69\n\n` +
      `{\n` +
      `    "error" : "access_denied",\n` +
      `    "message" : "access denied"\n` +
      `}`;
  }
}
