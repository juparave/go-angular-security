import { Component } from '@angular/core';
import { Store } from '@ngrx/store';
import { getCurrentUserAction } from './store/actions/auth.actions';

@Component({
  selector: 'app-root',
  standalone: false,
  template: `<router-outlet></router-outlet>`,
})
export class AppComponent {
  title = 'Austral QC Dashboard';

  constructor(private store: Store) { }

  ngOnInit() {
    this.store.dispatch(getCurrentUserAction());
  }
}
