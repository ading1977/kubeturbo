package master

import (
	"github.com/turbonomic/kubeturbo/pkg/cluster"
	"github.com/turbonomic/kubeturbo/pkg/discovery/monitoring/types"

	"k8s.io/client-go/kubernetes"
)

type ClusterMonitorConfig struct {
	clusterInfoScraper *cluster.ClusterScraper
	targetIdentifier   string
}

func NewClusterMonitorConfig(kclient *kubernetes.Clientset, targetIdentifier string) *ClusterMonitorConfig {
	k8sClusterScraper := cluster.NewClusterScraper(kclient)
	return &ClusterMonitorConfig{
		clusterInfoScraper: k8sClusterScraper,
		targetIdentifier:   targetIdentifier,
	}
}

// Implement MonitoringWorkerConfig interface.
func (c ClusterMonitorConfig) GetMonitorType() types.MonitorType {
	return types.StateMonitor
}

// Implement MonitoringWorkerConfig interface.
func (c ClusterMonitorConfig) GetMonitoringSource() types.MonitoringSource {
	return types.ClusterSource
}
