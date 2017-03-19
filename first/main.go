package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

func main() {
	//

	//config, err := clientcmd.BuildConfigFromFlags("", "$HOME/.kube/config")
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// create a Pod
	pod, err := c.Pods(v1.Namespace.Default).Create(&v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "spacemonkey",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "Spacemonkey",
					Image: "quay.io/yazpik/spacemonkey:latest",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	})

	// expose the pod to a service
	pod.SetLabels(map[string]string{
		"pod-group": "pod-group",
	})
	//include Labels
	pod, err = c.Pods(v1.NamespaceDefault).Update(pod)

	//create a service
	svc, err := c.Services(v1.NamespaceDefault).Create(&v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name: "spacemonkey-svc",
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Selector: pod.Labels,
			Ports: []v1.ServicePort{
				{
					Port: 8888,
				},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// print the service is exposed
	fmt.Println(svc.Specs.Ports[0].NodePort)

}
