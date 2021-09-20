
(async () => {
    const { getChromeVersion } = require('./index');
    const version = await getChromeVersion();
    console.log(version);
})();