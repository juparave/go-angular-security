import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { StoreModule } from '@ngrx/store';
import { EffectsModule } from '@ngrx/effects';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { articleReducers } from 'src/app/store/reducers/article.reducer';
import { HeaderComponent } from './header/header.component';
import { LoginEffect } from './store/effects/auth-login.effects';
import { RegisterEffect } from './store/effects/auth-register.effects';
import { HttpClientJsonpModule, HttpClientModule } from '@angular/common/http';
import { LogoutEffect } from './store/effects/auth-logout.effects';
import { GetCurrentUserEffect } from './store/effects/auth-get-current-user.effects';
import { RefreshTokensEffect } from './store/effects/auth-refresh-tokens.effects';
import { authReducer } from './store/reducers/auth.reducers';
import { environment } from 'src/environments/environment';

@NgModule({
  declarations: [AppComponent, HeaderComponent],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,

    HttpClientModule,
    HttpClientJsonpModule,
    AppRoutingModule,

    // StoreModule.forRoot(
    //   {
    //     articles: articleReducers,
    //   },
    //   {}
    // ),
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
    }),
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
