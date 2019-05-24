package executor

import (
	"github.com/golang/glog"
	apiappsv1beta1 "k8s.io/api/apps/v1beta1"
	apicorev1 "k8s.io/api/core/v1"
	apiextv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedappsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	typedextv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

// k8sController defines a common interface for kubernetes controller actions
// Currently supported controllers include:
// - ReplicationController
// - ReplicaSet
// - Deployment
type k8sController interface {
	get(name string) (*k8sControllerSpec, error)
	update(spec *k8sControllerSpec) error
}

// k8sControllerSpec defines a set of objects that we want to update:
// - replicas: The replicas of a controller to update for horizontal scale
// - podSpec: The pod specification of a controller to update for resize
type k8sControllerSpec struct {
	replicas int32
	podSpec  apicorev1.PodSpec
	isResize bool
	updated  bool
}

// ReplicationController
type replicationController struct {
	k8sController
	client typedcorev1.ReplicationControllerInterface
	rc     *apicorev1.ReplicationController
}

func (rc *replicationController) get(name string) (*k8sControllerSpec, error) {
	var err error
	rc.rc, err = rc.client.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &k8sControllerSpec{
		replicas: *rc.rc.Spec.Replicas,
		podSpec:  rc.rc.Spec.Template.Spec,
	}, nil
}

func (rc *replicationController) update(spec *k8sControllerSpec) error {
	var err error
	if spec.isResize {
		glog.V(2).Infof("Scaling down replica of ReplicaController to 0")
		*rc.rc.Spec.Replicas = 0
		if rc.rc, err = rc.client.Update(rc.rc); err != nil {
			return err
		}
	}
	*rc.rc.Spec.Replicas = spec.replicas
	rc.rc.Spec.Template.Spec = spec.podSpec
	_, err = rc.client.Update(rc.rc)
	return err
}

func (rc *replicationController) String() string {
	return "ReplicationController"
}

// ReplicaSet
type replicaSet struct {
	k8sController
	client typedextv1beta1.ReplicaSetInterface
	rs     *apiextv1beta1.ReplicaSet
}

func (rs *replicaSet) get(name string) (*k8sControllerSpec, error) {
	var err error
	rs.rs, err = rs.client.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &k8sControllerSpec{
		replicas: *rs.rs.Spec.Replicas,
		podSpec:  rs.rs.Spec.Template.Spec,
	}, nil
}

func (rs *replicaSet) update(spec *k8sControllerSpec) error {
	var err error
	if spec.isResize {
		glog.V(2).Infof("Scaling down replicas of ReplicaSet to 0")
		*rs.rs.Spec.Replicas = 0
		if rs.rs, err = rs.client.Update(rs.rs); err != nil {
			return err
		}
	}
	*rs.rs.Spec.Replicas = spec.replicas
	rs.rs.Spec.Template.Spec = spec.podSpec
	_, err = rs.client.Update(rs.rs)
	return err
}

func (rs *replicaSet) String() string {
	return "ReplicaSet"
}

// Deployment
type deployment struct {
	k8sController
	client typedappsv1beta1.DeploymentInterface
	dep    *apiappsv1beta1.Deployment
}

func (dep *deployment) get(name string) (*k8sControllerSpec, error) {
	var err error
	dep.dep, err = dep.client.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &k8sControllerSpec{
		replicas: *dep.dep.Spec.Replicas,
		podSpec:  dep.dep.Spec.Template.Spec,
	}, nil
}

func (dep *deployment) update(spec *k8sControllerSpec) error {
	*dep.dep.Spec.Replicas = spec.replicas
	dep.dep.Spec.Template.Spec = spec.podSpec
	_, err := dep.client.Update(dep.dep)
	return err
}

func (dep *deployment) String() string {
	return "Deployment"
}
