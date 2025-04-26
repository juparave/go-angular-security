# Angular Project Conventions

This document outlines the key conventions and patterns used in our Angular project, with a focus on state management, route guarding, subscription management, and shared components organization.

## Table of Contents

1. [Project Structure](#project-structure)
2. [NgRx State Management](#ngrx-state-management)
3. [Route Guards](#route-guards)
4. [Subscription Management](#subscription-management)
5. [API Calls](#api-calls)
6. [Shared Module Organization](#shared-module-organization)
7. [Component Organization](#component-organization)
8. [Naming Conventions](#naming-conventions)
9. [Testing](#testing)
10. [Forms and Validation](#forms-and-validation)

## Project Structure

```
src/
├── app/
│   ├── components/             # Shared components
│   ├── features/               # Feature modules
│   │   ├── auth/               # Authentication feature
│   │   ├── basic/              # Basic subscription features
│   │   ├── pro/                # Pro subscription features
│   │   └── trial/              # Trial subscription features
│   ├── guards/                 # Route guards
│   ├── models/                 # TypeScript interfaces
│   ├── services/               # HTTP and utility services
│   ├── shared/                 # Shared modules and components
│   │   ├── components/         # Standalone shared components
│   │   ├── modules/            # Shared module bundles
│   │   ├── pipes/              # Custom pipes
│   │   ├── util/               # Utility functions
│   │   └── validators/         # Custom form validators
│   └── store/                  # NgRx state management
│       ├── actions/            # Action creators
│       ├── effects/            # Side effects
│       ├── interfaces/         # State interfaces
│       ├── reducers/           # State reducers
│       ├── selectors/          # Memoized selectors
│       └── types/              # Type definitions
```

## NgRx State Management

### State Organization

- Each feature of the application has its own dedicated state slice
- State is defined in interfaces that extend the app state
- Example state slices include: `auth`, `subscription`, etc.

### State Interface Pattern

```typescript
// Feature state interface (e.g., src/app/store/interfaces/subscription-state.ts)
export interface SubscriptionState {
  subscription: Subscription | null;
  isLoading: boolean;
  isUpdating: boolean;
  validationErrors: BackendErrors | null;
  redirectUrl: string | null;
}

// Root app state (src/app/store/interfaces/app-state.ts)
export interface AppState {
  readonly auth: AuthState;
  readonly subscription: SubscriptionState;
  // Additional state slices...
}
```

### Actions

- Actions are grouped by feature
- Use enum for action types to avoid string duplication
- Follow the pattern of `[Feature] Action description`
- Create separate actions for request, success, and failure states

```typescript
export enum ActionTypes {
  GET_SUBSCRIPTION = '[Subscription] Get subscription',
  GET_SUBSCRIPTION_SUCCESS = '[Subscription] Get subscription success',
  GET_SUBSCRIPTION_FAILURE = '[Subscription] Get subscription failure',
}

export const getSubscriptionAction = createAction(ActionTypes.GET_SUBSCRIPTION);
export const getSubscriptionSuccessAction = createAction(
  ActionTypes.GET_SUBSCRIPTION_SUCCESS,
  props<{ subscription: Subscription }>()
);
export const getSubscriptionFailureAction = createAction(
  ActionTypes.GET_SUBSCRIPTION_FAILURE,
  props<{ errors: BackendErrors }>()
);
```

### Reducers

- Use the NgRx `createReducer` and `on` functions
- Handle loading states in the reducer
- Clear validation errors when starting a new request
- Return a new state object, don't mutate the existing state

```typescript
export const subscriptionReducer = createReducer(
  initialState,
  on(
    getSubscriptionAction,
    (state): SubscriptionState => ({
      ...state,
      isLoading: true,
      validationErrors: null,
    })
  ),
  on(
    getSubscriptionSuccessAction,
    (state, action): SubscriptionState => ({
      ...state,
      isLoading: false,
      subscription: action.subscription,
    })
  ),
  on(
    getSubscriptionFailureAction,
    (state, action): SubscriptionState => ({
      ...state,
      isLoading: false,
      validationErrors: action.errors,
    })
  )
);
```

### Selectors

- Use `createSelector` for memoization
- Create a base selector for each feature slice
- Derive more specific selectors from the base selector
- Include loading states in the selector results

```typescript
// Base selector
export const selectSubscriptionState = (state: AppState): SubscriptionState => state.subscription;

// Derived selectors
export const selectSubscription = createSelector(
  selectSubscriptionState,
  (state: SubscriptionState): { subscription: Subscription | null; isLoading: boolean } => ({
    subscription: state.subscription,
    isLoading: state.isLoading
  })
);

export const selectHasActivePro = createSelector(
  selectSubscription,
  (state): { hasActivePro: boolean; isLoading: boolean } => {
    const subscription = state.subscription;
    
    const hasActivePro = !!subscription && 
      subscription.plan === 'pro' && 
      isActiveSubscription(subscription);
    
    return {
      hasActivePro,
      isLoading: state.isLoading
    };
  }
);
```

### Effects

- Organize effects by feature and action type
- Use a class-per-file approach for effects
- Inject services into effects for API calls
- Use `switchMap` for canceling previous requests

```typescript
@Injectable()
export class GetSubscriptionEffect {
  constructor(
    private actions$: Actions,
    private subscriptionService: SubscriptionService
  ) {}

  getSubscription$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(getSubscriptionAction),
      switchMap(() => {
        return this.subscriptionService.getSubscription().pipe(
          map((subscription: Subscription) => {
            return getSubscriptionSuccessAction({ subscription });
          }),
          catchError((errorResponse: HttpErrorResponse) => {
            return of(getSubscriptionFailureAction({
              errors: { subscription: [errorResponse.error.message] },
            }));
          })
        );
      })
    );
  });
}
```

## Route Guards

### Class + Function Approach

Use a class-based guard with a functional wrapper:

```typescript
@Injectable({
  providedIn: 'root',
})
export class SubscriptionGuard {
  constructor(private store: Store<AppState>, public router: Router) {}

  createGuard(selector: any, redirectTo: string = '/upgrade') {
    return (route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree> => {
      // Implementation using the selector to determine access
    };
  }
}

export const canActivatePro: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(SubscriptionGuard).createGuard(selectHasActivePro, '/upgrade')(route, state);
};
```

### Guard Naming

- Use `canActivate[Feature]` naming for functional guards (e.g., `canActivatePro`)
- Use `[Feature]Guard` naming for guard classes (e.g., `SubscriptionGuard`)

## Subscription Management

### Subscription Status Checking

Use dedicated selectors for different subscription types and states:

```typescript
export const selectHasActivePro = createSelector(
  selectSubscription,
  (state): { hasActivePro: boolean; isLoading: boolean } => {
    const subscription = state.subscription;
    
    const hasActivePro = !!subscription && 
      subscription.plan === 'pro' && 
      (subscription.status === 'active' || 
       (subscription.status === 'canceled' && 
         subscription.currentPeriodEnd && 
         new Date() < new Date(subscription.currentPeriodEnd)));
    
    return {
      hasActivePro,
      isLoading: state.isLoading
    };
  }
);
```

### Subscription Route Protection

Apply subscription guards to protect premium features:

```typescript
const routes: Routes = [
  { 
    path: 'features/pro', 
    loadChildren: () => import('./features/pro/pro.module').then(m => m.ProModule),
    canActivate: [canActivatePro]
  }
];
```

## API Calls

### Service Pattern

- Create dedicated services for API calls
- Return typed Observables from service methods
- Use environment variables for API URLs

```typescript
@Injectable({
  providedIn: 'root'
})
export class SubscriptionService {
  constructor(private http: HttpClient) {}

  getSubscription(): Observable<Subscription> {
    return this.http.get<Subscription>(`${environment.apiUrl}/subscriptions/current`);
  }
  
  // Additional methods for CRUD operations
}
```

## Shared Module Organization

### Module Structure

Our shared modules are organized as follows:

```
shared/
├── components/         # Standalone shared components
│   ├── autocomplete/   # Example component
│   └── page-not-found/ # Example component
├── modules/            # Shared module bundles
│   ├── alert/          # Alert module with component
│   └── backend-error-messages/ # Error handling module
├── pipes/              # Custom pipes
├── util/               # Utility functions
└── validators/         # Custom form validators
```

### SharedModule

The main `SharedModule` exports commonly used Angular Material modules:

```typescript
@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    MatToolbarModule,
    MatInputModule,
    // ...other Angular Material modules
  ],
  exports: [
    MatToolbarModule,
    MatInputModule,
    // ...same Angular Material modules
  ],
})
export class SharedModule { }
```

### Feature-specific Shared Modules

Feature-specific shared modules should:
- Group related components, directives, and pipes
- Export all declarations that need to be used outside the module
- Import only what they need

```typescript
@NgModule({
  declarations: [AlertComponent],
  imports: [CommonModule],
  exports: [AlertComponent],
})
export class AlertModule {}
```

## Component Organization

### Standalone vs Module-based Components

- Use standalone components for simple, self-contained UI elements
- Use module-based components when they have complex dependencies or need to be grouped

### Standalone Components

```typescript
@Component({
  selector: 'app-page-not-found',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, RouterModule],
  templateUrl: './page-not-found.component.html',
  styleUrl: './page-not-found.component.scss'
})
export class PageNotFoundComponent { }
```

### Component Modules

```typescript
@NgModule({
  declarations: [BackendErrorMessagesComponent],
  imports: [CommonModule, AlertModule],
  exports: [BackendErrorMessagesComponent],
})
export class BackendErrorMessagesModule {}
```

## Naming Conventions

### Files

- Feature modules: `feature-name.module.ts`
- Components: `component-name.component.ts`
- Services: `service-name.service.ts`
- Guards: `guard-name.guard.ts`
- Models/interfaces: `model-name.ts`
- State: `feature-name-state.ts`
- Actions: `feature-name.actions.ts`
- Reducers: `feature-name.reducers.ts`
- Selectors: `feature-name.selectors.ts`
- Effects: `feature-name-action-name.effects.ts`
- Pipes: `pipe-name.pipe.ts`
- Validators: `validator-name.validator.ts`

### Variables and Methods

- Observables end with ` (e.g., `user)
- Boolean variables start with is/has/should (e.g., `isLoading`, `hasAccess`)
- Actions follow the pattern:
  - `[get/update/delete/etc][Entity]Action`
  - `[get/update/delete/etc][Entity]SuccessAction`
  - `[get/update/delete/etc][Entity]FailureAction`

## Testing

### Component Testing

Each component should have a spec file with basic tests:

```typescript
describe('PageNotFoundComponent', () => {
  let component: PageNotFoundComponent;
  let fixture: ComponentFixture<PageNotFoundComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PageNotFoundComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PageNotFoundComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
```

## Forms and Validation

### Custom Validators

Create reusable validators in the shared/validators directory:

```typescript
// Example of a password confirmation validator
export function MustMatch(
  controlName: string,
  matchingControlName: string
): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const originalControl = control.get(controlName);
    const matchingControl = control.get(matchingControlName);

    if (matchingControl!.errors && !matchingControl!.errors['mustMatch']) {
      return null;
    }

    if (originalControl!.value !== matchingControl!.value) {
      matchingControl!.setErrors({ mustMatch: true });
      return { mustMatch: true };
    } else {
      matchingControl!.setErrors(null);
    }
    return null;
  };
}
```

### Custom Pipes

Create reusable pipes in the shared/pipes directory:

```typescript
@Pipe({
  name: 'role',
  standalone: true
})
export class RolePipe implements PipeTransform {
  transform(roles: string | string[] | undefined, roleName: string): boolean {
    if (!roles) {
      return false;
    }
    if (typeof roles === 'string') {
      roles = roles.split(',');
    }
    return roles.includes(roleName);
  }
}
```

## Conclusion

Following these conventions ensures consistency across our codebase and makes it easier for team members to understand and maintain the project. If you have questions or suggestions for improvements, please discuss with the team.