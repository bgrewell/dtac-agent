## v1.1.9 (2024-01-08)

### Fix

- fix bug with control_port handling

### Refactor

- modify logging output formatter

## v1.1.8 (2023-12-14)

### Fix

- **plugins**: fix routing issue caused by lack of action in route map

## v1.1.7 (2023-11-30)

### Fix

- removed old erroneous configuration elements

## v1.1.6 (2023-11-30)

### Fix

- modify endpoint list to take user auth into context and update required auth groups

### Refactor

- enable ability to disable/enable authn/authz subsystems

## v1.1.5 (2023-11-30)

### Fix

- **validation**: move validator(s) out of api adapters and into middleware to centralize

## v1.1.4 (2023-11-30)

### Fix

- **authn**: migrated from sha256 password hashes to salted bcrypt

## v1.1.3 (2023-11-30)

### Fix

- update endpoint validation/self-describing code
- fixed missing error check

## v1.1.2 (2023-11-29)

### Fix

- **authz**: add support for users database and fully implement casbin based RBAC authorization

## v1.1.1 (2023-11-29)

### Fix

- **grpc.go**: modify to follow standard practice of using grpc request metadata for authentication

## v1.1.0 (2023-11-28)

## v1.0.8 (2023-11-27)

### Fix

- fix issue with plugin configurations

## v1.0.7 (2023-11-27)

### Fix

- **ownership_windows.go**: disable support for plugin ownership checks on Windows

### Refactor

- fix plugin writeable check

## v1.0.6 (2023-11-27)

### Refactor

- refactor endpoint struct locations and update code as needed

## v1.0.5 (2023-11-14)

### Fix

- fix plugin loading issues

### Refactor

- change formatting of some logging
- **rest.go**: refactor REST adapter to use default logging framework

## v1.0.4 (2023-11-14)

### Fix

- **config.go**: fix missing config value 'enabled' for default tls config
- **postinstall.sh**: fix typo in config pass generation

## v1.0.3 (2023-11-14)

### Fix

- modified goreleaser and scripts to handle python dependencies

## v1.0.2 (2023-11-14)

### Fix

- update post-installer scripts to install dtac-tools

## v1.0.1 (2023-11-14)

### Refactor

- fix lint issues

## v1.0.0 (2023-11-14)

### Feat

- finish plugin rework and fixes
- add logging support to plugins
- add support for endpoints to describe their output

### Refactor

- remove .venv
- moved types out of internal to pkg
- **plugins**: simplify plugin creation
- **scripts/rhel/postinstall.sh**: modify post-install

## v0.2.2 (2023-10-24)

### Refactor

- update versioning starting at 0.2.1

## v0.2.1 (2023-10-24)

### Feat

- decouple output formatter
- add random password gen
- add endpoint to expose user/group service is running under

### Fix

- fix cli cert verification
- modify to save ca cert
- fix incorrect permissions

### Refactor

- clear CA priv key memory
- modify token retreival code
- add token function to cli
- fix typo
- refactor config
