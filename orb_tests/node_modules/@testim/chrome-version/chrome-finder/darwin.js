const { execSync } = require('child_process');
const path = require('path');
const { canAccess, newLineRegex, sort } = require('./util');

function darwin() {
  const suffixes = [
    // '/Contents/MacOS/Google Chrome Canary', 
    '/Contents/MacOS/Google Chrome', 
    // '/Contents/MacOS/Chromium'
  ];

  const LSREGISTER = '/System/Library/Frameworks/CoreServices.framework' +
    '/Versions/A/Frameworks/LaunchServices.framework' +
    '/Versions/A/Support/lsregister';

  const installations = [];

  execSync(
    `${LSREGISTER} -dump` +
    ' | grep -E -i \'(google chrome( canary)?|chromium).app$\'' +
    ' | awk \'{$1=""; print $0}\'')
    .toString()
    .split(newLineRegex)
    .forEach((inst) => {
      suffixes.forEach(suffix => {
        const execPath = path.join(inst.trim(), suffix);
        if (canAccess(execPath)) {
          installations.push(execPath);
        }
      });
    });

  // Retains one per line to maintain readability.
  const priorities = [
    // { regex: new RegExp(`^${process.env.HOME}/Applications/.*Chromium.app`), weight: 49 },
    { regex: new RegExp(`^${process.env.HOME}/Applications/.*Chrome.app`), weight: 50 },
    // { regex: new RegExp(`^${process.env.HOME}/Applications/.*Chrome Canary.app`), weight: 51 },
    // { regex: /^\/Applications\/.*Chromium.app/, weight: 99 },
    { regex: /^\/Applications\/.*Chrome.app/, weight: 100 },
    // { regex: /^\/Applications\/.*Chrome Canary.app/, weight: 101 },
    // { regex: /^\/Volumes\/.*Chromium.app/, weight: -3 },
    { regex: /^\/Volumes\/.*Chrome.app/, weight: -2 },
    // { regex: /^\/Volumes\/.*Chrome Canary.app/, weight: -1 }
  ];

  return sort(installations, priorities);
}

module.exports = darwin;
