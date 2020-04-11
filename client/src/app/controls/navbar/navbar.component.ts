import { Component } from '@angular/core';
import { ApiService } from '../../api.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { ChangePasswordModalComponent } from 'src/app/modals/change-password-modal/change-password-modal.component';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html'
})
export class NavBarComponent {
  username: string;
  isMenuCollapsed: boolean = true;

  constructor(api: ApiService, private modal: NgbModal) {
    this.username = api.getUser()?.username;
  }

  toggleNavBar() {
    this.isMenuCollapsed = !this.isMenuCollapsed;
  }

  changePassword() {
    this.onItemClicked();
    this.modal.open(ChangePasswordModalComponent);
    return false;
  }

  onItemClicked() {
    this.isMenuCollapsed = true;
  }
}
