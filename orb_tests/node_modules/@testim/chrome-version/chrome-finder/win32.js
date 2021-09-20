const path = require('path');
const { canAccess } = require('./util');

function win32() {
  const installations = [];
  const suffixes = [
    '\\Google\\Chrome SxS\\Application\\chrome.exe',
    '\\Google\\Chrome\\Application\\chrome.exe',
    '\\chrome-win32\\chrome.exe',
    '\\Chromium\\Application\\chrome.exe',
    // '\\Google\\Chrome Beta\\Application\\chrome.exe',
  ];
  const prefixes =
    [process.env.LOCALAPPDATA, process.env.PROGRAMFILES, process.env['PROGRAMFILES(X86)']];

  prefixes.forEach(prefix => suffixes.forEach(suffix => {
    if (prefix) {
      const chromePath = path.join(prefix, suffix);
      if (canAccess(chromePath)) {
        installations.push(chromePath);
      }
    }
  }));
  return installations;
}

module.exports = win32;
