package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigRenderYml(t *testing.T) {
	cfg := &Config{
		OTLPHost:     "0.0.0.0",
		OTLPHTTPPort: 4318,
		OTLPGRPCPort: 4317,
		EnableZipkin: true,
		EnableProm:   true,
		PromTarget: []string{
			"localhost:9000",
			"other-host:9000",
		},
	}
	want := `yaml:
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
      grpc:
        endpoint: 0.0.0.0:4317
  zipkin:
    endpoint: 0.0.0.0:9411
  prometheus:
    config:
      scrape_configs:
        - job_name: 'prometheus'
          scrape_interval: 15s
          static_configs:
            - targets:
              - 'localhost:9000'
              - 'other-host:9000'
processors:
exporters:
  tui:
service:
  pipelines:
    traces:
      receivers: 
        - otlp
        - zipkin
      processors:
      exporters:
        - tui
    logs:
      receivers:
        - otlp
      processors:
      exporters:
        - tui
    metrics:
      receivers:
        - otlp
        - prometheus
      processors:
      exporters:
        - tui
`
	got, err := cfg.RenderYml()
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestConfigRenderYmlMinimum(t *testing.T) {
	cfg := &Config{
		OTLPHost:     "0.0.0.0",
		OTLPHTTPPort: 4318,
		OTLPGRPCPort: 4317,
		EnableZipkin: false,
		EnableProm:   false,
	}
	want := `yaml:
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
      grpc:
        endpoint: 0.0.0.0:4317
processors:
exporters:
  tui:
service:
  pipelines:
    traces:
      receivers: 
        - otlp
      processors:
      exporters:
        - tui
    logs:
      receivers:
        - otlp
      processors:
      exporters:
        - tui
    metrics:
      receivers:
        - otlp
      processors:
      exporters:
        - tui
`
	got, err := cfg.RenderYml()
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want error
	}{
		{
			name: "OK_Minimum",
			cfg:  &Config{},
			want: nil,
		},
		{
			name: "OK_Maximum",
			cfg: &Config{
				OTLPHost:     "0.0.0.0",
				OTLPHTTPPort: 4318,
				OTLPGRPCPort: 4317,
				EnableZipkin: true,
				EnableProm:   true,
				PromTarget: []string{
					"localhost:9000",
					"other-host:9000",
				},
			},
			want: nil,
		},
		{
			name: "NG_Prom",
			cfg: &Config{
				OTLPHost:     "0.0.0.0",
				OTLPHTTPPort: 4318,
				OTLPGRPCPort: 4317,
				EnableProm:   true,
			},
			want: errors.New("the target endpoints for the prometheus receiver (--prom-target) must be specified when prometheus receiver enabled"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cfg.Validate())
		})
	}
}
