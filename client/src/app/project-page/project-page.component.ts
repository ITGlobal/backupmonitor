import { Component, OnInit } from '@angular/core';
import { IProject, ApiService, IBackup, IAccessKey } from '../api.service';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-project-page',
  templateUrl: './project-page.component.html',
  styleUrls: ['./project-page.component.scss']
})
export class ProjectPageComponent implements OnInit {

  constructor(private api: ApiService, private route: ActivatedRoute, private router: Router) {
    this.id = '';
    this.activeTab = 'summary';
    this.backups = [];
    this.accessKeys = [];
  }

  id: string;
  isBusy: boolean;
  project?: IProject;
  backups: IBackup[];
  accessKeys: IAccessKey[];
  error?: string;
  activeTab: string;

  ngOnInit(): void {
    this.route.queryParamMap.subscribe((p) => {
      const tab = p.get('tab');
      this.activeTab = tab || 'summary';
    });

    this.route.paramMap.subscribe(
      (p) => {
        const id = p.get('id');
        if (!id) {
          this.router.navigate(['/projects']);
          return;
        }

        this.id = id;
        this.refresh();
      }
    )
  }

  refresh() {
    this.isBusy = true;
    this.error = undefined;

    this.api.getProject(this.id).subscribe(
      (project) => {
        this.project = project;

        this.api.getProjectBackups(this.id).subscribe(
          (backups) => {
            this.backups = backups;
            this.api.getProjectAccessKeys(this.id).subscribe(
              (accessKeys) => {
                this.accessKeys = accessKeys;
                this.isBusy = false;
              },
              (e) => {
                this.isBusy = false;
                this.error = e;
              });
          },
          (e) => {
            this.isBusy = false;
            this.error = e;
          });
      },
      (e) => {
        this.isBusy = false;
        this.error = e;
      });
  }

  refreshBackups() {
    this.isBusy = true;
    this.error = undefined;

    this.api.getProjectBackups(this.id).subscribe(
      (backups) => {
        this.backups = backups;
        this.isBusy = false;
      },
      (e) => {
        this.isBusy = false;
        this.error = e;
      });
  }

  refreshAccessKeys() {
    this.isBusy = true;
    this.error = undefined;

    this.api.getProjectAccessKeys(this.id).subscribe(
      (accessKeys) => {
        this.accessKeys = accessKeys;
        this.isBusy = false;
      },
      (e) => {
        this.isBusy = false;
        this.error = e;
      });
  }

  selectTab(tab: string) {
    const url = this.router.parseUrl(this.router.url);
    this.router.navigate(url.root.segments, {
      queryParams: {
        tab: tab
      }
    });
    return false;
  }

  getTabLinkClass(tab: string): string {
    return this.activeTab === tab ? 'active' : '';
  }

  getTabPaneClass(tab: string): string {
    return this.activeTab === tab ? 'show active' : '';
  }
}
