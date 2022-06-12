import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BackendErrorMessagesComponent } from './backend-error-messages.component';
import { AlertModule } from 'src/app/shared/modules/alert/alert.module';

@NgModule({
  declarations: [BackendErrorMessagesComponent],
  imports: [CommonModule, AlertModule],
  exports: [BackendErrorMessagesComponent],
})
export class BackendErrorMessagesModule {}
