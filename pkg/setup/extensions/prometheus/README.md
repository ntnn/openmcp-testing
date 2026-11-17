# Prometheus Extension

The prometheus extension can be used to perfomance test [service providers](#service-provider-testing-example) and [domain services](#domain-service-testing-example) in context of OpenControlPlane. It deploys a dedicated test prometheus instance and enables performance testing based on regular [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) queries.

Tests can be executed by bootstrapping an OpenControlPlane environment or [against an existing environment](#testing-against-existing-environments).

## Metrics Overview

The prometheus extension is configured to collect the following metrics per target cluster (platform and optional remote clusters):

- metrics provided by an application like [controller-runtime metrics](https://book.kubebuilder.io/reference/metrics-reference.html) and [prometheus go client collector metrics](https://github.com/prometheus/client_golang)
- metrics provided by the [kubelet](https://kubernetes.io/docs/reference/instrumentation/metrics/) (/metrics/resource)

## Query Parameters Mapping

Each `prometheus.Query` call results in Kubernetes job that executes promtool [range queries](https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries). The following parameters exit:

- (required) Window defines the start, end and step settings of the range query. The derived duration can also be used as placeholder in PromQL queries (see set of predefined queries as a reference).
- (required) Job name identifies the job resource to parse and return the query result. Since the prometheus extension uses an exlusive namespace, any name can be chosen as long as it unique across `prometheus.Query` calls within the current test context.
- (optional) Container and controller parameters can be used as placeholders in PromQL queries. Note that some of the predefined queries only work if those values are provided.

## Service Provider Testing Flow

1. Add the `prometheus.Prometheus` to `OpenMCPSetup.Extensions`.
2. Use `prometheus.Query` to generate test results and reports.

## Domain Service Testing Flow

1. Add the `prometheus.Prometheus` to `OpenMCPSetup.Extensions`.
2. Add the control plane or workload cluster that runs your domain service as an additional scrape target with `prometheus.AddClusters`. A cluster has to be ready when you call this function. Note that adding clusters restarts prometheus which results in losing all previously scraped data. This is only a concern if you mix service provider testing with domain service testing as part of the same test suite which is not recommended.
3. Use `prometheus.Query` to generate test results and reports.

## Reference Usage

The [observability e2e test](../../../../e2e/observability_test.go) verifies that the Prometheus extension (un)installs and runs without errors but does not assert the correctness of the metrics or data it produces. Real-world tests should set an appropriate query window and reasonable waiting times.

## Testing Against Existing Environments

Use the extension directly instead of bootstrapping an OpenControlPlane enviroment.

```go
promExt := prometheus.Prometheus{}
testenv = env.NewWithConfig(envconf.New().WithNamespace(openmcp.Namespace))
testenv.Setup(func(ctx context.Context, c *envconf.Config) (context.Context, error) {
	return promExt.Install(ctx, c)
}).Finish(func(ctx context.Context, c *envconf.Config) (context.Context, error) {
	return promExt.Uninstall(ctx, c)
})
os.Exit(testenv.Run(m))
```

The current kubeconfig context has to point to the platform cluster. To interact with additional [clusters](https://open-control-plane.io/developers/clusters-and-namespaces#cluster-overview), e.g. to create resources on the onboarding or a control plane cluster, use `clusterutils.ConfigByEnv` instead of the kind-based config utility functions.

It is recommended to verify that an environment meets an expected configuration by using `openMCPSetup.Verify`.

Note that using an existing prometheus deployment is currently not supported to avoid depending on config assumptions that can't be adjusted.
