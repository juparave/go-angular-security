import { Action, createAction, props } from '@ngrx/store';
import { Article } from 'src/app/model/article';

export enum ArticleActionType {
  ADD_ITEM = '[ARTICLE] Add Article',
}

export class AddItemAction implements Action {
  readonly type = ArticleActionType.ADD_ITEM;
  //add an optional payload
  constructor(public payload: Article) {}
}

export type ArticleAction = AddItemAction;

export const addItemAction = createAction(
  ArticleActionType.ADD_ITEM,
  props<{ payload: Article }>()
);
