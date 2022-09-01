import {AbstractControl, NG_VALIDATORS, ValidationErrors, Validator, ValidatorFn} from '@angular/forms';
import {Directive} from '@angular/core';

@Directive({
  selector: '[ngxEmptyInput]',
  providers: [{provide: NG_VALIDATORS,
    useExisting: EmptyInputDirective,
    multi: true}],
})
export class EmptyInputDirective implements Validator {
  validate(control: AbstractControl): ValidationErrors | null {
    return emptyInputValidator()(control);
  }
}

export function emptyInputValidator(): ValidatorFn {
  return (control): ValidationErrors | null => {
    const {value} = control;
    if (!value) return null;
    const trimmed = control.value.trim();
    return trimmed === '' ? {emptyInput: {value: control.value}} : null;
  };
}
