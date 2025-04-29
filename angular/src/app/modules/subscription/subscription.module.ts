import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms'; // Import FormsModule for ngModel

// Angular Material Modules
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';

import { SubscriptionRoutingModule } from './subscription-routing.module';
import { SubscriptionPageComponent } from './subscription-page/subscription-page.component';
import { SubscriptionStatusComponent } from './subscription-status/subscription-status.component';


@NgModule({
  declarations: [
    SubscriptionPageComponent,
    SubscriptionStatusComponent
  ],
  imports: [
    CommonModule,
    FormsModule, // Add FormsModule
    SubscriptionRoutingModule,
    // Add Angular Material Modules
    MatCardModule,
    MatButtonModule,
    MatSlideToggleModule,
    MatListModule,
    MatIconModule
  ]
})
export class SubscriptionModule { }
