package resources

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDefaultServiceMonitor() v1.ServiceMonitor {
	return v1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceMonitor",
			APIVersion: "monitoring.coreos.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "",
			Labels: map[string]string{},
		},
		Spec: v1.ServiceMonitorSpec{
			Endpoints: nil,
			Selector:  metav1.LabelSelector{},
		},
	}
}

// CreateServiceMonitor creates a ServiceMonitor definition
func CreateServiceMonitor(app v1alpha1.Application, portName string) v1.ServiceMonitor {
	monitor := generateDefaultServiceMonitor()

	monitor.Name = app.Metadata.Name
	monitor.Namespace = app.Metadata.Namespace
	monitor.Labels["app"] = app.Metadata.Name

	monitor.Spec.Selector.MatchLabels = map[string]string{
		"app": app.Metadata.Name,
	}

	monitor.Spec.Endpoints = []v1.Endpoint{
		{
			Port: portName,
			Path: app.Prometheus.Path,
		},
	}

	return monitor
}
