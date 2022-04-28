// import the interface
import { Action, createReducer, on } from '@ngrx/store';
import { Article } from 'src/app/model/article';
import { addItemAction } from 'src/app/store/actions/article.action';

// create a dummy initial state
const initialState: Array<Article> = [
  {
    id: 1,
    content: 'Computer Engineering',
    user_id: 1,
  },
];

const reducers = createReducer(
  initialState,
  on(
    addItemAction,
    (state, action): Array<Article> => [...state, action.payload]
  )
);

export function articleReducers(
  state: Array<Article> | undefined,
  action: Action
) {
  return reducers(state, action);
}
