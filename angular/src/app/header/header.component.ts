import { Component, OnDestroy, OnInit } from '@angular/core';
import { AuthService } from 'src/app/services/auth/auth.service';
import { Router } from '@angular/router';
import { AppState } from 'src/app/store/interfaces/app-state';
import {
  selectCurrentUser,
  selectIsAnonymous,
  selectIsLoggedIn,
} from '../store/selectors/auth.selectors';
import { User } from 'src/app/model/user';
import { select, Store } from '@ngrx/store';
import { Observable, Subscription, map, tap } from 'rxjs';
import { logoutAction } from '../store/actions/auth.actions';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnInit, OnDestroy {
  public username = '';
  private subscription?: Subscription;
  isLoggedIn$!: Observable<boolean | null>;
  isAnonymous$!: Observable<boolean>;
  currentUser$!: Observable<User | null>;

  constructor(
    private store: Store<AppState>,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    this.isLoggedIn$ = this.store.select(selectIsLoggedIn).pipe(map(state => state?.isLoggedIn));
    this.isAnonymous$ = this.store.select(selectIsAnonymous).pipe(map(state => state?.isAnon));
    this.currentUser$ = this.store.select(selectCurrentUser);
    // this.subscription = this.store
    //   .pipe(
    //     select(currentUserSelector),
    //     tap((currentUser: User | null) => {
    //       if (currentUser != null && currentUser.id != null) {
    //         this.username = currentUser.email;
    //       } else {
    //         this.username = '';
    //       }
    //     })
    //   )
    //   .subscribe();
  }

  logout(): void {
    this.store.dispatch(logoutAction());
    // this.authService.logout().subscribe();
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }
}
