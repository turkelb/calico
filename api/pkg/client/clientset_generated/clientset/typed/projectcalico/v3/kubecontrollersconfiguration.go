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

// KubeControllersConfigurationsGetter has a method to return a KubeControllersConfigurationInterface.
// A group's client should implement this interface.
type KubeControllersConfigurationsGetter interface {
	KubeControllersConfigurations() KubeControllersConfigurationInterface
}

// KubeControllersConfigurationInterface has methods to work with KubeControllersConfiguration resources.
type KubeControllersConfigurationInterface interface {
	Create(ctx context.Context, kubeControllersConfiguration *v3.KubeControllersConfiguration, opts v1.CreateOptions) (*v3.KubeControllersConfiguration, error)
	Update(ctx context.Context, kubeControllersConfiguration *v3.KubeControllersConfiguration, opts v1.UpdateOptions) (*v3.KubeControllersConfiguration, error)
	UpdateStatus(ctx context.Context, kubeControllersConfiguration *v3.KubeControllersConfiguration, opts v1.UpdateOptions) (*v3.KubeControllersConfiguration, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v3.KubeControllersConfiguration, error)
	List(ctx context.Context, opts v1.ListOptions) (*v3.KubeControllersConfigurationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.KubeControllersConfiguration, err error)
	KubeControllersConfigurationExpansion
}

// kubeControllersConfigurations implements KubeControllersConfigurationInterface
type kubeControllersConfigurations struct {
	client rest.Interface
}

// newKubeControllersConfigurations returns a KubeControllersConfigurations
func newKubeControllersConfigurations(c *ProjectcalicoV3Client) *kubeControllersConfigurations {
	return &kubeControllersConfigurations{
		client: c.RESTClient(),
	}
}

// Get takes name of the kubeControllersConfiguration, and returns the corresponding kubeControllersConfiguration object, and an error if there is any.
func (c *kubeControllersConfigurations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v3.KubeControllersConfiguration, err error) {
	result = &v3.KubeControllersConfiguration{}
	err = c.client.Get().
		Resource("kubecontrollersconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KubeControllersConfigurations that match those selectors.
func (c *kubeControllersConfigurations) List(ctx context.Context, opts v1.ListOptions) (result *v3.KubeControllersConfigurationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v3.KubeControllersConfigurationList{}
	err = c.client.Get().
		Resource("kubecontrollersconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kubeControllersConfigurations.
func (c *kubeControllersConfigurations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("kubecontrollersconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a kubeControllersConfiguration and creates it.  Returns the server's representation of the kubeControllersConfiguration, and an error, if there is any.
func (c *kubeControllersConfigurations) Create(ctx context.Context, kubeControllersConfiguration *v3.KubeControllersConfiguration, opts v1.CreateOptions) (result *v3.KubeControllersConfiguration, err error) {
	result = &v3.KubeControllersConfiguration{}
	err = c.client.Post().
		Resource("kubecontrollersconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubeControllersConfiguration).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a kubeControllersConfiguration and updates it. Returns the server's representation of the kubeControllersConfiguration, and an error, if there is any.
func (c *kubeControllersConfigurations) Update(ctx context.Context, kubeControllersConfiguration *v3.KubeControllersConfiguration, opts v1.UpdateOptions) (result *v3.KubeControllersConfiguration, err error) {
	result = &v3.KubeControllersConfiguration{}
	err = c.client.Put().
		Resource("kubecontrollersconfigurations").
		Name(kubeControllersConfiguration.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubeControllersConfiguration).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *kubeControllersConfigurations) UpdateStatus(ctx context.Context, kubeControllersConfiguration *v3.KubeControllersConfiguration, opts v1.UpdateOptions) (result *v3.KubeControllersConfiguration, err error) {
	result = &v3.KubeControllersConfiguration{}
	err = c.client.Put().
		Resource("kubecontrollersconfigurations").
		Name(kubeControllersConfiguration.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubeControllersConfiguration).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the kubeControllersConfiguration and deletes it. Returns an error if one occurs.
func (c *kubeControllersConfigurations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("kubecontrollersconfigurations").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kubeControllersConfigurations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("kubecontrollersconfigurations").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched kubeControllersConfiguration.
func (c *kubeControllersConfigurations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.KubeControllersConfiguration, err error) {
	result = &v3.KubeControllersConfiguration{}
	err = c.client.Patch(pt).
		Resource("kubecontrollersconfigurations").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
