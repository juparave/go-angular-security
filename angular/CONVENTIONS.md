# Angular Project Conventions

This document outlines the key conventions and patterns used in our Angular project, with a focus on state management, route guarding, subscription management, and overall project organization.

## Table of Contents

1. [Project Structure](#project-structure)
2. [NgRx State Management](#ngrx-state-management)
3. [Route Guards](#route-guards)
4. [Subscription Management](#subscription-management)
5. [API Calls](#api-calls)
6. [Shared Module Organization](#shared-module-organization)
7. [Feature Module Organization](#feature-module-organization)
8. [Component Organization](#component-organization)
9. [Naming Conventions](#naming-conventions)
10. [Testing](#testing)
11. [Forms and Validation](#forms-and-validation)
12. [Styling](#styling)

## Project Structure

Our project follows a modular architecture with a clear separation of concerns:

```
src/
├── app/
│   ├── models/                 # TypeScript interfaces
│   ├── modules/                # Feature modules
│   │   ├── app/                # Main app feature
│   │   │   ├── base/           # Base layout
│   │   │   ├── custom-sidenav/ # Navigation
│   │   │   └── dashboard/      # Dashboard
│   │   └── login/              # Authentication features
│   │       ├── form/           # Login form
│   │       ├── request-reset-password/
│   │       ├── reset-password/
│   │       └── sign-up/
│   ├── services/               # HTTP and utility services
│   │   ├── auth/               # Authentication services and guards
│   │   └── ...                 # Other services
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
├── assets/                     # Static assets
│   ├── images/                 # Image files
│   └── scss/                   # Global SCSS
│       ├── _mixins.scss        # SCSS mixins
│       ├── _variables.scss     # SCSS variables
│       └── styles.scss         # Main styles
└── environments/               # Environment configuration
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

### Guards Location and Organization

Route guards are located in the `services/auth` directory:

```
services/
└── auth/
    ├── auth.guard.ts
    ├── unauth.guard.ts
    └── subscription.guard.ts
```

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

### Service Organization

Services are organized in the `services` directory:

```
services/
├── auth/
│   ├── auth.service.ts
│   └── refresh-token.interceptor.ts
├── file-upload.service.ts
├── loading.interceptors.ts
├── loading.service.ts
├── notification.service.ts
├── persistance.service.ts
└── subscription.service.ts
```

### Service Pattern

- Create dedicated services for API calls
- Return typed Observables from service methods
- Use environment variables for API URLs

```typescript
@Injectable({
  providedIn: 'root'
})
export class SubscriptionService {
  private baseUrl = `${environment.apiUrl}/subscriptions`;

  constructor(private http: HttpClient) {}

  getSubscription(): Observable<Subscription> {
    return this.http.get<Subscription>(`${this.baseUrl}/current`);
  }
  
  // Additional methods for CRUD operations
}
```

### Environment Configuration

For consistent API URL handling, follow this standardized approach:

1. Define the base API URL in the environment files:
   ```typescript
   // environment.ts
   export const environment = {
     production: false,
     apiUrl: 'http://localhost:3000/api'
   };
   ```

2. In service classes, define a base URL property for the resource:
   ```typescript
   // Define a base URL property in the service class
   private baseUrl = `${environment.apiUrl}/resource`;
   ```

3. Reference this property for all endpoints within the service:
   ```typescript
   // Get all resources
   getAll(): Observable<Resource[]> {
     return this.http.get<Resource[]>(this.baseUrl);
   }
   
   // Get single resource by ID
   getById(id: string): Observable<Resource> {
     return this.http.get<Resource>(`${this.baseUrl}/${id}`);
   }
   ```

4. For services that handle multiple resource types, organize URLs by resource:
   ```typescript
   private urls = {
     users: `${environment.apiUrl}/users`,
     profiles: `${environment.apiUrl}/profiles`,
     settings: `${environment.apiUrl}/settings`
   };
   ```

### HTTP Interceptors

Use interceptors for cross-cutting concerns:

```typescript
@Injectable()
export class RefreshTokenInterceptor implements HttpInterceptor {
  constructor(private store: Store<AppState>) {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    // Implementation for handling token refresh
  }
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

## Feature Module Organization

Our application is organized into feature modules:

```
modules/
├── app/                # Main application feature
│   ├── base/           # Base layout component
│   ├── custom-sidenav/ # Navigation sidebar
│   └── dashboard/      # Dashboard component
└── login/              # Authentication feature
    ├── form/           # Login form
    ├── request-reset-password/
    ├── reset-password/
    └── sign-up/
```

### Feature Module Structure

Each feature module includes:
- Components for the feature UI
- A module file for dependency management
- Feature-specific routing if needed

```typescript
@NgModule({
  declarations: [
    LoginComponent,
    LoginFormComponent,
    RequestResetPasswordComponent,
    ResetPasswordComponent,
    SignUpComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule.forChild(routes),
    SharedModule,
    BackendErrorMessagesModule
  ]
})
export class LoginModule {}
```

## Component Organization

### Component Structure

Each component consists of:
- HTML template (*.component.html)
- SCSS styles (*.component.scss)
- TypeScript logic (*.component.ts)
- Test specifications (*.component.spec.ts)

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

- Observables end with `$` (e.g., `user$`)
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

### Form Structure

Forms should follow these patterns:
- Use Reactive Forms over Template-driven forms
- Group related form controls with FormGroup
- Apply validators at both field and form level
- Use custom validators for complex validation logic

```typescript
// Example login form
this.loginForm = this.formBuilder.group({
  email: ['', [Validators.required, Validators.email]],
  password: ['', [Validators.required, Validators.minLength(8)]]
});
```

### Error Handling

Use the BackendErrorMessages component for displaying server-side errors:

```html
<app-backend-error-messages 
  [backendErrors]="validationErrors$ | async">
</app-backend-error-messages>
```

## Styling

### SCSS Organization

Styles are organized in a dedicated SCSS directory:

```
assets/
└── scss/
    ├── _mixins.scss    # Reusable SCSS mixins
    ├── _variables.scss # Global SCSS variables
    └── styles.scss     # Main styles entry point
```

### Component Styling

- Each component has its own SCSS file
- Component styles should be scoped to the component
- Use BEM naming for CSS classes when appropriate
- Leverage Angular Material's theming system

```scss
// Example component SCSS structure
.container {
  width: 100%;
}

.mat-mdc-card {
  width: 380px;
  position: absolute;
  top: 20%;
  left: 50%;
  transform: translate(-50%, -50%);
}
```

### Theme Variables

Use the variables defined in _variables.scss for consistent styling:

```scss
// Example of using global variables
.alert {
  &.alert-success {
    background: $success-color;
    color: $success-contrast;
  }
}
```

## Conclusion

Following these conventions ensures consistency across our codebase and makes it easier for team members to understand and maintain the project. If you have questions or suggestions for improvements, please discuss with the team.