# Medik

Medik is a tool to help you run sanity checks and produce reports for your dev environment. It checks for aspects that go beyond what a `package.json`, a `devbox.json` or even a `.devcontainer.json` ensure. Medik is not meant to setup your environment, but rather check if your setup is done right.

## Features

Medik will allow you to check several aspects of your dev environment

- [x] Environment variables
  - [x] Existence
  - [x] Value
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

## `medik.yaml`

The main configuration file for your dev environment is `medik.yaml` it follows a simple structure of `protocols` where you describe checks to validate your environment. A `medik.yaml` main contain a default protocol that will be run every time `medik examine` is called and scenario specific protocols that can be run on demand

### `protocol`

A protocol is a set of properties, they might be: `vitals` or `checks`. `vitals` are **must-have** validations, if any fail the environment is considered **unhealthy**. `checks` are optional validations, they show a warn but the environment is still considered **healthy**

Protocols also have a name so you can use them in a specific scenario.

```yaml
protocols: # This list the protocols
  - name: release # A protocol to build a release of a project
    vitals:               # Checks that cannot fail
      - type: env         # Check if the following environment variables are set
        vars:
          - SECRET_KEY
      - type: env.dir     # Check if the following environment variables are set and point to valid directories
        vars:
          - BUILD_DIR
      - type: file        # Check if a file exists (relative to `medik.yaml`)
          - ./android/app/app-upload.keystore
      - type: service     # Check if a service is running on the local machine
        port: tcp:8081
  - name: test            # A protocol to check the environment to run a test suit
    vitals:
      - type: cmd
        cmd: pg_isready   # Use a custom command to check if the database is up and running
    checks:
      - type: cmd
        cmd: npx eslint   # Emit a warning if the linting fails 
```
