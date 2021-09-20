"use strict";

const findChrome = require('./chrome-finder');
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const path = require('path');
const readdir = util.promisify(require('fs').readdir);
const { execSync } = require('child_process');

async function getChromeVersionFromCli() {

    let chromePath;
    try {
        chromePath = findChrome();
    } catch (err) {
        return null;
    }

    const res = await exec(chromePath.replace(/ /g, '\\ ') + ' --version');

    const version = res.stdout.substr(14).trim();
    return version;

}

async function getChromeVersionWin() {
    
    let chromePath;
    try {
        chromePath = findChrome();
    } catch (err) {
        return null;
    }

    const versionPath = path.dirname(chromePath);

    const contents = await readdir(versionPath);

    const versions = contents.filter(a => /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/g.test(a));

    // returning oldest in case there is an updated version and chrome still hasn't relaunched
    const oldest = versions.sort((a,b) => a>b)[0];

    return oldest;
   
}

function getChromeVersionFromOsa() {
    try {
        const version = execSync('osascript -e \'tell application "Google Chrome" to get version\'').toString().trim();
        return version;
    } catch (err) {
        return null;
    }
}

async function getChromeVersion() {

    const os = process.platform;
    
    if (os === 'darwin') return getChromeVersionFromOsa();
    if (os === 'linux') return getChromeVersionFromCli();
    if (os.includes('win')) return getChromeVersionWin();

    console.log(`${os} is not supported`);

    return null;

}

if (require.main == module) {
    getChromeVersion().then(v => console.log(v));
}

module.exports = {
    getChromeVersion 
};