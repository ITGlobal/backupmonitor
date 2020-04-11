import { Component, Input } from '@angular/core';
import { IAccessKey } from 'src/app/api.service';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { ClipboardService } from 'ngx-clipboard'

@Component({
  selector: 'app-view-access-key-modal',
  templateUrl: './view-access-key-modal.component.html',
  styleUrls: ['./view-access-key-modal.component.scss']
})
export class ViewAccessKeyModalComponent {
  constructor(private modal: NgbActiveModal, private clipboard: ClipboardService) { }

  @Input() public accessKey: IAccessKey;
  copyStatus?: string;

  copy() {
    this.clipboard.copy(this.accessKey.key);
    this.copyStatus = 'Copied to clipboard';
    window.setTimeout(() => { this.copyStatus = undefined }, 1000);
  }

  dismiss() {
    this.modal.dismiss();
  }
}
