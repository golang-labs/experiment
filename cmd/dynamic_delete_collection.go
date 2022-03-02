package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func init() {
	klog.InitFlags(nil)
	flag.Set("v", "9")
}

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error creating dynamic client: %v\n", err)
		os.Exit(1)
	}

	Group := "asm.alauda.io"
	Version := "v1alpha1"
	Resource := "gatewaydeploys"

	gvr := schema.GroupVersionResource{
		Group:    Group,
		Version:  Version,
		Resource: Resource,
	}
	us, err := client.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			fmt.Printf("err is not notfound\n")
		} else {
			fmt.Printf("err is notfound\n")
		}

		fmt.Printf("error get all gatewaydeploys: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d\n", len(us.Items))

	nsMap := make(map[string]interface{})
	for _, item := range us.Items {
		ns := item.GetNamespace()
		if _, ok := nsMap[ns]; !ok {
			nsMap[ns] = struct{}{}
		}
	}

	fmt.Printf("%v\n", nsMap)

	if client.Resource(gvr).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{}); err != nil {
		fmt.Printf("error delete gatewaydeploy: %v\n", err)
		os.Exit(1)
	}
}
