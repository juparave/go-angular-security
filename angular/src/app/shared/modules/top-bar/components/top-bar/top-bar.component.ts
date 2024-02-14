import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable, map } from 'rxjs';
import { User } from 'src/app/model/user';
import { logoutAction } from 'src/app/store/actions/auth.actions';
import { AppState } from 'src/app/store/interfaces/app-state';
import {
  selectCurrentUser,
  selectIsAnonymous,
  selectIsLoggedIn,
} from 'src/app/store/selectors/auth.selectors';

@Component({
  selector: 'app-top-bar',
  templateUrl: './top-bar.component.html',
  styleUrls: ['./top-bar.component.scss'],
})
export class TopBarComponent implements OnInit {
  isLoggedIn$!: Observable<boolean | null>;
  isAnonymous$!: Observable<boolean>;
  currentUser$!: Observable<User | null>;

  constructor(private store: Store<AppState>) { }

  ngOnInit(): void {
    this.initializeValues();
  }

  initializeValues() {
    this.isLoggedIn$ = this.store.select(selectIsLoggedIn).pipe(map(state => state?.isLoggedIn));
    this.isAnonymous$ = this.store.select(selectIsAnonymous).pipe(map(state => state?.isAnon));
    this.currentUser$ = this.store.select(selectCurrentUser);
  }

  doLogout() {
    this.store.dispatch(logoutAction());
  }
}
