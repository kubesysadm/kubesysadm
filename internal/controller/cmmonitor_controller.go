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
	"fmt"
	"github.com/go-logr/logr"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoringv1beta1 "sysadm.cn/kubesysadm/api/v1beta1"
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
	b := ctrl.NewControllerManagedBy(mgr)
	b.For(&monitoringv1beta1.CmMonitor{})
	b.Watches(&corev1.ConfigMap{}, handler.EnqueueRequestsFromMapFunc(r.getChangedMap),
		builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}))
	b.Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(r.getChangedSecret),
		builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}))

	return b.Complete(r)

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
	ns := secret.GetNamespace()
	secretName := secret.GetName()

	var cmMonitorList monitoringv1beta1.CmMonitorList
	e := r.List(ctx, &cmMonitorList)
	if e != nil {
		logger.Error(e, "get secret monitored error")
		return nil
	}

	for _, item := range cmMonitorList.Items {
		logger.Info("monitor secret item:", "namespace:", item.Spec.NameSpace, "name:", item.Spec.Name,
			"kind:", item.Spec.Kind)
		kind := strings.TrimSpace(strings.ToLower(item.Spec.Kind))
		if kind == "secret" {
			if ns == item.Spec.NameSpace && secretName == item.Spec.Name {
				lastVersion := item.Status.LastVersion
				item.Status.LastVersion = secret.GetResourceVersion()
				e := r.Client.Status().Update(ctx, &item)
				if e != nil {
					logger.Error(e, "update cmMonitor status error")
				}

				logger.Info("last version:", "version: ", lastVersion)

				if lastVersion == "" || lastVersion >= secret.GetResourceVersion() {
					continue
				}
				r.restartWorkloadWithSecretChanged(ctx, ns, secretName)

			}
		}
	}
	return nil
}

func (r *CmMonitorReconciler) restartWorkloadWithConfigMapChanged(ctx context.Context, ns, cm string) {
	logger := log.FromContext(ctx)
	listOpts := &client.ListOptions{Namespace: ns}
	r.restartDeploys(ctx, logger, cm, "cm", listOpts)
	r.restartStatefulSet(ctx, logger, cm, "cm", listOpts)
	r.restartDaemonSet(ctx, logger, cm, "cm", listOpts)
}

func (r *CmMonitorReconciler) restartWorkloadWithSecretChanged(ctx context.Context, ns, secretName string) {
	logger := log.FromContext(ctx)
	listOpts := &client.ListOptions{Namespace: ns}
	r.restartDeploys(ctx, logger, secretName, "secret", listOpts)
	r.restartStatefulSet(ctx, logger, secretName, "secret", listOpts)
	r.restartDaemonSet(ctx, logger, secretName, "secret", listOpts)
}

func (r *CmMonitorReconciler) restartDeploys(ctx context.Context, logger logr.Logger, cm, kind string, listOpts *client.ListOptions) {
	deployList := appsv1.DeploymentList{}
	e := r.Client.List(ctx, &deployList, listOpts)
	if e != nil {
		logger.Error(e, "get Deployment list error.")
		return
	}

	for _, item := range deployList.Items {
		if isResourceUsedByPod(item.Spec.Template, cm, kind) {
			patchData := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubesysadm.sysadm.cn/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
			patch := client.RawPatch(types.MergePatchType, []byte(patchData))
			e := r.Patch(ctx, &item, patch)
			if e != nil {
				logger.Error(e, "patch deployment error.", "namespace: ", item.Namespace, "name:", item.Name)
			}
		}

	}
}

func (r *CmMonitorReconciler) restartStatefulSet(ctx context.Context, logger logr.Logger, cm, kind string, listOpts *client.ListOptions) {
	stsList := appsv1.StatefulSetList{}
	e := r.Client.List(ctx, &stsList, listOpts)
	if e != nil {
		logger.Error(e, "get StatefulSet list error.")
		return
	}

	for _, item := range stsList.Items {
		if isResourceUsedByPod(item.Spec.Template, cm, kind) {
			patchData := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubesysadm.sysadm.cn/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
			patch := client.RawPatch(types.MergePatchType, []byte(patchData))
			e := r.Patch(ctx, &item, patch)
			if e != nil {
				logger.Error(e, "patch statefulset error.", "namespace: ", item.Namespace, "name:", item.Name)
			}
		}

	}
}

