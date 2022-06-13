import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { User } from 'src/app/model/user';
import { logoutAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import {
  currentUserSelector,
  isAnonymousSelector,
  isLoggedInSelector,
} from 'src/app/store/selectors/auth.selectors';

@Component({
  selector: 'app-top-bar',
  templateUrl: './top-bar.component.html',
  styleUrls: ['./top-bar.component.scss'],
})
export class TopBarComponent implements OnInit {
  isLoggedIn$!: Observable<boolean | null>;
  isAnonumous$!: Observable<boolean>;
  currentUser$!: Observable<User | null>;

  constructor(private store: Store<AppState>) {}

  ngOnInit(): void {
    this.initializeValues();
  }

  initializeValues() {
    this.isLoggedIn$ = this.store.select(isLoggedInSelector);
    this.isAnonumous$ = this.store.select(isAnonymousSelector);
    this.currentUser$ = this.store.select(currentUserSelector);
  }

  doLogout() {
    this.store.dispatch(logoutAction());
  }
}
