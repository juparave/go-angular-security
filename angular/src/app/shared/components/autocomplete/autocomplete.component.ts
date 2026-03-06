import { CommonModule } from '@angular/common';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule, MatAutocompleteTrigger } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { Observable, map, startWith } from 'rxjs';

type MyOptions = {
  label: string;
  value: string;
};

/**
 * AutocompleteComponent is a custom component used for displaying an autocomplete input field.
 * It provides suggestions based on the user's input and allows selection from a list of options.
 * 
 * ref: https://stackblitz.com/edit/angular-custom-mat-auto-complete?file=src%2Fapp%2Fautocomplete-display-example.ts,src%2Fapp%2Fautocomplete-display-example.html
 * ref: https://material.angular.io/components/autocomplete/overview
 */
@Component({
  selector: 'app-autocomplete',
  standalone: true,
  imports: [
    CommonModule,
    MatAutocompleteModule,
    MatInputModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatButtonModule,
    MatIconModule,
  ],
  templateUrl: './autocomplete.component.html',
  styleUrl: './autocomplete.component.scss'
})
export class AutocompleteComponent implements OnInit {
  @ViewChild('inputAutoComplete') inputAutoComplete: any;

  @Input() myControl!: FormControl;
  @Input() options$!: Observable<MyOptions[]>;
  @Input() allOptions!: MyOptions[];
  @Input() label!: string;
  @Input() placeholder!: string;

  constructor() { }

  ngOnInit() {
    this.options$ = this.myControl.valueChanges
      .pipe(
        startWith(''),
        map(value => this._filter(value || ''))
      );
  }

  /**
   * Filters the options based on the provided value.
   * 
   * @param value - The value to filter the options by.
   * @returns An array of filtered options.
   */
  private _filter(value: string): MyOptions[] {
    const filterValue = value.toLowerCase();
    return this.allOptions.filter(option => option.label.toLowerCase().includes(filterValue));
  }

  /**
   * Returns the label of the option that matches the given value.
   * If no match is found, an empty string is returned.
   * 
   * @param value - The value to search for.
   * @returns The label of the matching option, or an empty string if no match is found.
   */
  displayFn(value: string): string {
    return this.allOptions?.find(option => option.value === value)?.label || '';
  }


  /**
   * Clears the input field and focuses on it.
   * @param evt - The event object.
   */
  clearInput(evt: any): void {
    evt.stopPropagation();
    this.myControl?.reset();
    this.inputAutoComplete?.nativeElement.focus();
  }

  /**
   * Opens or closes the panel of the autocomplete component based on the current state of the trigger.
   * @param evt - The event object.
   * @param trigger - The MatAutocompleteTrigger instance.
   */
  openOrClosePanel(evt: any, trigger: MatAutocompleteTrigger): void {
    evt.stopPropagation();
    if (trigger.panelOpen)
      trigger.closePanel();
    else
      trigger.openPanel();
  }

}
