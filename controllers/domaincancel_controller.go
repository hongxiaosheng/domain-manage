/*
Copyright 2021.

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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	domainmanagev1alpha1 "cmit.com/crd/domain-manage/api/v1alpha1"
)

// DomainCancelReconciler reconciles a DomainCancel object
type DomainCancelReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

//+kubebuilder:rbac:groups=domain-manage.cmit.com,resources=domaincancels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=domain-manage.cmit.com,resources=domaincancels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=domain-manage.cmit.com,resources=domaincancels/finalizers,verbs=update
//+kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DomainCancel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *DomainCancelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("domaincancel", req.NamespacedName)

	// your logic here
	reqLogger := r.Log.WithValues("dc", req.NamespacedName)

	// 创建istio的clientset
	istioClient, err := versionedclient.NewForConfig(r.Config)
	if err != nil {
		// Request object not found.
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Other error.
		return ctrl.Result{}, err
	}

	// Create DomainCancel instance
	atInstance := &domainmanagev1alpha1.DomainCancel{}

	// Try to get cloud native dc instance.
	err = r.Get(ctx, req.NamespacedName, atInstance)
	if err != nil {
		// Request object not found.
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Other error.
		return ctrl.Result{}, err
	}

	DelValue(atInstance.Spec.Domains, reqLogger, istioClient)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DomainCancelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&domainmanagev1alpha1.DomainCancel{}).
		Complete(r)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
