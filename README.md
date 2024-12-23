# Medik

Medik is a tool to help you run sanity checks and produce reports for your dev environment. It checks for aspects that go beyond what a `package.json`, a `devbox.json` or even a `.devcontainer.json` ensure. Medik is not meant to setup your environment, but rather check if your setup is done right.

## Features

Medik will allow you to check several aspects of your dev environment

- [ ] Environment variables
  - [ ] Existence
  - [ ] Value
- [ ] Files and folders
  - [ ] Withing the environment
  - [ ] On the host system
- [ ] Existing programs
  - [ ] Check versioning
  - [ ] Check location
- [ ] Services
  - [ ] Port status
  - [ ] Reachable hosts
  - [ ] Network settings
  - [ ] Running processes

Each aspect is configured via rules set in the `medik.yaml` file. The will be also a `medik.overrides.yaml` file which every developer my use to apply some overrides in order to tweak Medik to work correctly in each users personal environment.

## Docs

### Environment Variables