import { Component, OnInit } from '@angular/core';
import { ApiService } from '../api.service';
import { FaConfig } from '@fortawesome/angular-fontawesome';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.scss'],
  providers: [FaConfig],
})
export class LoginPageComponent implements OnInit {

  constructor(private api: ApiService, private router: Router) {
    this.isUsernameValid = true;
    this.isPasswordValid = true;
  }

  username: string;
  password: string;

  isUsernameValid: boolean;
  isPasswordValid: boolean;

  isBusy: boolean;
  error?: string;

  ngOnInit(): void {
  }

  onSubmit() {
    if (this.isBusy) {
      return;
    }

    this.isUsernameValid = true;
    this.isPasswordValid = true;

    if (!this.username) {
      this.isUsernameValid = false;
    }

    if (!this.password) {
      this.isPasswordValid = false;
    }

    if (!this.isUsernameValid || !this.isPasswordValid) {
      return;
    }

    this.isBusy = true;
    this.error = undefined;

    this.api.authorize(this.username, this.password)
      .subscribe(
        () => {
          this.isBusy = false;

          const u = this.router.parseUrl(this.router.url);
          const url = u.queryParams['returnUrl'] || '/';
          this.router.navigate([url])
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
