/**
 * @Interface DynamicInput
 * Dynamic form input model
 */
export interface DynamicInput {
  // label to be shown with descriptive name
  label?: string;

  // field name for property -> payload
  name?: string;

  // type of value - string, number, boolean
  type?: string;

  // type of input component to use
  input?: string;

  // longer description for tooltip or help text
  description?: string;

  // is this a required field?
  required?: boolean;

  // dynamic properties list for different types of inputs and configs
  props?: any; // {[propName: string]: [] | {} | any}
}

/**
 * @Interface DynamicForm
 * Dynamic form config model
 */
export interface DynamicFormConfig {
  [propName: string]: DynamicInput;
}

