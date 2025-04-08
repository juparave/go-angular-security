// Angular Core Modules
import { NgModule } from '@angular/core'; // Decorator for defining Angular modules.
import { BrowserModule } from '@angular/platform-browser'; // Required for applications running in a browser. Provides core directives like NgIf, NgFor.
import { BrowserAnimationsModule } from '@angular/platform-browser/animations'; // Required for using Angular's animation features.

// Application Specific Modules and Components
import { AppRoutingModule } from './app-routing.module'; // Handles application routing configuration.
import { AppComponent } from './app.component'; // The root component of the application.

// NgRx State Management Modules
import { StoreModule } from '@ngrx/store'; // Provides the state management container (the Store).
import { EffectsModule } from '@ngrx/effects'; // Provides support for handling side effects (e.g., async operations like HTTP requests) triggered by actions.
import { StoreDevtoolsModule } from '@ngrx/store-devtools'; // Provides developer tools integration for inspecting state and actions.

// NgRx Effects for Authentication Flow
import { LoginEffect } from './store/effects/auth-login.effects'; // Handles side effects related to the login process.
import { RegisterEffect } from './store/effects/auth-register.effects'; // Handles side effects related to the user registration process.
import { LogoutEffect } from './store/effects/auth-logout.effects'; // Handles side effects related to the logout process.
import { GetCurrentUserEffect } from './store/effects/auth-get-current-user.effects'; // Handles fetching the currently logged-in user's data.
import { RefreshTokensEffect } from './store/effects/auth-refresh-tokens.effects'; // Handles refreshing authentication tokens.

// HTTP Client and Interceptors
import { HTTP_INTERCEPTORS, provideHttpClient, withInterceptorsFromDi, withJsonpSupport } from '@angular/common/http'; // Core HTTP client functionalities and the token for registering interceptors.
import { RefreshTokenInterceptor } from './services/auth/refresh-token.interceptor'; // Custom HTTP interceptor to automatically refresh JWT tokens when they expire.
import { LoadingInterceptor } from './services/loading.interceptors'; // Custom HTTP interceptor to show/hide a global loading indicator during HTTP requests.

// NgRx Reducers and Environment Configuration
import { authReducer } from './store/reducers/auth.reducers'; // Reducer function responsible for managing the 'auth' slice of the application state.
import { environment } from 'src/environments/environment'; // Environment-specific configuration (e.g., API URLs, production flags).

// Array defining application-wide providers, specifically HTTP interceptors.
// These interceptors will process every outgoing HTTP request.
const APP_PROVIDERS = [
  {
    provide: HTTP_INTERCEPTORS, // Token indicating this provider is for HTTP interceptors.
    useClass: RefreshTokenInterceptor, // The class implementing the refresh token logic.
    multi: true, // Specifies that this provider adds to a collection (multiple interceptors can exist).
  },
  {
    provide: HTTP_INTERCEPTORS, // Token indicating this provider is for HTTP interceptors.
    useClass: LoadingInterceptor, // The class implementing the loading indicator logic.
    multi: true, // Specifies that this provider adds to a collection.
  },
];

// Decorator that marks AppModule as an Angular module class (also called an NgModule class).
@NgModule({
  // Components, directives, and pipes that belong to this module.
  declarations: [
    AppComponent // The root component is declared here.
  ],
  // The main application view, called the root component, which hosts all other app views.
  // Only the root module should set the bootstrap property.
  bootstrap: [
    AppComponent // Specifies AppComponent as the component to load when the application starts.
  ],
  // Other modules whose exported classes are needed by component templates declared in this NgModule.
  imports: [
    BrowserModule, // Imports functionalities required for browser applications.
    BrowserAnimationsModule, // Imports animation functionalities.
    AppRoutingModule, // Imports the application's routing configuration.
    // Configures the NgRx store at the root level.
    // The `auth` key maps to the state slice managed by `authReducer`.
    StoreModule.forRoot({
      auth: authReducer,
    }),
    // Registers NgRx effects classes at the root level.
    // These effects listen for actions and perform side effects.
    EffectsModule.forRoot([
      RegisterEffect,
      LoginEffect,
      LogoutEffect,
      GetCurrentUserEffect,
      RefreshTokensEffect,
    ]),
    // Configures the NgRx Store Devtools extension.
    StoreDevtoolsModule.instrument({
      maxAge: 25, // Retains the last 25 states in the DevTools history.
      logOnly: environment.production, // Restricts the extension to log-only mode in production environments.
      connectInZone: true // Ensures DevTools operations run inside Angular's zone for better integration.
    })],
  // Creators of services that this NgModule contributes to the global collection of services;
  // they become accessible in all parts of the app. (Also used by Dependency Injection).
  providers: [
    APP_PROVIDERS, // Provides the custom HTTP interceptors defined above application-wide.
    // Provides and configures the modern HttpClient.
    // `withInterceptorsFromDi()` enables interceptors provided via DI (like ours).
    // `withJsonpSupport()` adds support for JSONP requests if needed.
    provideHttpClient(withInterceptorsFromDi(), withJsonpSupport())
  ],
})
export class AppModule { }
