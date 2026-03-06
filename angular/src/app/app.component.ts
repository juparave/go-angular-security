import { Component } from '@angular/core';
import { Store } from '@ngrx/store';
import { getCurrentUserAction } from './store/actions/auth.actions';
import { getSubscriptionAction } from './store/actions/subscription.actions';

@Component({
  selector: 'app-root',
  standalone: false,
  template: `<router-outlet></router-outlet>`,
})
export class AppComponent {
  constructor(private store: Store) { }

  ngOnInit() {
    // Fetch current user data on app start
    this.store.dispatch(getCurrentUserAction());

    // Subscription data will be fetched by the AuthGetCurrentUserEffects
    // after the user is successfully retrieved.
  }
}
