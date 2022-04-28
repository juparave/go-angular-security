import { Article } from './article';

export interface AppState {
  readonly articles: Array<Article>;
}
