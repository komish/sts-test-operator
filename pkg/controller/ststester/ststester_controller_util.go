package ststester

import (
	"context"
	ststestv1alpha1 "github.com/komish/sts-test-operator/pkg/apis/ststest/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"             // use this to get NamespacedName type when getting things from the API
	"sigs.k8s.io/controller-runtime/pkg/client" // use this to get things from the API!
)

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *ststestv1alpha1.StsTester) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func newStatefulsetForCR(cr *ststestv1alpha1.StsTester) *appsv1.StatefulSet {
	// pre-load this so we can use it multiple places
	stsName := cr.Name + "-statefulset"

	// Here are the labels that we'll apply.
	// TODO Extend to also take labels from the cr meta
	stsLabels := map[string]string{
		"sts.komish.net/owning-cr": cr.Name,
	}

	// These need to be used as selectors
	podLabels := map[string]string{
		"sts.komish.net/owning-sts": stsName,
	}

	sts := &appsv1.StatefulSet{
		// The statefulset metadata
		ObjectMeta: metav1.ObjectMeta{
			Name:      stsName,
			Namespace: cr.Namespace,
			Labels:    stsLabels,
		},
		// The statefulset's Spec
		Spec: appsv1.StatefulSetSpec{
			// We pre-define the labels for the pods, so we match them here
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			// This comes from the CR, our CR spec should have a replicas field
			Replicas: &cr.Spec.Replicas,
			// Template here refers to a pod template
			Template: corev1.PodTemplateSpec{
				// Objectmeta again, this time for the pod
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name + "-pod",
					Namespace: cr.Namespace,
					// Make sure we use the same labels as we defined for Selector
					Labels: podLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   "busybox",
							Command: []string{"sleep", "9001"},
						},
					},
				},
			},
		},
	}

	// We already have a pointer, so we don't need to return
	// the memory address here.
	return sts
}

func getConfigMapForCR(client client.Client, cr *ststestv1alpha1.StsTester) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(
		context.TODO(),
		types.NamespacedName{
			Namespace: cr.ObjectMeta.Namespace,
			Name:      cr.Spec.ConfigMapName},
		cm,
	)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
