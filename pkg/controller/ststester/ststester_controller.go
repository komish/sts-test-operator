package ststester

import (
	"context"

	ststestv1alpha1 "github.com/komish/sts-test-operator/pkg/apis/ststest/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	// corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ststester")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new StsTester Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStsTester{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ststester-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource StsTester
	err = c.Watch(&source.Kind{Type: &ststestv1alpha1.StsTester{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/* OG Scaffolding
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ststestv1alpha1.StsTester{},
	})
	if err != nil {
		return err
	}
	*/

	// Watch for changes to secondary resource Pods and requeue the owner StsTester

	// We'll deploy a statefulset so we'll watch statefulsets.
	if err = c.Watch(
		&source.Kind{Type: &appsv1.StatefulSet{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &ststestv1alpha1.StsTester{},
		}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileStsTester implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileStsTester{}

// ReconcileStsTester reconciles a StsTester object
type ReconcileStsTester struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a StsTester object and makes changes based on the state read
// and what is in the StsTester.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileStsTester) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling StsTester")

	// Fetch the StsTester instance
	instance := &ststestv1alpha1.StsTester{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// If we got here, then we got a matching instance in the request. Now check to see
	// if we have our (arbitrary) config map in the API. And that we can even see the configmap.
	cm, err := getConfigMapForCR(r.client, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// We want to inform the user how to find information on this here. We'll still exit
			// because this is required.
			reqLogger.Error(err, "An example error about the config map being missing goes here. https://example.com/more-info-on-this",
				"ConfigMap.Name", instance.Spec.ConfigMapName)
		}
		// We should requeue here
		return reconcile.Result{}, err
	}
	reqLogger.Info("Found expected ConfigMap", "ConfigMap.Name", cm.ObjectMeta.Name)

	// Define a new Pod object
	sts := newStatefulsetForCR(instance)

	// Set StsTester instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, sts, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sts.Name, Namespace: sts.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Statefulset", "Statefulset.Namespace", sts.Namespace, "Statefulset.Name", sts.Name)
		err = r.client.Create(context.TODO(), sts)
		if err != nil {
			return reconcile.Result{}, err
		}

		// sts created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Statefulset already exists", "Statefulset.Namespace", found.Namespace, "Statefulset.Name", found.Name)
	return reconcile.Result{}, nil
}
