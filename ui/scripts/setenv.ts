const { writeFile } = require('fs');

// read environment variables from .env file
require('dotenv').config();

const targetPath = `./src/environments/environment.env.ts`;

// we have access to our environment variables
// in the process.env object thanks to dotenv
const environmentFileContent = `
export const environment = {
  PS: {
    // site id
    SID: '${process.env.PS_SID}',
    // group key
    GROUP_KEY: '${process.env.PS_GROUP_KEY}',
  },
};

`;

// write the content to the respective file
writeFile(targetPath, environmentFileContent, (err) => {
  if (err) {
    console.warn(err);
  }
  console.info(`Wrote variables to ${targetPath}`);
});
