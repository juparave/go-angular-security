import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { Store, select } from '@ngrx/store';
import { Observable } from 'rxjs';
import { Subscription } from 'src/app/models/subscription';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectSubscription } from 'src/app/store/selectors/subscription.selectors';

@Component({
  selector: 'app-subscription-status',
  templateUrl: './subscription-status.component.html',
  styleUrls: ['./subscription-status.component.scss'],
  standalone: true,
  imports: [CommonModule, MatCardModule, MatIconModule, MatProgressSpinnerModule],
})
export class SubscriptionStatusComponent implements OnInit {
  private store = inject<Store<AppState>>(Store);

  // Observable for the subscription data from the store
  subscriptionData$:
    | Observable<{ subscription: Subscription | null; isLoading: boolean }>
    | undefined;

  ngOnInit(): void {
    // Select the subscription data from the store
    this.subscriptionData$ = this.store.pipe(select(selectSubscription));
  }
}
