package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var unsorted = Slice{
	{Kind: "Pod", Name: "foo"},
	{Kind: "Prometheus", Name: "bar"},
	{Kind: "ConfigMap", Name: "foo"},
	{Kind: "Prometheus", Name: "foo"},
	{Kind: "CustomResourceDefinition", Name: "prometheus"},
	{Kind: "Pod", Name: "bar"},
	{Kind: "Alertmanager", Name: "baz"},
	{Kind: "Service", Name: "baz"},
	{Kind: "Namespace", Name: "kube-system"},
}

func TestSlice_SortApplyOrder(t *testing.T) {
	expected := Slice{
		{Kind: "Namespace", Name: "kube-system"},
		{Kind: "ConfigMap", Name: "foo"},
		{Kind: "CustomResourceDefinition", Name: "prometheus"},
		{Kind: "Service", Name: "baz"},
		{Kind: "Pod", Name: "bar"},
		{Kind: "Pod", Name: "foo"},
		{Kind: "Alertmanager", Name: "baz"},
		{Kind: "Prometheus", Name: "bar"},
		{Kind: "Prometheus", Name: "foo"},
	}

	assert.Equal(t, expected, unsorted.Sort(ApplyOrder))
}

func TestSlice_SortDeleteOrder(t *testing.T) {
	expected := Slice{
		{Kind: "Alertmanager", Name: "baz"},
		{Kind: "Prometheus", Name: "bar"},
		{Kind: "Prometheus", Name: "foo"},
		{Kind: "Service", Name: "baz"},
		{Kind: "Pod", Name: "bar"},
		{Kind: "Pod", Name: "foo"},
		{Kind: "CustomResourceDefinition", Name: "prometheus"},
		{Kind: "ConfigMap", Name: "foo"},
		{Kind: "Namespace", Name: "kube-system"},
	}

	assert.Equal(t, expected, unsorted.Sort(DeleteOrder))
}
