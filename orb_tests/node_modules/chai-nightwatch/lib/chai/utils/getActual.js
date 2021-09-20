/*!
 * Chai - getActual utility
 * Copyright(c) 2012-2014 Jake Luer <jake@alogicalparadox.com>
 * MIT Licensed
 */

/**
 * # getActual(object, [actual])
 *
 * Returns the `actual` value for an Assertion
 *
 * @param {Object} obj (constructed Assertion)
 * @param {Arguments} chai.Assertion.prototype.assert args
 */

module.exports = function (obj, args) {
  return args.length > 4 ? args[4] : obj._obj;
};
