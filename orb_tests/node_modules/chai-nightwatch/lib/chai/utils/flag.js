/*!
 * Chai - flag utility
 * Copyright(c) 2012-2014 Jake Luer <jake@alogicalparadox.com>
 * MIT Licensed
 */

/**
 * ### flag(object, key, [value])
 *
 * Get or set a flag value on an object. If a
 * value is provided it will be set, else it will
 * return the currently set value or `undefined` if
 * the value is not set.
 *
 *     utils.flag(this, 'foo', 'bar'); // setter
 *     utils.flag(this, 'foo'); // getter, returns `bar`
 *
 * @param {Object} obj constructed Assertion
 * @param {String} key
 * @param {*} value (optional)
 * @name flag
 * @api private
 */

module.exports = function (obj, key, value) {
  const flags = obj.__flags || (obj.__flags = new Map());

  if (arguments.length === 1) {
    return flags;
  }

  if (arguments.length === 2) {
    return flags.get(key);
  }

  flags.set(key, value);
};
