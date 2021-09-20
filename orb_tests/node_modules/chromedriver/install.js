'use strict';
// @ts-check

const fs = require('fs');
const helper = require('./lib/chromedriver');
const axios = require('axios').default;
const path = require('path');
const del = require('del');
const child_process = require('child_process');
const os = require('os');
const url = require('url');
const https = require('https');
const extractZip = require('extract-zip');
const { getChromeVersion } = require('@testim/chrome-version');
const HttpsProxyAgent = require('https-proxy-agent');
const getProxyForUrl = require("proxy-from-env").getProxyForUrl;

const skipDownload = process.env.npm_config_chromedriver_skip_download || process.env.CHROMEDRIVER_SKIP_DOWNLOAD;
if (skipDownload === 'true') {
  console.log('Found CHROMEDRIVER_SKIP_DOWNLOAD variable, skipping installation.');
  process.exit(0);
}

const libPath = path.join(__dirname, 'lib', 'chromedriver');
let cdnUrl = process.env.npm_config_chromedriver_cdnurl || process.env.CHROMEDRIVER_CDNURL || 'https://chromedriver.storage.googleapis.com';
const configuredfilePath = process.env.npm_config_chromedriver_filepath || process.env.CHROMEDRIVER_FILEPATH;

// adapt http://chromedriver.storage.googleapis.com/
cdnUrl = cdnUrl.replace(/\/+$/, '');
const platform = validatePlatform();
const detect_chromedriver_version = process.env.npm_config_detect_chromedriver_version || process.env.DETECT_CHROMEDRIVER_VERSION;
let chromedriver_version = process.env.npm_config_chromedriver_version || process.env.CHROMEDRIVER_VERSION || helper.version;
let chromedriverBinaryFilePath;
let downloadedFile = '';

(async function install() {
  try {
    if (detect_chromedriver_version === 'true') {
      // Refer http://chromedriver.chromium.org/downloads/version-selection
      const chromeVersion = await getChromeVersion();
      console.log("Your Chrome version is " + chromeVersion);
      const chromeVersionWithoutPatch = /^(.*?)\.\d+$/.exec(chromeVersion)[1];
      await getChromeDriverVersion(getRequestOptions(cdnUrl + '/LATEST_RELEASE_' + chromeVersionWithoutPatch));
      console.log("Compatible ChromeDriver version is " + chromedriver_version);
    }
    if (chromedriver_version === 'LATEST') {
      await getChromeDriverVersion(getRequestOptions(`${cdnUrl}/LATEST_RELEASE`));
    } else {
      const latestReleaseForVersionMatch = chromedriver_version.match(/LATEST_(\d+)/);
      if (latestReleaseForVersionMatch) {
        const majorVersion = latestReleaseForVersionMatch[1];
        await getChromeDriverVersion(getRequestOptions(`${cdnUrl}/LATEST_RELEASE_${majorVersion}`));
      }
    }
    const tmpPath = findSuitableTempDirectory();
    const chromedriverBinaryFileName = process.platform === 'win32' ? 'chromedriver.exe' : 'chromedriver';
    chromedriverBinaryFilePath = path.resolve(tmpPath, chromedriverBinaryFileName);
    const chromedriverIsAvailable = await verifyIfChromedriverIsAvailableAndHasCorrectVersion();
    if (!chromedriverIsAvailable) {
      console.log('Current existing ChromeDriver binary is unavailable, proceeding with download and extraction.');
      await downloadFile(tmpPath);
      await extractDownload(tmpPath);
    }
    await copyIntoPlace(tmpPath, libPath);
    fixFilePermissions();
    console.log('Done. ChromeDriver binary available at', helper.path);
  } catch (err) {
    console.error('ChromeDriver installation failed', err);
    process.exit(1);
  }
})();

function validatePlatform() {
  /** @type string */
  let thePlatform = process.platform;
  if (thePlatform === 'linux') {
    if (process.arch === 'arm64' || process.arch === 'x64') {
      thePlatform += '64';
    } else {
      console.log('Only Linux 64 bits supported.');
      process.exit(1);
    }
  } else if (thePlatform === 'darwin' || thePlatform === 'freebsd') {
    if (process.arch === 'x64' || process.arch === 'arm64') {
      thePlatform = 'mac64';
    } else {
      console.log('Only Mac 64 bits supported.');
      process.exit(1);
    }
  } else if (thePlatform !== 'win32') {
    console.log('Unexpected platform or architecture:', process.platform, process.arch);
    process.exit(1);
  }
  return thePlatform;
}

