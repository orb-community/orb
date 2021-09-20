const { execSync, execFileSync } = require('child_process');
const path = require('path');
const fs = require('fs');
const { canAccess, sort, isExecutable, newLineRegex } = require('./util');


function findChromeExecutablesForLinuxDesktop(folder) {
  const argumentsRegex = /(^[^ ]+).*/; // Take everything up to the first space
  const chromeExecRegex = '^Exec=\/.*\/(google|chrome|chromium)-.*';

  let installations = [];
  if (canAccess(folder)) {
    // Output of the grep & print looks like:
    //    /opt/google/chrome/google-chrome --profile-directory
    //    /home/user/Downloads/chrome-linux/chrome-wrapper %U
    let execPaths;
    execPaths = execSync(`find "${folder}" -type f -exec grep -E "${chromeExecRegex}" "{}" \\; | awk -F '=' '{print $2}'`);

    execPaths = execPaths
      .toString()
      .split(newLineRegex)
      .map((execPath) => execPath.replace(argumentsRegex, '$1'));

    execPaths.forEach((execPath) => canAccess(execPath) && installations.push(execPath));
  }

  return installations;
}


/**
 * Look for linux executables in 2 ways
 * 1. Look into the directories where .desktop are saved on gnome based distro's
 * 2. Look for google-chrome-stable & google-chrome executables by using the which command
 */
function linux() {
  let installations = [];

  // 2. Look into the directories where .desktop are saved on gnome based distro's
  const desktopInstallationFolders = [
    path.join(require('os').homedir(), '.local/share/applications/'),
    '/usr/share/applications/',
  ];
  desktopInstallationFolders.forEach(folder => {
    installations = installations.concat(findChromeExecutablesForLinuxDesktop(folder));
  });

  // Look for google-chrome-stable & google-chrome executables by using the which command
  const executables = [
    'google-chrome-stable',
    'google-chrome',
    // 'chromium',
    // 'chromium-browser',
    // 'chromium/chrome',   // on toradex machines "chromium" is a directory. seen on Angstrom v2016.12
  ];
  executables.forEach((executable) => {
    // see http://tldp.org/LDP/Linux-Filesystem-Hierarchy/html/
    const validChromePaths = [
      '/usr/bin',
      '/usr/local/bin',
      '/usr/sbin',
      '/usr/local/sbin',
      '/opt/bin',
      '/usr/bin/X11',
      '/usr/X11R6/bin'
    ].map((possiblePath) => {
      try {
        const chromePathToTest = possiblePath + '/' + executable;
        if (fs.existsSync(chromePathToTest) && canAccess(chromePathToTest) && isExecutable(chromePathToTest)) {
          installations.push(chromePathToTest);
          return chromePathToTest;
        }
      } catch (err) {
        // not installed on this path or inaccessible
      }
      return undefined;
    }).filter((foundChromePath) => foundChromePath);

    // skip asking "which" command if the binary was found by searching the known paths.
    if (validChromePaths && validChromePaths.length > 0) {
      return;
    }

    try {
      const chromePath =
        execFileSync('which', [executable]).toString().split(newLineRegex)[0];
      if (canAccess(chromePath)) {
        installations.push(chromePath);
      }
    } catch (err) {
      // cmd which not installed.
    }
  });

  const priorities = [
    // { regex: /chromium$/, weight: 52 },
    { regex: /chrome-wrapper$/, weight: 51 },
    { regex: /google-chrome-stable$/, weight: 50 },
    { regex: /google-chrome$/, weight: 49 },
    // { regex: /chromium-browser$/, weight: 48 },
    { regex: /chrome$/, weight: 47 },
  ];

  return sort(Array.from(new Set(installations.filter(Boolean))), priorities);
}

module.exports = linux;
