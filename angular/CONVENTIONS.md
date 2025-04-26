# Angular Project Conventions

This document outlines the key conventions and patterns used in our Angular
project, with a focus on state management, route guarding, and feature
organization.

## Table of Contents

1. [Project Structure](#project-structure)
2. [NgRx State Management](#ngrx-state-management)
3. [Route Guards](#route-guards)
4. [Subscription Management](#subscription-management)
5. [API Calls](#api-calls)
6. [Naming Conventions](#naming-conventions)
7. [Module Organization](#module-organization)

## Project Structure

```
src/
├── app/
│   ├── components/             # Shared components
│   ├── features/               # Feature modules
│   ├── guards/                 # Route guards
│   ├── models/                 # TypeScript interfaces
│   ├── services/               # HTTP and utility services
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
  someData: DataType | null;
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
  GET_ENTITY = '[Feature] Get entity',
  GET_ENTITY_SUCCESS = '[Feature] Get entity success',
  GET_ENTITY_FAILURE = '[Feature] Get entity failure',
}

export const getEntityAction = createAction(ActionTypes.GET_ENTITY);
export const getEntitySuccessAction = createAction(
  ActionTypes.GET_ENTITY_SUCCESS,
  props<{ entity: Entity }>()
);
export const getEntityFailureAction = createAction(
  ActionTypes.GET_ENTITY_FAILURE,
  props<{ errors: BackendErrors }>()
);
```

### Reducers

- Use the NgRx `createReducer` and `on` functions
- Handle loading states in the reducer
- Clear validation errors when starting a new request
- Return a new state object, don't mutate the existing state

```typescript
export const featureReducer = createReducer(
  initialState,
  on(
    getEntityAction,
    (state): FeatureState => ({
      ...state,
      isLoading: true,
      validationErrors: null,
    })
  ),
  on(
    getEntitySuccessAction,
    (state, action): FeatureState => ({
      ...state,
      isLoading: false,
      entity: action.entity,
    })
  ),
  on(
    getEntityFailureAction,
    (state, action): FeatureState => ({
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
export const selectFeatureState = (state: AppState): FeatureState => state.feature;

// Derived selectors
export const selectEntity = createSelector(
  selectFeatureState,
  (state: FeatureState): { entity: Entity | null; isLoading: boolean } => ({
    entity: state.entity,
    isLoading: state.isLoading
  })
);
```

### Effects

- Organize effects by feature and action type
- Use a class-per-file approach for effects
- Inject services into effects for API calls
- Use `switchMap` for canceling previous requests

```typescript
@Injectable()
export class GetEntityEffect {
  constructor(
    private actions$: Actions,
    private entityService: EntityService
  ) {}

  getEntity$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(getEntityAction),
      switchMap(() => {
        return this.entityService.getEntity().pipe(
          map((entity: Entity) => {
            return getEntitySuccessAction({ entity });
          }),
          catchError((errorResponse: HttpErrorResponse) => {
            return of(getEntityFailureAction({
              errors: { entity: [errorResponse.error.message] },
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
export class FeatureGuard {
  constructor(private store: Store<AppState>, public router: Router) {}

  createGuard(selector: any, redirectTo: string = '/fallback') {
    return (route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree> => {
      // Implementation...
    };
  }
}

export const canActivateFeature: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(FeatureGuard).createGuard(selectFeatureAccess)(route, state);
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
      isActiveSubscription(subscription);
    
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
export class EntityService {
  constructor(private http: HttpClient) {}

  getEntity(): Observable<Entity> {
    return this.http.get<Entity>(`${environment.apiUrl}/entities/current`);
  }
}
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

### Variables and Methods

- Observables end with `$` (e.g., `user$`)
- Boolean variables start with is/has/should (e.g., `isLoading`, `hasAccess`)
- Actions follow the pattern:
  - `[get/update/delete/etc][Entity]Action`
  - `[get/update/delete/etc][Entity]SuccessAction`
  - `[get/update/delete/etc][Entity]FailureAction`

## Module Organization

### Lazy Loading

- Use lazy loading for feature modules
- Apply route guards at the feature module level
- Structure routes hierarchically

```typescript
const routes: Routes = [
  {
    path: 'app',
    children: [
      { 
        path: 'features/premium', 
        loadChildren: () => import('./features/premium/premium.module').then(m => m.PremiumModule),
        canActivate: [canActivatePremium]
      }
    ],
    canActivate: [canActivateAuth]
  }
];
```

### Feature Module Structure

Feature modules should:
- Have their own routing
- Declare components used only within the feature
- Import only what they need

```typescript
@NgModule({
  declarations: [
    FeatureComponent,
    FeatureDetailComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    // Other needed modules
  ]
})
export class FeatureModule {}
```

## Conclusion

Following these conventions ensures consistency across our codebase and makes it easier for team members to understand and maintain the project. If you have questions or suggestions for improvements, please discuss with the team.
