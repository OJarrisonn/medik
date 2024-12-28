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

A protocol is a set of exams divided between: `vitals` and `checks`. `vitals` are **must-have** validations, if any fail the environment is considered **unhealthy**. `checks` are optional validations, they show a warning if the validation fails but the environment is still considered **healthy**

Protocols are defined by name under `protocols` section.

```yaml
protocols: # This list the protocols
  release:                # A protocol to build a release of a project
    vitals:               # Checks that cannot fail
      - exam: env.is-set  # Check if the following environment variables are set
        vars:
          - SECRET_KEY
      - exam: env.dir     # Check if the following environment variables are set and point to valid directories
        vars:
          - BUILD_DIR
      - exam: file.exists # Check if a file exists (relative to `medik.yaml`)
        paths:
          - ./android/app/app-upload.keystore
      - exam: service.is-up  # Check if a service is running on the local machine
        ports:
          - tcp:8081
  test:                   # A protocol to check the environment to run a test suit
    vitals:
      - exam: cmd.custom
        cmd:
          run: pg_isready   # Use a custom command to check if the database is up and running
          exit-code: 0      # The command should exit with code 0
    checks:
      - exam: cmd.custom
        cmd:
          run: npx eslint   # Emit a warning if the linting fails
          exit-code: 0      # The command should exit with code 0
```

### Exams

Exams are the actual checks that Medik will run on your environment. They are defined by a type and a set of parameters. The following exams are available:

### `env`

The set of exams related to environment variables. The field `vars` is a list of environment variables to check and is mandatory for all of the below listed exams.

- `env.is-set`: Check if an environment variable is set
- `env.not-empty`: Check if an environment variable is set and not empty
- `env.regex`: Check if an environment variable is set and matches a regular expression
  - `regex`: The regular expression to match
- `env.options`: Check if an environment variable is set and matches one of the given options
  - `options`: A list of options to match
- `env.int`: Check if an environment variable is set and is an integer
- `env.float`: Check if an environment variable is set and is a float
- `env.int-range`: Check if an environment variable is set and is an integer within a range
  - `min`: The minimum value (inclusive)
  - `max`: The maximum value (inclusive)
- `env.float-range`: Check if an environment variable is set and is a float within a range
  - `min`: The minimum value (inclusive)
  - `max`: The maximum value (inclusive)
- `env.file`: Check if an environment variable is set and points to an existing file
- `env.dir`: Check if an environment variable is set and points to an existing directory
- `env.ipv4`: Check if an environment variable is set and is a valid IPv4 address
- `env.ipv6`: Check if an environment variable is set and is a valid IPv6 address
- `env.ip`: Check if an environment variable is set and is a valid IP address (IPv4 or IPv6)
- `env.hostname`: Check if an environment variable is set and is a valid hostname (valid url)

### `file`

> This is work in progress, not implemented yet

The set of exams related to files and folders. The field `paths` is a list of path to be checked and is mandatory for all of the below listed exams.

- [ ] `file.exists`: Check if a file exists
- [ ] `file.not-exists`: Check if a file does not exist
- [ ] `file.is-empty`: Check if a file is empty
- [ ] `file.is-not-empty`: Check if a file is not empty
- [ ] `file.is-hidden`: Check if a file is hidden
- [ ] `file.is-file`: Check if a path points to a file
- [ ] `file.is-dir`: Check if a path points to a directory
- [ ] `file.is-symlink`: Check if a path points to a symlink
- [ ] `file.is-socket`: Check if a path points to a socket
- [ ] `file.is-pipe`: Check if a path points to a pipe
- [ ] `file.is-char-device`: Check if a path points to a character device
- [ ] `file.is-block-device`: Check if a path points to a block device
- [ ] `file.is-executable`: Check if a file is executable
- [ ] `file.is-readable`: Check if a file is readable
- [ ] `file.is-writable`: Check if a file is writable

### `service`

> This is work in progress, not implemented yet

The set of exams related to services running on the machine. The field `ports` is a list of ports to check and is mandatory for all of the below listed exams.

A port can be defined as `<protocol>:<port>` if running locally or `<protocol>://<host>:<port>` if running on a remote. Where `protocol` is either `tcp` or `udp`

- [ ] `service.is-up`: Check if a service is running
- [ ] `service.is-down`: Check if a service is not running
- [ ] `service.is-reachable`: Check if a service is reachable
- [ ] `service.is-not-reachable`: Check if a service is not reachable
- [ ] `service.is-listening`: Check if a service is listening
- [ ] `service.is-not-listening`: Check if a service is not listening

### `bin.exists`

> This is work in progress, not implemented yet

This exam checks if a binary is available in the system. The field `bin` is mandatory.



```yaml

### `cmd.custom`

> This is work in progress, not implemented yet

This is a generic exam that allows you to run a custom command and check its output. The field `cmd` is mandatory. The command in the field `cmd.run` will be run in the shell. Several validations can be applied to the output of the command:

- `exit-code`: Check if the exit code of the command is the expected one
- `stdout-contains`: Check if the stdout contains a given string
- `stdout-not-contains`: Check if the stdout does not contain a given string
- `stdout-regex`: Check if the stdout matches a regular expression
- `stderr-contains`: Check if the stderr contains a given string
- `stderr-not-contains`: Check if the stderr does not contain a given string
- `stderr-regex`: Check if the stderr matches a regular expression

Using multiple validations will require all of them to pass for the exam to be successful.

```yaml
protocols:
  test:
    vitals:
      - exam: cmd.custom
        cmd:
          run: pg_isready
          exit-code: 0
          stdout-contains: accepting connections
```
