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

## Service Provider Testing Example

1. Add the prometheus extension to your OpenMCPSetup:

```go
setup.OpenMCPOperatorSetup{
	Name:         "openmcp-operator",
	Image:        "ghcr.io/openmcp-project/images/openmcp-operator:v1.3.0",
	Extensions: []extensions.Extension{
		&prometheus.Prometheus{},
	},
}
```

2. Add a `prometheus.Query` in your test suite (note that in a real test you want to provide a time window that actually represents the start and end time where you e.g. created load on the controller under test):

```go
params := prometheus.QueryParameters{
	Container:  container,
	Controller: controller,
	Window: prometheus.ObservabilityWindow{
		Start: time.Now().Add(-15 * time.Minute),
		End:   time.Now(),
	},
	JobName: fmt.Sprintf("query-%s", container),
}
results, err := prometheus.Query(ctx, cfg, queries, params)
```

## Domain Service Testing Example

1. Add the prometheus extension to your OpenMCPSetup (same as before).
2. Add additional clusters to monitor with `prometheus.AddClusters` to scrape control planes and workload clusters that have been provisioned by your test suite:

```go
Setup(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
	p := prometheus.Prometheus{}
	if err := p.AddClusters(ctx, c, prometheus.TargetCluster{
		Name:  "test", // control plane name
		Label: "test-controlplane",
	}); err != nil {
		t.Errorf("add scrape job for test control plane failed: %v", err)
	}
	return ctx
}).
```

Note that adding additional clusters restarts prometheus which results in losing all previously scraped data. This is only a concern if you mix evaluation of service provider performance with domain service performance as part of the same test which is currently not recommended.

3. Add a `prometheus.Query` in your test suite (same as before).

## Reference Usage

Both scenarios are covered by the simplified [observability e2e test](../../../../e2e/observability_test.go) example.

## Testing Against Existing Environments

Use the extension directly instead of bootstrapping an OpenControlPlane enviroment.

```go
promExt := prometheus.Prometheus{}
testenv = env.NewWithConfig(envconf.New().WithNamespace(openmcp.Namespace))
testenv.Setup(func(ctx context.Context, c *envconf.Config) (context.Context, error) {
	if err := promExt.Install(ctx, c); err != nil {
		return nil, err
	}
	return ctx, nil
}).Finish(func(ctx context.Context, c *envconf.Config) (context.Context, error) {
	if err := promExt.Uninstall(ctx, c); err != nil {
		return nil, err
	}
	return ctx, nil
})
os.Exit(testenv.Run(m))
```

The current kubeconfig context has to point to the platform cluster. To interact with additional [clusters](https://open-control-plane.io/developers/clusters-and-namespaces#cluster-overview), e.g. to create resources on the onboarding or a control plane cluster, use `clusterutils.ConfigByEnv` instead of the kind-based config utility functions.

It is recommended to verify that an environment meets an expected configuration by using `openMCPSetup.Verify`.

Note that using an existing prometheus deployment is currently not supported to avoid depending on config assumptions that can't be adjusted.
