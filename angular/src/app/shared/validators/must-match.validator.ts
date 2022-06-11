import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

// custom validator to check that two fields match
export function MustMatch(
  controlName: string,
  matchingControlName: string
): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const originalControl = control.get(controlName);
    const matchingControl = control.get(matchingControlName);

    if (matchingControl!.errors && !matchingControl!.errors['mustMatch']) {
      // return if another validator has already found an error on the matchingControl
      return null;
    }

    // set error on matchingControl if validation fails
    if (originalControl!.value !== matchingControl!.value) {
      matchingControl!.setErrors({ mustMatch: true });
      return { mustMatch: true };
    } else {
      matchingControl!.setErrors(null);
    }
    return null;
  };
}
