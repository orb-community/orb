import { Directive } from '@angular/core';
import { AbstractControl, NG_VALIDATORS, ValidationErrors, Validator } from '@angular/forms';
import { tagInputValidator } from 'app/shared/directives/tag-input.validator';

@Directive({
  selector: '[ngxOrbTagInput]',
  providers: [{
    provide: NG_VALIDATORS,
    useExisting: tagInputValidator,
    multi: true,
  }],
})
export class ValidTagInputDirective implements Validator {

  validate(control: AbstractControl): ValidationErrors|null {
    return tagInputValidator()(control);
  }
}
