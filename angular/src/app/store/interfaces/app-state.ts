import { Article } from 'src/app/model/article';
import { AuthState } from './auth-state';

export interface AppState {
  readonly articles: Array<Article>;
  readonly auth: AuthState;
}
