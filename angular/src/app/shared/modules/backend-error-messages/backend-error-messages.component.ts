import { Component, Input, OnInit } from '@angular/core';
import { BackendErrors } from 'src/app/store/types/backend-errors.interface';

@Component({
  selector: 'app-backend-error-messages',
  standalone: false,
  styles: [],
  templateUrl: './backend-error-messages.component.html',
})
export class BackendErrorMessagesComponent implements OnInit {
  @Input('backendErrors') backendErrors!: BackendErrors | null;

  errorMessages: string[] = [];

  constructor() { }

  ngOnInit(): void {
    if (this.backendErrors) {
      if (typeof this.backendErrors['error'] === 'string') {
        // Handle single string error message
        this.errorMessages = [this.backendErrors['error']];
      } else {
        // Handle object of arrays error message
        this.errorMessages = Object.keys(this.backendErrors).flatMap(
          (name: string) => {
            const messages = this.backendErrors![name];
            return messages.map(message => `${name}: ${message}`);
          }
        );
      }
    } else {
      this.errorMessages = [];
    }
  }
}
