const chromedriver = require('chromedriver');

const testUrl = 'http://localhost';
const defaultTimeout = 15000;

module.exports = {
  src_folders : ['tests'],
  page_objects_path : './page_objects',
  skip_testcases_on_fail: false,

  webdriver: {
    start_process: true,
  },

  test_settings: {
    default: {
      launch_url: testUrl,
      globals: {
        waitForConditionTimeout: defaultTimeout,
      },
      webdriver: {
        server_path: chromedriver.path,
        port: 9515,
      },
      desiredCapabilities: {
        browserName: 'chrome',
      }
    },

    headless: {
      launch_url: testUrl,
      globals: {
        waitForConditionTimeout: defaultTimeout,
      },
      webdriver: {
        server_path: chromedriver.path,
        port: 9515,
      },
      desiredCapabilities: {
        browserName: 'chrome',
        chromeOptions: {
          w3c: false,
          args: ['--headless', '--no-sandbox', '--disable-dev-shm-usage'],
        }
      }
    },

    beta: {
      launch_url: "http://beta.orb.live"
    }
  }
}
