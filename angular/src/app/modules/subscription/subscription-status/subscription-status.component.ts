import { Component, OnInit } from '@angular/core';
import { Store, select } from '@ngrx/store';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Subscription } from 'src/app/models/subscription';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectSubscription } from 'src/app/store/selectors/subscription.selectors';

@Component({
  selector: 'app-subscription-status',
  templateUrl: './subscription-status.component.html',
  styleUrls: ['./subscription-status.component.scss']
})
export class SubscriptionStatusComponent implements OnInit {
  // Observable for the subscription data from the store
  subscriptionData$: Observable<{ subscription: Subscription | null; isLoading: boolean }> | undefined;

  constructor(private store: Store<AppState>) { }

  ngOnInit(): void {
    // Select the subscription data from the store
    this.subscriptionData$ = this.store.pipe(select(selectSubscription));
  }
}