async function downloadFile(dirToLoadTo) {
  if (detect_chromedriver_version !== 'true' && configuredfilePath) {
    downloadedFile = configuredfilePath;
    console.log('Using file: ', downloadedFile);
    return;
  } else {
    const fileName = `chromedriver_${platform}.zip`;
    const tempDownloadedFile = path.resolve(dirToLoadTo, fileName);
    downloadedFile = tempDownloadedFile;
    const formattedDownloadUrl = `${cdnUrl}/${chromedriver_version}/${fileName}`;
    console.log('Downloading from file: ', formattedDownloadUrl);
    console.log('Saving to file:', downloadedFile);
    await requestBinary(getRequestOptions(formattedDownloadUrl), downloadedFile);
  }
}

function verifyIfChromedriverIsAvailableAndHasCorrectVersion() {
  if (!fs.existsSync(chromedriverBinaryFilePath))
    return Promise.resolve(false);
  const forceDownload = process.env.npm_config_chromedriver_force_download === 'true' || process.env.CHROMEDRIVER_FORCE_DOWNLOAD === 'true';
  if (forceDownload)
    return Promise.resolve(false);
  console.log('ChromeDriver binary exists. Validating...');
  const deferred = new Deferred();
  try {
    fs.accessSync(chromedriverBinaryFilePath, fs.constants.X_OK);
    const cp = child_process.spawn(chromedriverBinaryFilePath, ['--version']);
    let str = '';
    cp.stdout.on('data', data => str += data);
    cp.on('error', () => deferred.resolve(false));
    cp.on('close', code => {
      if (code !== 0)
        return deferred.resolve(false);
      const parts = str.split(' ');
      if (parts.length < 3)
        return deferred.resolve(false);
      if (parts[1].startsWith(chromedriver_version)) {
        console.log(`ChromeDriver is already available at '${chromedriverBinaryFilePath}'.`);
        return deferred.resolve(true);
      }
      deferred.resolve(false);
    });
  }
  catch (error) {
    deferred.resolve(false);
  }
  return deferred.promise;
}

function findSuitableTempDirectory() {
  const now = Date.now();
  const candidateTmpDirs = [
    process.env.npm_config_tmp,
    process.env.XDG_CACHE_HOME,
    // Platform specific default, including TMPDIR/TMP/TEMP env
    os.tmpdir(),
    path.join(process.cwd(), 'tmp')
  ];

  for (let i = 0; i < candidateTmpDirs.length; i++) {
    if (!candidateTmpDirs[i]) continue;
    // Prevent collision with other versions in the dependency tree
    const namespace = chromedriver_version;
    const candidatePath = path.join(candidateTmpDirs[i], namespace, 'chromedriver');
    try {
      fs.mkdirSync(candidatePath, { recursive: true });
      const testFile = path.join(candidatePath, now + '.tmp');
      fs.writeFileSync(testFile, 'test');
      fs.unlinkSync(testFile);
      return candidatePath;
    } catch (e) {
      console.log(candidatePath, 'is not writable:', e.message);
    }
  }
  console.error('Can not find a writable tmp directory, please report issue on https://github.com/giggio/chromedriver/issues/ with as much information as possible.');
  process.exit(1);
}

function getRequestOptions(downloadPath) {
  /** @type import('axios').AxiosRequestConfig */
  const options = { url: downloadPath, method: "GET" };
  const urlParts = url.parse(downloadPath);
  const isHttps = urlParts.protocol === 'https:';
  const proxyUrl = getProxyForUrl(downloadPath);

  if (proxyUrl) {
    const proxyUrlParts = url.parse(proxyUrl);
    options.proxy = {
      host: proxyUrlParts.hostname,
      port: proxyUrlParts.port ? parseInt(proxyUrlParts.port) : 80,
      protocol: proxyUrlParts.protocol
    };
  }

  if (isHttps) {
    // Use certificate authority settings from npm
    let ca = process.env.npm_config_ca;
    if (ca)
      console.log('Using npmconf ca.');

    if (!ca && process.env.npm_config_cafile) {
      try {
        ca = fs.readFileSync(process.env.npm_config_cafile, { encoding: 'utf8' });
      } catch (e) {
        console.error('Could not read cafile', process.env.npm_config_cafile, e);
      }
      console.log('Using npmconf cafile.');
    }

    if (proxyUrl) {
      console.log('Using workaround for https-url combined with a proxy.');
      const httpsProxyAgentOptions = url.parse(proxyUrl);
      // @ts-ignore
      httpsProxyAgentOptions.ca = ca;
      // @ts-ignore
      httpsProxyAgentOptions.rejectUnauthorized = !!process.env.npm_config_strict_ssl;
      // @ts-ignore
      options.httpsAgent = new HttpsProxyAgent(httpsProxyAgentOptions);
      options.proxy = false;
    } else {
      options.httpsAgent = new https.Agent({
        rejectUnauthorized: !!process.env.npm_config_strict_ssl,
        ca: ca
      });
    }
  }

  // Use specific User-Agent
  if (process.env.npm_config_user_agent) {
    options.headers = { 'User-Agent': process.env.npm_config_user_agent };
  }

  return options;
}

