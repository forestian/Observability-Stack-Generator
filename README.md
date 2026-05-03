# observability-stack-generator

`obsgen` is an MVP Go CLI that generates starter Kubernetes observability stack configuration files.

It creates readable Helm values and helper scripts for:

- Grafana Loki
- Grafana Mimir
- Grafana Tempo
- Grafana Alloy
- MinIO or generic S3-compatible object storage

The tool is a generator only. It does not call Kubernetes, run Helm, create cloud resources, generate real secrets, or deploy anything directly.

## Usage

```sh
go run . init --name demo --namespace monitoring --storage minio --profile dev --output ./demo-stack
```

```sh
go run . init --name demo-s3 --namespace monitoring --storage s3 --profile production --output ./demo-s3-stack
```

```sh
go run . version
```

## Install from GitHub Releases

Download a prebuilt binary from the GitHub Releases page.

Linux/macOS:

```sh
tar -xzf <archive>.tar.gz
chmod +x obsgen
./obsgen version
```

Windows:

Download the Windows archive, extract it, and run:

```powershell
obsgen.exe version
```

## Commands

### `obsgen init`

Flags:

- `--name` default `observability-stack`
- `--namespace` default `monitoring`
- `--output` default `./observability-stack`
- `--storage` default `minio`, allowed `minio` or `s3`
- `--profile` default `dev`, allowed `dev` or `production`
- `--force` overwrite generated files if the output directory already exists

When `--profile production` is used without an explicit `--storage` flag, the generated stack defaults to S3-compatible object storage.

## Development

```sh
go test ./...
```
