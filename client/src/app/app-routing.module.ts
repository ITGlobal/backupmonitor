import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginPageComponent } from './login-page/login-page.component';
import { AuthGuard } from './auth.guard';
import { LogoutPageComponent } from './logout-page/logout-page.component';
import { ProjectPageComponent } from './project-page/project-page.component';
import { ProjectsPageComponent } from './projects-page/projects-page.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { CreateProjectPageComponent } from './create-project-page/create-project-page.component';
import { EditProjectPageComponent } from './edit-project-page/edit-project-page.component';

const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'logout', component: LogoutPageComponent, canActivate: [AuthGuard] },
  { path: 'projects', component: ProjectsPageComponent, canActivate: [AuthGuard] },
  { path: 'new-project', component: CreateProjectPageComponent, canActivate: [AuthGuard] },
  { path: 'projects/:id/edit', component: EditProjectPageComponent, canActivate: [AuthGuard] },
  { path: 'projects/:id', component: ProjectPageComponent, canActivate: [AuthGuard] },
  { path: '', redirectTo: '/projects', canActivate: [AuthGuard], pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
