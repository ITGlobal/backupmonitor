import { Component, OnInit } from '@angular/core';
import { ApiService, IProject } from '../api.service';

@Component({
  selector: 'app-main-page',
  templateUrl: './projects-page.component.html',
  styleUrls: ['./projects-page.component.scss']
})
export class ProjectsPageComponent implements OnInit {
  constructor(private api: ApiService) {
    this.projects = [];
  }

  isBusy: boolean;
  projects: IProject[];
  error?: string;

  ngOnInit(): void {
    this.refresh()
  }

  refresh() {
    this.isBusy = true;
    this.error = undefined;

    this.api.getProjects().subscribe(
      (projects) => {
        this.projects = projects;
        this.isBusy = false;
      },
      (e) => {
        this.isBusy = false;
        this.error = e;
      });
  }

  dismissError() {
    this.error = undefined;
  }
}
