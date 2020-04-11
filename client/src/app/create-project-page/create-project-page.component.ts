import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { IProjectCreateParams, ApiService } from '../api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-create-project-page',
  templateUrl: './create-project-page.component.html',
  styleUrls: ['./create-project-page.component.scss']
})
export class CreateProjectPageComponent implements OnInit {

  constructor(private api: ApiService, private router: Router) {
    this.project = {
      id: '',
      name: '',
      isActive: true,
      notifications: {
        enabled: false,
        slack: [],
        telegram: [],
        webhook: []
      },
      backupFrequency: 24 * 3600,
      backupRetention: 10,
    };
  }

  project: IProjectCreateParams;
  form: FormGroup;

  isBusy: boolean;
  error?: string;

  get id() { return this.form.get('id'); }
  get name() { return this.form.get('name'); }
  get backupFrequency() { return this.form.get('backupFrequency'); }
  get backupRetention() { return this.form.get('backupRetention'); }
  get isActive() { return this.form.get('isActive'); }
  get notifications() { return this.form.get('notifications'); }

  ngOnInit(): void {
    this.form = new FormGroup({
      'id': new FormControl(this.project.id, [
        Validators.required,
        Validators.minLength(2),
        Validators.pattern(/^[a-z][a-z0-9-]*[a-z0-9]$/)
      ]),
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

  }

  onSubmit() {
    if (this.isBusy) {
      return;
    }
    
    this.isBusy = true;
    this.error = undefined;

    const model: IProjectCreateParams = this.form.value;
    model.backupFrequency = parseInt(model.backupFrequency as any);
    model.backupRetention = parseInt(model.backupRetention as any);

    this.api.createProject(model)
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
    this.router.navigate(['/projects']);
  }
}
