import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { SubscriptionPageComponent } from './subscription-page/subscription-page.component';
import { SubscriptionRoutingModule } from './subscription-routing.module';


@NgModule({
  declarations: [
    SubscriptionPageComponent
  ],
  imports: [
    CommonModule,
    SubscriptionRoutingModule
  ]
})
export class SubscriptionModule { }
