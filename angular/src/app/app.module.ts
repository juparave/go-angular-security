import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { LoginEffect } from './store/effects/auth-login.effects';
import { RegisterEffect } from './store/effects/auth-register.effects';
import {
  HttpClientJsonpModule,
  HttpClientModule,
  HTTP_INTERCEPTORS,
} from '@angular/common/http';
import { LogoutEffect } from './store/effects/auth-logout.effects';
import { GetCurrentUserEffect } from './store/effects/auth-get-current-user.effects';
import { RefreshTokensEffect } from './store/effects/auth-refresh-tokens.effects';
import { authReducer } from './store/reducers/auth.reducers';
import { environment } from 'src/environments/environment';
import { RefreshTokenInterceptor } from './services/auth/refresh-token.interceptor';
import { LoadingInterceptor } from './services/loading.interceptors';

const APP_PROVIDERS = [
  {
    provide: HTTP_INTERCEPTORS,
    useClass: RefreshTokenInterceptor,
    multi: true,
  },
  {
    provide: HTTP_INTERCEPTORS,
    useClass: LoadingInterceptor,
    multi: true,
  },
];

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,

    HttpClientModule,
    HttpClientJsonpModule,
    AppRoutingModule,

    StoreModule.forRoot({
      auth: authReducer,
    }),
    EffectsModule.forRoot([
      RegisterEffect,
      LoginEffect,
      LogoutEffect,
      GetCurrentUserEffect,
      RefreshTokensEffect,
    ]),
    StoreDevtoolsModule.instrument({
      maxAge: 25,
      logOnly: environment.production,
      connectInZone: true
    }),
  ],
  providers: [APP_PROVIDERS],
  bootstrap: [AppComponent]
})
export class AppModule { }
