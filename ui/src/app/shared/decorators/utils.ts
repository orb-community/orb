import { debounce } from 'lodash';

/**
 * Debounce a method
 */
export function Debounce(ms) {
  return function (target: any, key: any, descriptor: any) {
    const oldFunc = descriptor.value;
    const newFunc = debounce(oldFunc, ms);
    descriptor.value = function () {
      return newFunc.apply(this, arguments);
    };
  };
}
