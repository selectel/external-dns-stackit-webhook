# Selectel Webhook - ExternalDNS

[![GoTemplate](https://img.shields.io/badge/go/template-black?logo=go)](https://github.com/golang-standards/project-layout)
[![CI](https://github.com/selectel/external-dns-stackit-webhook/actions/workflows/main.yml/badge.svg)](https://github.com/selectel/external-dns-stackit-webhook/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/selectel/external-dns-stackit-webhook)](https://goreportcard.com/report/github.com/selectel/external-dns-stackit-webhook)
[![GitHub release](https://img.shields.io/github/release/selectel/external-dns-stackit-webhook.svg)](https://github.com/selectel/external-dns-stackit-webhook/releases)
[![Last Commit](https://img.shields.io/github/last-commit/selectel/external-dns-stackit-webhook/main.svg)](https://github.com/selectel/external-dns-stackit-webhook/commits/main)
[![GitHub issues](https://img.shields.io/github/issues/selectel/external-dns-stackit-webhook.svg)](https://github.com/selectel/external-dns-stackit-webhook/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/selectel/external-dns-stackit-webhook.svg)](https://github.com/selectel/external-dns-stackit-webhook/pulls)
[![GitHub stars](https://img.shields.io/github/stars/selectel/external-dns-stackit-webhook.svg?style=social&label=Star&maxAge=2592000)](https://github.com/selectel/external-dns-stackit-webhook/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/selectel/external-dns-stackit-webhook.svg?style=social&label=Fork&maxAge=2592000)](https://github.com/selectel/external-dns-stackit-webhook/network)

ExternalDNS serves as an add-on for Kubernetes designed to automate the management of Domain Name System (DNS)
records for Kubernetes services by utilizing various DNS providers. While Kubernetes traditionally manages DNS
records internally, ExternalDNS augments this functionality by transferring the responsibility of DNS records
management to an external DNS provider such as Selectel. Consequently, the Selectel webhook enables the management
of your Selectel domains within your Kubernetes cluster using
[ExternalDNS](https://github.com/kubernetes-sigs/external-dns).

For utilizing ExternalDNS with Selectel, it is mandatory to establish a Selectel project, a service account
within the project, generate an authentication token for the service account, authorize the service account
to create and read dns zones, and finally, establish a Selectel zone.

## Kubernetes Deployment

The Selectel webhook is presented as a standard Open Container Initiative (OCI) image released in the
[GitHub container registry](https://github.com/selectel/external-dns-stackit-webhook/pkgs/container/external-dns-stackit-webhook).
The deployment is compatible with all Kubernetes-supported methods. The subsequent example
demonstrates the deployment as a
[sidecar container](https://kubernetes.io/docs/concepts/workloads/pods/#workload-resources-for-managing-pods)
within the ExternalDNS pod.

```shell 
kubectl create secret generic external-dns-webhook --from-literal=password='<Service-User-Password>'
```

```shell
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: external-dns
  namespace: default
  labels:
    app.kubernetes.io/name: external-dns
    app.kubernetes.io/instance: external-dns
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: external-dns
  labels:
    app.kubernetes.io/name: external-dns
    app.kubernetes.io/instance: external-dns
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["list","watch"]
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get","watch","list"]
  - apiGroups: [""]
    resources: ["services","endpoints"]
    verbs: ["get","watch","list"]
  - apiGroups: ["extensions","networking.k8s.io"]
    resources: ["ingresses"]
    verbs: ["get","watch","list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: external-dns-viewer
  labels:
    app.kubernetes.io/name: external-dns
    app.kubernetes.io/instance: external-dns
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: external-dns
subjects:
  - kind: ServiceAccount
    name: external-dns
    namespace: default
---
apiVersion: v1
kind: Service
metadata:
  name: external-dns
  namespace: default
  labels:
    app.kubernetes.io/name: external-dns
    app.kubernetes.io/instance: external-dns
spec:
  type: ClusterIP
  selector:
    app.kubernetes.io/name: external-dns
    app.kubernetes.io/instance: external-dns
  ports:
    - name: http
      port: 7979
      targetPort: http
      protocol: TCP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: external-dns
  namespace: default
  labels:
    app.kubernetes.io/name: external-dns
    app.kubernetes.io/instance: external-dns
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: external-dns
      app.kubernetes.io/instance: external-dns
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: external-dns
        app.kubernetes.io/instance: external-dns
    spec:
      serviceAccountName: external-dns
      securityContext:
        fsGroup: 65534
      containers:
        - name: external-dns
          securityContext:
            capabilities:
              drop:
              - ALL
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            runAsUser: 65534
          image: registry.k8s.io/external-dns/external-dns:v0.14.0
          imagePullPolicy: IfNotPresent
          args:
            - --log-level=info
            - --log-format=text
            - --interval=1m
            - --source=service
            - --source=ingress
            - --policy=sync # set it upsert-only if you don't want it to delete records
            - --provider=webhook
          ports:
            - name: http
              protocol: TCP
              containerPort: 7979
          livenessProbe:
            failureThreshold: 2
            httpGet:
              path: /healthz
              port: http
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 5
          readinessProbe:
            failureThreshold: 6
            httpGet:
              path: /healthz
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 5
        - name: webhook
          securityContext:
            capabilities:
              drop:
                - ALL
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            runAsUser: 65534
          image: ghcr.io/selectel/external-dns-stackit-webhook:v0.1.0
          imagePullPolicy: IfNotPresent
          args:
            - --project-id=c158c736-0300-4044-95c4-b7d404279b35 # your project id
            - --account-id=000000 # your account id
            - --username=Username # your service user's name
          ports:
            - name: http
              protocol: TCP
              containerPort: 8888
          livenessProbe:
            failureThreshold: 2
            httpGet:
              path: /healthz
              port: http
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 5
          readinessProbe:
            failureThreshold: 6
            httpGet:
              path: /healthz
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 5
          env:
            - name: PASSWORD
              valueFrom:
                secretKeyRef:
                  name: external-dns-webhook
                  key: password
EOF
```

## Configuration

The configuration of the Selectel webhook can be accomplished through command line arguments and environment variables.
Below are the options that are available.

- `--project-id`/`PROJECT_ID` (required): Specifies the project id to authorize.
- `--account-id`/`ACCOUNT_ID` (required): Specifies the account id to authorize.
- `--username`/`USERNAME` (required): Specifies the username of your service user to authorize.
- `--password`/`PASSWORD` (required): Specifies the password of your service user to authorize.
- `--worker`/`WORKER`  (optional): Specifies the number of workers to employ for querying the API. Given that we
  need to iterate over all zones and records, it can be parallelized. However, it is important to avoid
  setting this number excessively high to prevent receiving 429 rate limiting from the API (default 10).
- `--base-url`/`BASE_URL` (optional): Identifies the Base URL for utilizing the API (
  default "https://api.selectel.ru/domains/v2").
- `--api-port`/`API_PORT` (optional): Specifies the port to listen on (default 8888).
- `--domain-filter`/`DOMAIN_FILER` (optional): Establishes a filter for DNS zone names (default []).
- `--dry-run`/`DRY_RUN` (optional): Specifies whether to perform a dry run (default false).
- `--log-level`/`LOG_LEVEL` (optional): Defines the log level (default "info"). Possible values are: debug, info, warn,
  error.

## Development

Run the app:

```bash
export PROJECT_ID="c158c736-0300-4044-95c4-b7d404279b35"
export ACCOUNT_ID="123456"
export USERNAME="username"
export PASSWORD ="password"

make run
```

Lint the code:

```bash
make lint
```

Test the code:

```bash
make test
```
