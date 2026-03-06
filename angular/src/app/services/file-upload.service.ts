import { Injectable } from '@angular/core';
import { HttpClient, HttpRequest, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class FileUploadService {
  private baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  upload(files: File[]): Observable<HttpEvent<any>> {
    const formData: FormData = new FormData();

    files.forEach((file: File) => {
      formData.append('files', file, file.name);
    });

    const req = new HttpRequest('POST', `${this.baseUrl}/upload/`, formData, {
      reportProgress: true,
      responseType: 'json'
    });

    return this.http.request(req);
  }

  getFiles(): Observable<any> {
    return this.http.get(`${this.baseUrl}/upload/`);
  }
}

