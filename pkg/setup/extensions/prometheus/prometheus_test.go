package prometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_applyDefaults(t *testing.T) {
	tests := []struct {
		name  string
		input Prometheus
		want  Prometheus
	}{
		{
			name:  "all defaults applied on zero value",
			input: Prometheus{},
			want: Prometheus{
				Namespace:             defaultNamespace,
				ClusterResourceSuffix: defaultSuffix,
				Image:                 defaultImage,
			},
		},
		{
			name: "existing values not overwritten",
			input: Prometheus{
				Namespace:             "custom-ns",
				ClusterResourceSuffix: "custom-suffix",
				Image:                 "custom-image:latest",
			},
			want: Prometheus{
				Namespace:             "custom-ns",
				ClusterResourceSuffix: "custom-suffix",
				Image:                 "custom-image:latest",
			},
		},
		{
			name: "partial: only namespace set",
			input: Prometheus{
				Namespace: "my-ns",
			},
			want: Prometheus{
				Namespace:             "my-ns",
				ClusterResourceSuffix: "openmcp-testing",
				Image:                 defaultImage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.applyDefaults()
			assert.Equal(t, tt.want, tt.input)
		})
	}
}
