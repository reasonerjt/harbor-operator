package portal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/apis/goharbor.io/v1alpha2"
)

const (
	PublicPort = 80
)

func (r *Reconciler) GetService(ctx context.Context, portal *goharborv1alpha2.Portal) (*corev1.Service, error) {
	name := r.NormalizeName(ctx, portal.GetName())
	namespace := portal.GetNamespace()

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       PublicPort,
					TargetPort: intstr.FromInt(port),
				},
			},
			Selector: map[string]string{
				"portal.goharbor.io/name":      name,
				"portal.goharbor.io/namespace": namespace,
			},
		},
	}, nil
}