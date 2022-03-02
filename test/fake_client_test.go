package test

import (
	"context"
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

func Test_NewSimpleDynamicClientWithCustomListKinds(t *testing.T) {
	scheme := runtime.NewScheme()

	Group := "asm.alauda.io"
	Version := "v1alpha1"
	Kind := "toy"
	Resource := "toys"

	client := fake.NewSimpleDynamicClientWithCustomListKinds(scheme,
		map[schema.GroupVersionResource]string{
			{Group: Group, Version: Version, Resource: Resource}: Kind + "List",
		},
		newUnstructured(Group+"/"+Version, Kind, "ns-foo", "name-foo"),
		newUnstructured(Group+"/"+Version, Kind, "ns-foo", "name-bar"),
		newUnstructured(Group+"/"+Version, Kind, "ns-foo", "name-baz"),
	)

	listFirst, _ := client.Resource(schema.GroupVersionResource{Group: Group, Version: Version, Resource: Resource}).List(context.TODO(), metav1.ListOptions{})
	fmt.Println(len(listFirst.Items))

	_, err := client.Resource(schema.GroupVersionResource{Group: Group, Version: Version, Resource: Resource}).Namespace("ns-foo").Create(context.TODO(), newUnstructured(Group+"/"+Version, Kind, "ns-foo", "name-foo2"), metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	listSecond, _ := client.Resource(schema.GroupVersionResource{Group: Group, Version: Version, Resource: Resource}).List(context.TODO(), metav1.ListOptions{})
	fmt.Println(len(listSecond.Items))
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}
}
