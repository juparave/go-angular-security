import { Component, Input, OnInit } from '@angular/core';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';

@Component({
  selector: 'app-backend-error-messages',
  styles: [],
  templateUrl: './backend-error-messages.component.html',
})
export class BackendErrorMessagesComponent implements OnInit {
  @Input('backendErrors') backendErrors!: BackendErrors | null;

  errorMessages: string[] = [];

  constructor() {}

  ngOnInit(): void {
    if (this.backendErrors) {
      this.errorMessages = Object.keys(this.backendErrors).map(
        (name: string) => {
          const messages = this.backendErrors![name].join(' ');
          return `${name}: ${messages}`;
        }
      );
    } else {
      this.errorMessages = [];
    }
  }
}
