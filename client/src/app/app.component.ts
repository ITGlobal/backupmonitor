import { Component } from '@angular/core';
import { FaConfig, FaIconLibrary } from '@fortawesome/angular-fontawesome';
import { fas, prefix } from '@fortawesome/free-solid-svg-icons';
import { fab } from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styles: []
})
export class AppComponent {
  title = 'backup-monitor-frontend';

  constructor(library: FaIconLibrary, faConfig: FaConfig) {
    library.addIconPacks(fas, fab);
    faConfig.defaultPrefix = prefix;
  }
}
