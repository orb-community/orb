# chrome-version

Finds the version of Chrome that is installed on your machine.

## Installation
```sh
npm install @testim/chrome-version
```

## Use
```js
(async () => {
    const { getChromeVersion } = require('chrome-version');
    const version = await getChromeVersion();
    console.log(version);
})();
```

* If no version of chrome is installed on your machine `getChromeVersion` will return null.
