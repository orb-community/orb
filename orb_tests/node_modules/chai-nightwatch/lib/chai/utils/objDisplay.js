const inspect = require('./inspect');
const config = require('../config');

/**
 * ### .objDisplay (object)
 *
 * Determines if an object or an array matches
 * criteria to be inspected in-line for error
 * messages or should be truncated.
 *
 * @param {*} obj object to inspect
 * @name objDisplay
 * @api public
 */

module.exports = function (obj) {
  const str = inspect(obj);
  const type = Object.prototype.toString.call(obj);

  if (config.truncateThreshold && str.length >= config.truncateThreshold) {
    if (type === '[object Function]') {
      return !obj.name || obj.name === ''
        ? '[Function]'
        : '[Function: ' + obj.name + ']';
    }

    if (type === '[object Array]') {
      return '[ Array(' + obj.length + ') ]';
    }

    if (type === '[object Object]') {
      const keys = Object.keys(obj);
      const kstr = keys.length > 2
        ? keys.splice(0, 2).join(', ') + ', ...'
        : keys.join(', ');

      return '{ Object (' + kstr + ') }';
    }
  }

  return str;
};
