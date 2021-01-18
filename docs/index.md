# residential-dns

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square)

Residential DNS

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| deployment.args | string | `nil` | The list of arguments to pass to the controller. |
| deployment.env | string | `nil` | The list of environment variables to set on the pod spec template. |
| deployment.imageTag | string | `"v0.0.0"` | The residential-dns Docker image tag to run. |
| deployment.podAnnotations | string | `nil` | A map of annotations to add to the pod spec template. |
| serviceAccount.annotations | string | `nil` | A map of annotations to add to the ServiceAccount. |
| serviceAccount.create | bool | `true` | Whether to create a ServiceAccount object and associate it with the Deployment. |
| serviceAccount.name | string | `"residential-dns"` | The name to use for the ServiceAccount object. |
