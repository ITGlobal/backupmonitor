import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { HighlightModule, HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';
import { NgbModalModule, NgbDropdownModule, NgbCollapseModule } from '@ng-bootstrap/ng-bootstrap';
import { ClipboardModule } from 'ngx-clipboard';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginPageComponent } from './login-page/login-page.component';
import { NavBarComponent } from './controls/navbar/navbar.component';
import { LogoutPageComponent } from './logout-page/logout-page.component';
import { ProjectPageComponent } from './project-page/project-page.component';
import { ProjectsPageComponent } from './projects-page/projects-page.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ProjectintegrationGuideComponent } from './project-page/project-integration-guide/project-integration-guide.component';
import { ProjectSummaryComponent } from './project-page/project-summary/project-summary.component';
import { ProjectBackupsComponent } from './project-page/project-backups/project-backups.component';
import { ProjectAccessKeysComponent } from './project-page/project-access-keys/project-access-keys.component';
import { ViewAccessKeyModalComponent } from './modals/view-access-key-modal/view-access-key-modal.component';
import { DeleteAccessKeyModalComponent } from './modals/delete-access-key-modal/delete-access-key-modal.component';
import { CreateAccessKeyModalComponent } from './modals/create-access-key-modal/create-access-key-modal.component';
import { DeleteBackupModalComponent } from './modals/delete-backup-modal/delete-backup-modal.component';
import { DeleteProjectModalComponent } from './modals/delete-project-modal/delete-project-modal.component';
import { EditProjectPageComponent } from './edit-project-page/edit-project-page.component';
import { CreateProjectPageComponent } from './create-project-page/create-project-page.component';
import { NotificationTargetsEditorComponent } from './controls/notification-targets-editor/notification-targets-editor.component';
import { AddNotificationTargetModalComponent } from './modals/add-notification-target-modal/add-notification-target-modal.component';
import { ProjectListItemComponent } from './projects-page/project-list-item/project-list-item.component';
import { ChangePasswordModalComponent } from './modals/change-password-modal/change-password-modal.component';

export function getHighlightLanguages() {
  return {
    bash: () => import('highlight.js/lib/languages/bash'),
    http: () => import('highlight.js/lib/languages/http'),
  };
}

@NgModule({
  declarations: [
    AppComponent,
    LoginPageComponent,
    NavBarComponent,
    LogoutPageComponent,
    ProjectPageComponent,
    ProjectsPageComponent,
    PageNotFoundComponent,
    ProjectintegrationGuideComponent,
    ProjectSummaryComponent,
    ProjectBackupsComponent,
    ProjectAccessKeysComponent,
    ViewAccessKeyModalComponent,
    DeleteAccessKeyModalComponent,
    CreateAccessKeyModalComponent,
    DeleteBackupModalComponent,
    DeleteProjectModalComponent,
    EditProjectPageComponent,
    CreateProjectPageComponent,
    NotificationTargetsEditorComponent,
    AddNotificationTargetModalComponent,
    ProjectListItemComponent,
    ChangePasswordModalComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    AppRoutingModule,
    FontAwesomeModule,
    HttpClientModule,
    FormsModule,
    HighlightModule,
    NgbModalModule,
    NgbDropdownModule,
    NgbCollapseModule,
    ClipboardModule,
  ],
  providers: [
    {
      provide: HIGHLIGHT_OPTIONS,
      useValue: {
        languages: getHighlightLanguages()
      }
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
