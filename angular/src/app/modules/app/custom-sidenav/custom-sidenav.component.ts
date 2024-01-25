import { Component, EventEmitter, Input, OnInit, Output, computed, signal } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable } from 'rxjs';
import { User } from 'src/app/model/user';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectCurrentUser } from '@store/selectors/auth.selectors';
import { logoutAction } from '@store/actions/auth.actions';

export type MenuItem = {
  icon: string;
  label: string;
  route?: string;
}


@Component({
  selector: 'app-custom-sidenav',
  templateUrl: './custom-sidenav.component.html',
  styleUrl: './custom-sidenav.component.scss'
})
export class CustomSidenavComponent implements OnInit {
  @Output() closeSidenav = new EventEmitter<void>();
  @Input() set collapsed(value: boolean) {
    this.sidenavCollapsed.set(value);
  }
  sidenavCollapsed = signal<boolean>(false);
  profilePicSize = computed(() => this.sidenavCollapsed() ? '40' : '124');
  currentUser$!: Observable<User | null>;

  menuItems = signal<MenuItem[]>([
    {
      label: 'Menu 1',
      icon: 'receipt',
      route: '/app/menu1',
    },
    {
      label: 'Menu 2',
      icon: 'cloud_upload',
      route: '/app/menu2',
    },
    {
      label: 'Menu 3',
      icon: 'add',
      route: '/app/menu3',
    },
    {
      label: 'Menu 4',
      icon: 'receipt_long',
      route: '/app/menu4',
    }
  ]);

  constructor(private store: Store<AppState>) { }

  ngOnInit(): void {
    this.initializeValues();
  }

  onClickMenuItem() {
    this.closeSidenav.emit();
  }

  initializeValues() {
    this.currentUser$ = this.store.select(selectCurrentUser);
  }

  doLogout() {
    this.store.dispatch(logoutAction());
  }
}
