import { Component, OnInit, computed, effect, inject, signal } from '@angular/core';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { Observable } from 'rxjs';
import { map, shareReplay } from 'rxjs/operators';
import { Store } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';
import { User } from 'src/app/models/user';
import { selectCurrentUser } from 'src/app/store/selectors/auth.selectors';
import { logoutAction } from 'src/app/store/actions/auth.actions';

@Component({
  selector: 'app-base',
  standalone: false,
  templateUrl: './base.component.html',
  styleUrls: ['./base.component.scss']
})
export class BaseComponent implements OnInit {
  private breakpointObserver = inject(BreakpointObserver);
  currentUser$!: Observable<User | null>;

  isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset)
    .pipe(
      map(result => result.matches),
      shareReplay()
    );

  collapsed = signal<boolean>(false);
  isShowing = signal<boolean>(false);

  sidenavWidth = computed(() => (this.collapsed()) ? '64px' : '250px');

  constructor(private store: Store<AppState>) {
    effect(() => {
      if (this.isShowing()) {
        this.collapsed.set(false);
      }
    }, { allowSignalWrites: true });
  }

  ngOnInit(): void {
    this.initializeValues();
  }

  initializeValues() {
    this.currentUser$ = this.store.select(selectCurrentUser)
  }

  closeSidenav() {
    // receive the close event from the child component
    this.isShowing.set(false);
  }

  doLogout() {
    this.store.dispatch(logoutAction());
  }
}
