/*
Copyright 2024 Wayne Wang<net_use@bzhy.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	
	monitoringv1beta1 "kubesysadm/api/v1beta1"
)

// CmMonitorReconciler reconciles a CmMonitor object
type CmMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.sysadm.cn,resources=cmmonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.sysadm.cn,resources=cmmonitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.sysadm.cn,resources=cmmonitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CmMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *CmMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("req info: ", "NS:", req.Namespace, "Name: ", req.Name, "toString:", req.String())

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CmMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1beta1.CmMonitor{}).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.getChangedMap),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.getChangedSecret),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

func (r *CmMonitorReconciler) getChangedMap(ctx context.Context, configMap client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)
	logger.Info("changed configMap: ", "NS:", configMap.GetNamespace(), "name: ", configMap.GetName(),
		"version:", configMap.GetResourceVersion())
	ns := configMap.GetNamespace()
	cmName := configMap.GetName()

	var cmMonitorList monitoringv1beta1.CmMonitorList
	e := r.List(ctx, &cmMonitorList)
	if e != nil {
		logger.Error(e, "get configMap monitored error")
		return nil
	}

	for _, item := range cmMonitorList.Items {
		logger.Info("monitor configmap item:", "namespace:", item.Spec.NameSpace, "name:", item.Spec.Name,
			"kind:", item.Spec.Kind)
		kind := strings.TrimSpace(strings.ToLower(item.Spec.Kind))
		if kind == "cm" || kind == "configmap" {
			if ns == item.Spec.NameSpace && cmName == item.Spec.Name {
				lastVersion := item.Status.LastVersion
				item.Status.LastVersion = configMap.GetResourceVersion()
				e := r.Client.Status().Update(ctx, &item)
				if e != nil {
					logger.Error(e, "update cmMonitor status error")
				}

				logger.Info("last version:", "version: ", lastVersion)

				if lastVersion == "" || lastVersion >= configMap.GetResourceVersion() {
					continue
				}
				r.restartWorkloadWithConfigMapChanged(ctx, ns, cmName)

			}
		}
	}
	return nil
}

func (r *CmMonitorReconciler) getChangedSecret(ctx context.Context, secret client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)
	logger.Info("changed secret: ", "NS:", secret.GetNamespace(), "name: ", secret.GetName(),
		"version:", secret.GetResourceVersion())

	return nil
}

func (r *CmMonitorReconciler) restartWorkloadWithConfigMapChanged(ctx context.Context, ns, cm string) {
	logger := log.FromContext(ctx)
	podList := corev1.PodList{}
	listOpts := &client.ListOptions{Namespace: ns}
	e := r.Client.List(ctx, &podList, listOpts)
	if e != nil {
		logger.Error(e, "get pod list error.")
		return
	}

	for _, item := range podList.Items {
		vols := item.Spec.Volumes
		for _, vol := range vols {
			if vol.ConfigMap != nil {
				cmName := vol.ConfigMap.LocalObjectReference.Name
				if strings.TrimSpace(strings.ToLower(cmName)) == strings.TrimSpace(strings.ToLower(cm)) {
					graceSecond := int64(5)
					deleteOpts := &client.DeleteOptions{GracePeriodSeconds: &graceSecond}
					e := r.Delete(ctx, &item, deleteOpts)
					if e != nil {
						logger.Error(e, "restart pod error")
						return
					}
				}
			}

		}
	}

}
