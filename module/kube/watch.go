package kube

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

const (
	RUNNING = "running"
	ERROR   = "error"
	DELETED = "deleted"
	PENDING = "pending"
)

func WatchPod(podName string, podNamespace string, c chan string) {

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{"app": "nginx"}.AsSelector()}

	w, err := CLIENT.Pods(podNamespace).Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
	}

	for {
		event, ok := <-w.ResultChan()

		if !ok {
			// Resource was deleted, and chanle closed, so return to main programme
			return
		}
		switch event.Type {
		case "DELETED":
			c <- DELETED
			w.Stop()
		case "ERROR":
			c <- ERROR
			w.Stop()
		default:
			if event.Object.(*api.Pod).Status.Phase == "running" {
				c <- RUNNING
			} else {
				c <- PENDING
			}
		}
	}
}
