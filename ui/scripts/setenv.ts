const { writeFile } = require('fs');

// read environment variables from .env file
require('dotenv').config();

const targetPath = `./src/environments/environment.env.ts`;

// PACTSAFE
const enablePS = () => {
  if (process.env.PS_SID !== '' && process.env.PS_GROUP_KEY !== '') {
    return `
  PS: {
    // site id
    SID: '${ process.env.PS_SID }',
    // group key
    GROUP_KEY: '${ process.env.PS_GROUP_KEY }',
  },
`;
  } else {
    return ``;
  }
};

// we have access to our environment variables
// in the process.env object thanks to dotenv
const environmentFileContent = `
export const environment = {
  ${enablePS()}
};
`;

// write the content to the respective file
writeFile(targetPath, environmentFileContent, (err) => {
  if (err) {
    console.warn(err);
  }
  console.info(`Wrote ${ environmentFileContent } to ${ targetPath }`);
});
