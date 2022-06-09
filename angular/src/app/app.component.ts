import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { getCurrentUserAction } from './store/actions/auth.actions';

@Component({
  selector: 'app-root',
  template: `<router-outlet></router-outlet>`,
})
export class AppComponent implements OnInit {
  title = 'angular';
  constructor(private store: Store) {}

  ngOnInit() {
    this.store.dispatch(getCurrentUserAction());
  }
}
