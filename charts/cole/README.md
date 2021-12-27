# cole

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

A Helm chart for Cole

## Get Repo Info

```console
helm repo add cole https://nicolastakashi.github.io/cole
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release nicolastakashi/cole
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Values

| Key                                     | Type   | Default                    | Description                                                                                  |
|-----------------------------------------|--------|----------------------------|----------------------------------------------------------------------------------------------|
| affinity                                | object | `{}`                       |                                                                                              |
| flags.grafana.containerName             | string | `"grafana"`                |                                                                                              |
| flags.grafana.log.format                | string | `"console"`                |                                                                                              |
| flags.grafana.namespace                 | string | `"grafana"`                |                                                                                              |
| flags.grafana.podLabelselector[0].name  | string | `"app.kubernetes.io/name"` |                                                                                              |
| flags.grafana.podLabelselector[0].value | string | `"grafana"`                |                                                                                              |
| flags.kubeconfig                        | string | `""`                       |                                                                                              |
| flags.log.level                         | string | `"debug"`                  |                                                                                              |
| flags.metrics.includeUname              | bool   | `false`                    |                                                                                              |
| fullnameOverride                        | string | `""`                       |                                                                                              |
| image.pullPolicy                        | string | `"Always"`                 |                                                                                              |
| image.repository                        | string | `"ntakashi/cole"`          |                                                                                              |
| image.tag                               | string | `"1.0.0"`                  |                                                                                              |
| imagePullSecrets                        | list   | `[]`                       |                                                                                              |
| nameOverride                            | string | `""`                       |                                                                                              |
| nodeSelector                            | object | `{}`                       |                                                                                              |
| podAnnotations                          | object | `{}`                       |                                                                                              |
| podSecurityContext                      | object | `{}`                       |                                                                                              |
| resources                               | object | `{}`                       |                                                                                              |
| securityContext.readOnlyRootFilesystem  | bool   | `true`                     |                                                                                              |
| service.port                            | int    | `80`                       |                                                                                              |
| service.type                            | string | `"ClusterIP"`              |                                                                                              |
| serviceAccount.annotations              | object | `{}`                       |                                                                                              |
| serviceAccount.create                   | bool   | `true`                     |                                                                                              |
| serviceAccount.name                     | string | `""`                       | If not set and create is true, a name is generated using the fullname template               |
| serviceMonitor.enabled                  | bool   | `false`                    |                                                                                              |
| serviceMonitor.interval                 | string | `""`                       | ref: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#endpoint |
| serviceMonitor.scrapeTimeout            | string | `""`                       | ref: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#endpoint |
| tolerations                             | list   | `[]`                       |                                                                                              |

