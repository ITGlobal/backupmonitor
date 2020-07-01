import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { map, tap, catchError } from 'rxjs/operators';
import { Observable, throwError } from 'rxjs';

export interface IUser {
  id: number;
  username: string;
}

export type BackupType = 'last' | 'archive';

export interface IBackup {
  id: string;
  filename: string;
  time: Date;
  type: BackupType;
}
export type BackupStatus = 'ok' | 'outdated' | 'none';

export interface INotificationParams {
  enabled: boolean;
  slack: string[];
  telegram: string[];
  webhook: string[];
}

export interface IProject {
  id: string;
  name: string;
  isActive: boolean;
  backupFrequency: number;
  backupRetention: number;
  notifications: INotificationParams;
  lastBackup?: IBackup;
  backupStatus: BackupStatus;
}

export interface IProjectCreateParams {
  id: string;
  name: string;
  isActive: boolean;
  backupFrequency: number;
  backupRetention: number;
  notifications: INotificationParams;
}

export interface IProjectUpdateParams {
  name?: string;
  isActive?: boolean;
  backupFrequency: number;
  backupRetention: number;
  notifications: INotificationParams;
}

export interface IAccessKey {
  id: number;
  label: string;
  key: string;
}

interface IAuthResponse {
  token: string;
  user: IUser;
}

const localStorageKeys = {
  token: 'api_token',
  user: 'api_user'
};

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private token: string | null;
  private user: IUser | null;

  constructor(private http: HttpClient) {
    this.token = localStorage.getItem(localStorageKeys.token);

    const j = localStorage.getItem(localStorageKeys.user);
    this.user = j && JSON.parse(j) || null;
  }

  public isAuthorized(): boolean {
    return !!this.token;
  }

  public getUser(): IUser {
    return this.user as IUser;
  }

  public authorize(username: string, password: string): Observable<void> {
    return this.http.post<IAuthResponse>('/api/authorize', { username, password })
      .pipe(
        catchError(ApiService.handleError)
      )
      .pipe(
        tap(r => {
          this.token = r.token;
          this.user = r.user;

          localStorage.setItem(localStorageKeys.token, this.token);
          localStorage.setItem(localStorageKeys.user, JSON.stringify(this.user));
        })
      )
      .pipe(
        map(_ => { })
      );
  }

  public unauthorize() {
    this.token = null;
    this.user = null;

    localStorage.removeItem(localStorageKeys.token);
    localStorage.removeItem(localStorageKeys.user);
  }

  public changePassword(oldPassword: string, newPassword: string): Observable<{}> {
    return this.http.post<{}>(`/api/me/password`, { oldPassword, newPassword }, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public getProjects(): Observable<IProject[]> {
    return this.http.get<IProject[]>('/api/projects', {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      )
      .pipe(
        map((xs) => {
          for (const i in xs) {
            xs[i] = ApiService.mapProject(xs[i]);
          }
          return xs;
        })
      );
  }

  public getProject(id: string): Observable<IProject> {
    return this.http.get<IProject>(`/api/projects/${id}`, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      )
      .pipe(
        map((x) => {
          x = ApiService.mapProject(x);
          return x;
        })
      );
  }

  public createProject(project: IProjectCreateParams): Observable<IProject> {
    return this.http.post<IProject>(`/api/projects`, project, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public updateProject(id: string, project: IProjectUpdateParams): Observable<IProject> {
    return this.http.put<IProject>(`/api/projects/${id}`, project, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public deleteProject(projectId: string): Observable<void> {
    return this.http.delete<void>(`/api/projects/${projectId}`, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public getProjectBackups(id: string): Observable<IBackup[]> {
    return this.http.get<IBackup[]>(`/api/projects/${id}/backup`, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      )
      .pipe(
        map((xs) => {
          for (const i in xs) {
            xs[i] = ApiService.mapBackup(xs[i]);
          }
          return xs;
        })
      );
  }

  public deleteProjectBackup(backupId: string): Observable<void> {
    return this.http.delete<void>(`/api/backup/${backupId}`, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public getProjectAccessKeys(id: string): Observable<IAccessKey[]> {
    return this.http.get<IAccessKey[]>(`/api/projects/${id}/keys`, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public createProjectAccessKey(projectId: string, label: string): Observable<IAccessKey> {
    return this.http.post<IAccessKey>(`/api/projects/${projectId}/keys`, { label }, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public deleteProjectAccessKey(projectId: string, accessKeyId: number): Observable<void> {
    return this.http.delete<void>(`/api/projects/${projectId}/keys/${accessKeyId}`, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public testSlackNotification(target: string): Observable<void> {
    return this.http.post<void>('/api/notify/slack', { target }, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public testTelegramNotification(target: string): Observable<void> {
    return this.http.post<void>('/api/notify/telegram', { target }, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  public testWebhookNotification(target: string): Observable<void> {
    return this.http.post<void>('/api/notify/webhook', { target }, {
      headers: {
        Authorization: `Bearer ${this.token}`
      }
    })
      .pipe(
        catchError(ApiService.handleError)
      );
  }

  private static handleError(error: HttpErrorResponse) {
    if (error.error?.message) {
      return throwError(error.error?.message);
    }
    if (error.message) {
      return throwError(error.message);
    }

    return throwError(error.status);
  }

  private static mapProject(obj: any): IProject {
    if (!obj) {
      return obj;
    }

    if (obj.lastBackup) {
      obj.lastBackup = ApiService.mapBackup(obj.lastBackup);
    }

    return obj as IProject;
  }

  private static mapBackup(obj: any): IBackup {
    if (!obj) {
      return obj;
    }

    if (obj.time) {
      obj.time = new Date(Date.parse(obj.time as string));
    }

    return obj as IBackup;
  }
}

