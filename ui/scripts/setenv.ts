const { writeFile } = require('fs');

// read environment variables from .env file
require('dotenv').config();

const targetPath = `./src/environments/environment.env.ts`;

const enableMaintenace = () => {
  if (process.env.MAINTENANCE) {
    return `
  MAINTENANCE: '${ process.env.MAINTENANCE }',
    `;
  } else {
    return '';
  }
};

const enableGTAG = () => {
  if (process.env.GTAGID) {
    return `
    GTAGID: '${ process.env.GTAGID }',
    `;
  } else {
    return '';
  }
};

// we have access to our environment variables
// in the process.env object thanks to dotenv
const environmentFileContent = `
export const environment = {${enableMaintenace()}${enableGTAG()}};
`;

// write the content to the respective file
writeFile(targetPath, environmentFileContent, (err) => {
  if (err) {
    console.warn(err);
  }
  console.info(`Wrote ${ environmentFileContent } to ${ targetPath }`);
});
