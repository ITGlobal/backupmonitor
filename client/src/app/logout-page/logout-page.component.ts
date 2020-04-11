import { Component, OnInit } from '@angular/core';
import { ApiService } from '../api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-logout-page',
  template: ''
})
export class LogoutPageComponent implements OnInit {

  constructor(private api: ApiService, private router: Router) { }

  ngOnInit(): void {
    this.api.unauthorize();
    this.router.navigate(['/']);
  }

}