func (r *CmMonitorReconciler) restartDaemonSet(ctx context.Context, logger logr.Logger, cm, kind string, listOpts *client.ListOptions) {
	daemonSetList := appsv1.DaemonSetList{}
	e := r.Client.List(ctx, &daemonSetList, listOpts)
	if e != nil {
		logger.Error(e, "get daemonSet list error.")
		return
	}

	for _, item := range daemonSetList.Items {
		if isResourceUsedByPod(item.Spec.Template, cm, kind) {
			patchData := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubesysadm.sysadm.cn/restartedAt": "%s"}}}}}`, time.Now().Format("20060102150405"))
			patch := client.RawPatch(types.MergePatchType, []byte(patchData))
			e := r.Patch(ctx, &item, patch)
			if e != nil {
				logger.Error(e, "patch daemonset error.", "namespace: ", item.Namespace, "name:", item.Name)
			}
		}

	}
}

func isResourceUsedByPod(podSpec corev1.PodTemplateSpec, name, kind string) bool {
	switch kind {
	case "cm":
		return isCmUsedByPod(podSpec, name)
	case "secret":
		return isSecretUsedByPod(podSpec, name)
	}

	return false

}

func isCmUsedByPod(podSpec corev1.PodTemplateSpec, name string) bool {
	// checking the configMap whether be mounted by the pod
	vols := podSpec.Spec.Volumes
	for _, vol := range vols {
		if vol.ConfigMap != nil {
			cm := vol.ConfigMap.LocalObjectReference.Name
			if strings.Compare(cm, name) == 0 {
				return true
			}
		}
	}

	// checking the configMap used with ENV
	containers := podSpec.Spec.Containers
	for _, c := range containers {
		envs := c.Env
		for _, env := range envs {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil {
				configMapKeyRef := env.ValueFrom.ConfigMapKeyRef
				cm := configMapKeyRef.LocalObjectReference.Name
				if strings.Compare(cm, name) == 0 {
					return true
				}
			}
		}
	}

	// checking the configMap used with ENV in initContainers
	initContainers := podSpec.Spec.InitContainers
	for _, c := range initContainers {
		envs := c.Env
		for _, env := range envs {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil {
				configMapKeyRef := env.ValueFrom.ConfigMapKeyRef
				cm := configMapKeyRef.LocalObjectReference.Name
				if strings.Compare(cm, name) == 0 {
					return true
				}
			}
		}
	}

	return false
}

func isSecretUsedByPod(podSpec corev1.PodTemplateSpec, name string) bool {
	// checking the secret whether be mounted by the pod
	vols := podSpec.Spec.Volumes
	for _, vol := range vols {
		if vol.Secret != nil {
			secretName := vol.Secret.SecretName
			if strings.Compare(secretName, name) == 0 {
				return true
			}
		}
	}

	// checking the secret used with ENV
	containers := podSpec.Spec.Containers
	for _, c := range containers {
		envs := c.Env
		for _, env := range envs {
			if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
				secretKeyRef := env.ValueFrom.SecretKeyRef
				secretName := secretKeyRef.LocalObjectReference.Name
				if strings.Compare(secretName, name) == 0 {
					return true
				}
			}
		}
	}

	// checking the secret used with ENV in initContainers
	initContainers := podSpec.Spec.InitContainers
	for _, c := range initContainers {
		envs := c.Env
		for _, env := range envs {
			if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
				secretKeyRef := env.ValueFrom.SecretKeyRef
				secretName := secretKeyRef.LocalObjectReference.Name
				if strings.Compare(secretName, name) == 0 {
					return true
				}
			}
		}
	}

	return false
}
