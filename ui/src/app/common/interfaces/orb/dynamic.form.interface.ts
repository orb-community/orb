/**
 * @Interface DynamicInput
 * Dynamic form input model
 */
export interface DynamicInput {
  // label to be shown with descriptive name
  label?: string;

  // type of value - string, number, boolean
  type?: string;

  // type of input component to use
  input?: string;

  // longer description for tooltip or help text
  description?: string;

  // dynamic properties list for different types of inputs and configs
  props?: {
    // is this a required field?
    // default|missing -> false
    required?: boolean;

    // short example string for input value - placeholder
    example?: string;

    // leave space for any extra
    [propName: string]: string[]|number[]|string|number|any; // {[propName: string]: [] | {} | any}
  };

  // optional name used when parsing results from service for ease of use
  // name for field
  name?: string;
}

/**
 * @Interface DynamicForm
 * Dynamic form config model
 */
export interface DynamicFormConfig {
  [propName: string]: DynamicInput;
}

