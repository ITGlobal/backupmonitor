import { Component, OnInit } from '@angular/core';
import { ApiService, IProject, IProjectUpdateParams, IProjectCreateParams } from '../api.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormGroup, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-edit-project-page',
  templateUrl: './edit-project-page.component.html',
  styleUrls: ['./edit-project-page.component.scss']
})
export class EditProjectPageComponent implements OnInit {

  constructor(private api: ApiService, private route: ActivatedRoute, private router: Router) { }

  id: string;
  project: IProject;
  form: FormGroup;

  isBusy: boolean;
  error?: string;

  get name() { return this.form.get('name'); }
  get backupFrequency() { return this.form.get('backupFrequency'); }
  get backupRetention() { return this.form.get('backupRetention'); }
  get isActive() { return this.form.get('isActive'); }
  get notifications() { return this.form.get('notifications'); }

  ngOnInit(): void {
    this.isBusy = true;

    this.route.paramMap.subscribe(
      (p) => {
        const id = p.get('id');
        if (!id) {
          this.router.navigate(['/projects']);
          return;
        }

        this.id = id;
        this.load();
      }
    )
  }

  private load() {
    this.isBusy = true;
    this.error = undefined;

    this.api.getProject(this.id).subscribe(
      (project) => {
        this.project = project;
        this.isBusy = false;

        this.form = new FormGroup({
          'name': new FormControl(this.project.name, [
            Validators.required,
            Validators.minLength(1),
            Validators.pattern(/^.+.*$/)
          ]),
          'backupFrequency': new FormControl(this.project.backupFrequency, [
            Validators.required,
            Validators.min(3600),
          ]),
          'backupRetention': new FormControl(this.project.backupRetention, [
            Validators.required,
            Validators.min(1),
          ]),
          'isActive': new FormControl(this.project.isActive),
          'notifications': new FormControl(this.project.notifications),
        });
      },
      (e) => {
        this.isBusy = false;
        this.error = e;
      });
  }

  onSubmit() {
    if (this.isBusy) {
      return;
    }

    this.isBusy = true;
    this.error = undefined;

    const model: IProjectUpdateParams = this.form.value;
    model.backupFrequency = parseInt(model.backupFrequency as any);
    model.backupRetention = parseInt(model.backupRetention as any);

    const e = this.validate(model);
    if (!!e) {
      this.isBusy = false;
      this.error = e;
      return;
    }

    this.api.updateProject(this.id, model)
      .subscribe(
        (p) => {
          this.isBusy = false;
          this.router.navigate(['/projects', p.id]);
        },
        (e) => {
          this.isBusy = false;
          this.error = e;
        });
  }

  cancel() {
    this.router.navigate(['/projects', this.id]);
  }

  validate(model: IProjectUpdateParams): string | null {
    if (!model.name) {
      return 'Name is not set';
    }

    return null;
  }
}
