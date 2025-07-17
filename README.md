# nfpm-helper

[nfpm](https://nfpm.goreleaser.com/) helper for packaging software from existing releases. 

## Build configuration

Filename: `nfpm-helper.yml`

Example:

```
name: oauth2-proxy
download:
  url_template: https://github.com/oauth2-proxy/oauth2-proxy/releases/download/v${VERSION}/${NAME}-v${VERSION}.linux-${ARCH}.tar.gz
strip_components: 1
outputs:
  - arch: amd64
  - arch: arm64
```

### Config

| Name               | Type                  | Required | Default | Description                |
|--------------------|-----------------------|----------|---------|----------------------------|
| `name`             | `string`              | Y        |         | Package name               |
| `download`         | `DownloadBaseConfig`  | N        |         | Download base config       |
| `strip_components` | `int`                 | N        | 0       | Strip components from path |
| `packaging`        | `PackagingBaseConfig` | N        |         | Packaging base config      |
| `outputs`          | `[]Output`            | Y        |         | Outputs                    |

### DownloadBaseConfig

| Name           | Type     | Required | Default | Description           |
|----------------|----------|----------|---------|-----------------------|
| `url_template` | `string` | N*       |         | Download URL template |

### PackagingBaseConfig

| Name                | Type     | Required | Default                       | Description       |
|---------------------|----------|----------|-------------------------------|-------------------|
| `filename_template` | `string` | N        | `${NAME}_${VERSION}_${ARCH}`  | Filename template |

### Output

| Name        | Type              | Required | Default | Description      |
|-------------|-------------------|----------|---------|------------------|
| `arch`      | `string`          | Y        |         | Arch             |
| `download`  | `DownloadConfig`  | N        |         | Download config  |
| `packaging` | `PackagingConfig` | N        |         | Packaging config |

### DownloadConfig

| Name           | Type                | Required | Default | Description           |
|----------------|---------------------|----------|---------|-----------------------|
| `url_template` | `string`            | N*       |         | Download URL template |
| `env`          | `map[string]string` | N        |         | Environment override  |

### PackagingConfig

| Name                | Type                | Required | Default | Description          |
|---------------------|---------------------|----------|---------|----------------------|
| `filename_template` | `string`            | N        |         | Filename template    |
| `env`               | `map[string]string` | N        |         | Environment override |

### Variables

Used in templates

| Name          | Description                  |
|---------------|------------------------------|
| `NAME`        | Package name (`Config.name`) |
| `VERSION`     | Package version              |
| `ARCH`        | Package arch                 |
| `ARCHIVE_DIR` | Unpacked archive directory.  |

### Arch

Follows `GOARCH`

Examples: amd64, 386, arm5, arm6, arm7, arm64, mips, mips64, mips64le, mipsle, ppc64, ppc64le, loong64, riscv64, s390x)

## Index configuration

Filename: `nfpm-helper.index.yml`

Example:

```
packages:
  - name: gost
    dir: packages/gost
  - name: oauth2-proxy
    dir: packages/oauth2-proxy
  - name: ory/hydra
    dir: packages/ory/hydra
  - name: ory/oathkeeper
    dir: packages/ory/oathkeeper
  - name: prometheus/node_exporter
    dir: packages/prometheus/node_exporter
  - name: starship
    dir: packages/starship
```

### Config

| Name       | Type      | Required | Default | Description |
|------------|-----------|----------|---------|-------------|
| `packages` | `[]Entry` | Y        |         | Packages    |

### Entry

| Name   | Type     | Required | Default | Description                                                    |
|--------|----------|----------|---------|----------------------------------------------------------------|
| `name` | `string` | Y        |         | Package name (for external reference)                          |
| `dir`  | `string` | Y        |         | Path to build context (directory containing `nfpm-helper.yml`) |

## Generate configuration

Filename: `nfpm-helper.gen.yml`

Example:
```
repositories:
  - source: https://github.com/ngyewch/nfpm-packaging.git
    type: git
    packages:
      - name: gost
        version: 3.1.0
        archs: 
          - amd64
      - name: oauth2-proxy
        version: 7.9.0
        archs: 
          - amd64
      - name: ory/hydra
        version: 2.3.0
        archs: 
          - amd64
      - name: ory/oathkeeper
        version: 0.40.9
        archs: 
          - amd64
      - name: prometheus/node_exporter
        version: 1.9.1
        archs: 
          - amd64
```

### Config

| Name           | Type                 | Required | Default | Description  |
|----------------|----------------------|----------|---------|--------------|
| `repositories` | `[]RepositoryConfig` | Y        |         | Repositories |

### RepositoryConfig

| Name       | Type              | Required | Default | Description                                                               |
|------------|-------------------|----------|---------|---------------------------------------------------------------------------|
| `source`   | `string`          | Y        |         | Source.                                                                   |
| `version`  | `string`          | N        |         | If `type` = `git`, `version` = branch, tag or commit hash.                |
| `type`     | `string`          | N        | `local` | If `local`, `source` is a path. If `git`, source is a git repository URL. |  
| `packages` | `[]PackageConfig` | Y        |         | Packages                                                                  |

### PackageConfig

| Name        | Type       | Required | Default | Description                           |
|-------------|------------|----------|---------|---------------------------------------|
| `name`      | `string`   | Y        |         | Package name as defined in the index. |
| `version`   | `string`   | Y        |         | Package version.                      |
| `archs`     | `[]string` | Y        |         | Archs to generate.                    |
| `packagers` | `[]string` | N        |         | Packagers to use.                     |