/**
 *
 * @param {import('axios').AxiosRequestConfig} requestOptions
 */
async function getChromeDriverVersion(requestOptions) {
  console.log('Finding Chromedriver version.');
  const response = await axios(requestOptions);
  chromedriver_version = response.data.trim();
  console.log(`Chromedriver version is ${chromedriver_version}.`);
}

/**
 *
 * @param {import('axios').AxiosRequestConfig} requestOptions
 * @param {string} filePath
 */
async function requestBinary(requestOptions, filePath) {
  const outFile = fs.createWriteStream(filePath);
  let response;
  try {
    response = await axios({ responseType: 'stream', ...requestOptions });
  } catch (error) {
    if (error && error.response) {
      if (error.response.status)
        console.error('Error status code:', error.response.status);
      if (error.response.data) {
        error.response.data.on('data', data => console.error(data.toString('utf8')));
        await new Promise((resolve) => {
          error.response.data.on('finish', resolve);
          error.response.data.on('error', resolve);
        });
      }
    }
    throw new Error('Error with http(s) request: ' + error);
  }
  let count = 0;
  let notifiedCount = 0;
  response.data.on('data', data => {
    count += data.length;
    if ((count - notifiedCount) > 1024 * 1024) {
      console.log('Received ' + Math.floor(count / 1024) + 'K...');
      notifiedCount = count;
    }
  });
  response.data.on('end', () => console.log('Received ' + Math.floor(count / 1024) + 'K total.'));
  const pipe = response.data.pipe(outFile);
  await new Promise((resolve, reject) => {
    pipe.on('finish', resolve);
    pipe.on('error', reject);
  });
}

async function extractDownload(dirToExtractTo) {
  if (path.extname(downloadedFile) !== '.zip') {
    fs.copyFileSync(downloadedFile, chromedriverBinaryFilePath);
    console.log('Skipping zip extraction - binary file found.');
    return;
  }
  console.log(`Extracting zip contents to ${dirToExtractTo}.`);
  try {
    await extractZip(path.resolve(downloadedFile), { dir: dirToExtractTo });
  } catch (error) {
    throw new Error('Error extracting archive: ' + error);
  }
}

async function copyIntoPlace(originPath, targetPath) {
  await del(targetPath, { force: true });
  console.log(`Copying from ${originPath} to target path ${targetPath}`);
  fs.mkdirSync(targetPath);

  // Look for the extracted directory, so we can rename it.
  const files = fs.readdirSync(originPath, { withFileTypes: true })
    .filter(dirent => dirent.isFile() && dirent.name.startsWith('chromedriver') && !dirent.name.endsWith(".debug") && !dirent.name.endsWith(".zip"))
    .map(dirent => dirent.name);
  const promises = files.map(name => {
    return new Promise((resolve) => {
      const file = path.join(originPath, name);
      const reader = fs.createReadStream(file);
      const targetFile = path.join(targetPath, name);
      const writer = fs.createWriteStream(targetFile);
      writer.on("close", () => resolve());
      reader.pipe(writer);
    });
  });
  await Promise.all(promises);
}


function fixFilePermissions() {
  // Check that the binary is user-executable and fix it if it isn't (problems with unzip library)
  if (process.platform != 'win32') {
    const stat = fs.statSync(helper.path);
    // 64 == 0100 (no octal literal in strict mode)
    if (!(stat.mode & 64)) {
      console.log('Fixing file permissions.');
      fs.chmodSync(helper.path, '755');
    }
  }
}

function Deferred() {
  this.resolve = null;
  this.reject = null;
  this.promise = new Promise(function (resolve, reject) {
    this.resolve = resolve;
    this.reject = reject;
  }.bind(this));
  Object.freeze(this);
}
