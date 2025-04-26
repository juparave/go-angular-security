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

    // Also fetch subscription data
    // This typically would be done after the user is authenticated
    // but we dispatch it here for simplicity
    this.store.dispatch(getSubscriptionAction());
  }
}
