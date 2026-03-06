import { Component } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-subscription-success',
  templateUrl: './subscription-success.component.html',
  styleUrl: './subscription-success.component.scss',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, RouterModule],
})
export class SubscriptionSuccessComponent {}
