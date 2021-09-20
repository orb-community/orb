/*!
 * Based on chai library
 *
 * Copyright(c) 2011-2014 Jake Luer <jake@alogicalparadox.com>
 * MIT Licensed
 */

const used = [];
/*!
 * Chai version
 */

exports.version = '2.2.0';

/*!
 * Assertion Error
 */

exports.AssertionError = require('assertion-error');

/*!
 * Utils for plugins (not exported)
 */

const util = require('./chai/utils');

/**
 * # .use(function)
 *
 * Provides a way to extend the internals of Chai
 *
 * @param {Function} fn
 * @returns {this} for chaining
 * @api public
 */

exports.use = function (fn) {
  if (!~used.indexOf(fn)) {
    fn(this, util);
    used.push(fn);
  }

  return this;
};

/*!
 * Utility Functions
 */

exports.util = util;

/*!
 * Primary `Assertion` prototype
 */

const assertion = require('./chai/assertion');
exports.use(assertion);

/*!
 * Core Assertions
 */

const core = require('./chai/core/assertions');
exports.use(core);

/*!
 * Expect interface
 */

const expect = require('./chai/interface/expect');
exports.use(expect);

exports.flag = require('./chai/utils/flag');
