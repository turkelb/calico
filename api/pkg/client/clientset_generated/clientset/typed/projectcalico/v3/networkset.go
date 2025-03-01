// Copyright (c) 2021 Tigera, Inc. All rights reserved.

// Code generated by client-gen. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	v3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	scheme "github.com/projectcalico/api/pkg/client/clientset_generated/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// NetworkSetsGetter has a method to return a NetworkSetInterface.
// A group's client should implement this interface.
type NetworkSetsGetter interface {
	NetworkSets(namespace string) NetworkSetInterface
}

// NetworkSetInterface has methods to work with NetworkSet resources.
type NetworkSetInterface interface {
	Create(ctx context.Context, networkSet *v3.NetworkSet, opts v1.CreateOptions) (*v3.NetworkSet, error)
	Update(ctx context.Context, networkSet *v3.NetworkSet, opts v1.UpdateOptions) (*v3.NetworkSet, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v3.NetworkSet, error)
	List(ctx context.Context, opts v1.ListOptions) (*v3.NetworkSetList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.NetworkSet, err error)
	NetworkSetExpansion
}

// networkSets implements NetworkSetInterface
type networkSets struct {
	client rest.Interface
	ns     string
}

// newNetworkSets returns a NetworkSets
func newNetworkSets(c *ProjectcalicoV3Client, namespace string) *networkSets {
	return &networkSets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the networkSet, and returns the corresponding networkSet object, and an error if there is any.
func (c *networkSets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v3.NetworkSet, err error) {
	result = &v3.NetworkSet{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("networksets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NetworkSets that match those selectors.
func (c *networkSets) List(ctx context.Context, opts v1.ListOptions) (result *v3.NetworkSetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v3.NetworkSetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("networksets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested networkSets.
func (c *networkSets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("networksets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a networkSet and creates it.  Returns the server's representation of the networkSet, and an error, if there is any.
func (c *networkSets) Create(ctx context.Context, networkSet *v3.NetworkSet, opts v1.CreateOptions) (result *v3.NetworkSet, err error) {
	result = &v3.NetworkSet{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("networksets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(networkSet).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a networkSet and updates it. Returns the server's representation of the networkSet, and an error, if there is any.
func (c *networkSets) Update(ctx context.Context, networkSet *v3.NetworkSet, opts v1.UpdateOptions) (result *v3.NetworkSet, err error) {
	result = &v3.NetworkSet{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("networksets").
		Name(networkSet.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(networkSet).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the networkSet and deletes it. Returns an error if one occurs.
func (c *networkSets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("networksets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *networkSets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("networksets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched networkSet.
func (c *networkSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.NetworkSet, err error) {
	result = &v3.NetworkSet{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("networksets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
