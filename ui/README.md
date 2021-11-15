# <a href="https://github.com/ns1labs/orb">Orb</a> UI

> This wiki is targeted at developers looking to build the UI from source, to either
> run it locally for development purposes or to create a custom UI docker image.

## Development

### Prerequisites

The following are needed to run the UI:

* [node](https://nodejs.org/en/blog/release/v12.21.0/)
* [npm](https://github.com/npm/cli/tree/v7.22.0)

*It is recomended to build the UI using [yarn](https://www.npmjs.com/package/yarn)*

### Install

```bash
# note: if you haven't checked out the full repo yet, and you're only interested in developing 
# the front end locally, you can do so by checking out only the ui folder.
# [read more...](
git clone git@github.com:ns1labs/orb.git --no-checkout --depth 1 ${path}

# however you clone the project
cd ${path}/ui
yarn install
```

### Usage

A developer build from the source can be achieved using the following command:

```bash
yarn build
```

*(Check [package.json](./package.json) file for available tasks.)*


### Development  

While developing, it is useful to serve UI locally and have your changes to the code having effect immediately.

The commands `yarn start` and `yarn start:withmock` will generate a dev build and serve it at `http://localhost:4200/`.

*(Note that `http://localhost:4200/` is for development use only, and is not intended to be used by the end-user.)*

*([proxy-config.json](./proxy-config.json) re-routes all outbound requests when running local serve task)

> When running develop builds locally with ng serve, it might be necessary to run the following commands to be able to watch the sum of files required.
> [stackoverflow inotify](https://stackoverflow.com/a/55763478)  

```sh
# insert the new value into the system config
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p

# check that the new value was applied
cat /proc/sys/fs/inotify/max_user_watches

# config variable name (not runnable)
fs.inotify.max_user_watches=524288
```

---

## Entities Data examples

See [data examples](https://github.com/ns1labs/orb/wiki/Orb-UI---Entities-Data-Examples) for examples of *Orb Entities*
to aid in UI design , form validation and related info.


---  

## QA & Testing

Quality Assurance & Test frameworks and scripts are still a *WORK IN PROGRESS*  
Check our [Wiki](https://github.com/ns1labs/orb/wiki/UI-QA-Automation-Tags) for more information.
