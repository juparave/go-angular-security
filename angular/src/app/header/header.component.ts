import { Component, OnDestroy, OnInit } from '@angular/core';
import { AuthService } from 'src/app/services/auth/auth.service';
import { Router } from '@angular/router';
import { AppState } from 'src/app/store/interfaces/app-state';
import {
  currentUserSelector,
  isAnonymousSelector,
  isLoggedInSelector,
} from '../store/selectors/auth.selectors';
import { User } from 'src/app/model/user';
import { select, Store } from '@ngrx/store';
import { Observable, Subscription, tap } from 'rxjs';
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
  ) {}

  ngOnInit(): void {
    this.isLoggedIn$ = this.store.select(isLoggedInSelector);
    this.isAnonymous$ = this.store.select(isAnonymousSelector);
    this.currentUser$ = this.store.select(currentUserSelector);
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
