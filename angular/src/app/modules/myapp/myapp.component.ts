import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';

@Component({
  selector: 'app-app',
  templateUrl: './myapp.component.html',
  styleUrls: ['./myapp.component.scss'],
})
export class MyAppComponent implements OnInit {
  constructor(private store: Store<AppState>) {}

  ngOnInit(): void {}
}
