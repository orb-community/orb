import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

/**
 * Validates Tag object to have at least one key
 * tags !== {};
 */
export function tagInputValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors|null => {
    const { value } = control;

    if (!value) {
      return null;
    }

    // check if object tags has at least one key defined
    // key input validation guarantee no keys can be empty strings
    const hasAtLeastOneKey = Object.entries(value).length !== 0;

    return !hasAtLeastOneKey ? {tagMustDefineAtLeastOneKey: true} : null;
  };
}
