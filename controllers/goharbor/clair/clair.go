package clair

import (
	"context"
	"time"

	"github.com/ovh/configstore"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/apis/goharbor.io/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/controllers/common"
	"github.com/goharbor/harbor-operator/pkg/controllers/config"
	"github.com/goharbor/harbor-operator/pkg/event-filter/class"
)

const (
	DefaultRequeueWait    = 2 * time.Second
	ConfigImageKey        = "docker-image"
	DefaultImage          = "goharbor/clair-photon:v2.0.0"
	ConfigAdapterImageKey = "adapter-docker-image"
	DefaultAdapterImage   = "goharbor/clair-adapter-photon:v2.0.0"
)

// Reconciler reconciles a Clair object.
type Reconciler struct {
	*common.Controller
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := r.InitResources()
	if err != nil {
		return errors.Wrap(err, "cannot initialize resources")
	}

	err = r.Controller.SetupWithManager(mgr)
	if err != nil {
		return errors.Wrap(err, "cannot setup common controller")
	}

	className, err := r.ConfigStore.GetItemValue(config.HarborClassKey)
	if err != nil {
		return errors.Wrap(err, "cannot get harbor class")
	}

	concurrentReconcile, err := r.ConfigStore.GetItemValueInt(config.ReconciliationKey)
	if err != nil {
		return errors.Wrap(err, "cannot get concurrent reconcile")
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(&class.Filter{
			ClassName: className,
		}).
		For(&goharborv1alpha2.Clair{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: int(concurrentReconcile),
		}).
		Complete(r)
}

// +kubebuilder:rbac:groups=containerregistry.ovhcloud.com,resources=clairs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=containerregistry.ovhcloud.com,resources=clairs/status,verbs=get;update;patch

func New(ctx context.Context, name, version string, configStore *configstore.Store) (*Reconciler, error) {
	configStore.Env(name)

	r := &Reconciler{}

	r.Controller = common.NewController(name, version, r, configStore)

	return r, nil
}