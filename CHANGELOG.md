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
