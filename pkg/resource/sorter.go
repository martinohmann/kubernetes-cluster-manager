// Adapted from https://github.com/helm/helm/blob/master/pkg/tiller/kind_sorter.go

package resource

import "sort"

// ResourceOrder is a slice of strings that defines the ordering of resources.
type ResourceOrder []string

// ApplyOrder is the resource order for apply operations.
var ApplyOrder ResourceOrder = []string{
	"Namespace",
	"ResourceQuota",
	"LimitRange",
	"PodSecurityPolicy",
	"PodDisruptionBudget",
	"Secret",
	"ConfigMap",
	"StorageClass",
	"PersistentVolume",
	"PersistentVolumeClaim",
	"ServiceAccount",
	"CustomResourceDefinition",
	"ClusterRole",
	"ClusterRoleBinding",
	"Role",
	"RoleBinding",
	"Service",
	"DaemonSet",
	"Pod",
	"ReplicationController",
	"ReplicaSet",
	"Deployment",
	"StatefulSet",
	"Job",
	"CronJob",
	"Ingress",
	"APIService",
}

// DeleteOrder is the resource order for delete operations.
var DeleteOrder ResourceOrder = []string{
	"APIService",
	"Ingress",
	"Service",
	"CronJob",
	"Job",
	"StatefulSet",
	"Deployment",
	"ReplicaSet",
	"ReplicationController",
	"Pod",
	"DaemonSet",
	"RoleBinding",
	"Role",
	"ClusterRoleBinding",
	"ClusterRole",
	"CustomResourceDefinition",
	"ServiceAccount",
	"PersistentVolumeClaim",
	"PersistentVolume",
	"StorageClass",
	"ConfigMap",
	"Secret",
	"PodDisruptionBudget",
	"PodSecurityPolicy",
	"LimitRange",
	"ResourceQuota",
	"Namespace",
}

type resourceSorter struct {
	order     map[string]int
	resources []*Resource
	isDelete  bool
}

func newResourceSorter(resources []*Resource, order ResourceOrder) *resourceSorter {
	o := make(map[string]int)

	for k, v := range order {
		o[v] = k
	}

	return &resourceSorter{
		resources: resources,
		order:     o,
		// TODO(mohmann): this is a little hacky and needs to be improved. We
		// need a way to make the sorter aware of the fact that custom
		// resources should be deleted before their CustomResourceDefinition.
		isDelete: order[len(order)-1] == "Namespace",
	}
}

// Len implements Len from sort.Interface.
func (s *resourceSorter) Len() int {
	return len(s.resources)
}

// Swap implements Swap from sort.Interface.
func (s *resourceSorter) Swap(i, j int) {
	s.resources[i], s.resources[j] = s.resources[j], s.resources[i]
}

// Less implements Less from sort.Interface.
func (s *resourceSorter) Less(i, j int) bool {
	a, b := s.resources[i], s.resources[j]

	aPos, aok := s.order[a.Kind]
	bPos, bok := s.order[b.Kind]

	if !aok && !bok {
		if a.Kind == b.Kind {
			return a.Name < b.Name
		}

		return a.Kind < b.Kind
	}

	if !aok {
		return s.isDelete
	}

	if !bok {
		return !s.isDelete
	}

	if aPos == bPos {
		return a.Name < b.Name
	}

	return aPos < bPos
}

func sortResources(resources []*Resource, order ResourceOrder) []*Resource {
	s := newResourceSorter(resources, order)

	sort.Sort(s)

	return resources
}
