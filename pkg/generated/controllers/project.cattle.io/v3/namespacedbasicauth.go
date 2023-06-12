/*
Copyright 2023 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	v3 "github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// NamespacedBasicAuthController interface for managing NamespacedBasicAuth resources.
type NamespacedBasicAuthController interface {
	generic.ControllerMeta
	NamespacedBasicAuthClient

	// OnChange runs the given handler when the controller detects a resource was changed.
	OnChange(ctx context.Context, name string, sync NamespacedBasicAuthHandler)

	// OnRemove runs the given handler when the controller detects a resource was changed.
	OnRemove(ctx context.Context, name string, sync NamespacedBasicAuthHandler)

	// Enqueue adds the resource with the given name to the worker queue of the controller.
	Enqueue(namespace, name string)

	// EnqueueAfter runs Enqueue after the provided duration.
	EnqueueAfter(namespace, name string, duration time.Duration)

	// Cache returns a cache for the resource type T.
	Cache() NamespacedBasicAuthCache
}

// NamespacedBasicAuthClient interface for managing NamespacedBasicAuth resources in Kubernetes.
type NamespacedBasicAuthClient interface {
	// Create creates a new object and return the newly created Object or an error.
	Create(*v3.NamespacedBasicAuth) (*v3.NamespacedBasicAuth, error)

	// Update updates the object and return the newly updated Object or an error.
	Update(*v3.NamespacedBasicAuth) (*v3.NamespacedBasicAuth, error)

	// Delete deletes the Object in the given name.
	Delete(namespace, name string, options *metav1.DeleteOptions) error

	// Get will attempt to retrieve the resource with the specified name.
	Get(namespace, name string, options metav1.GetOptions) (*v3.NamespacedBasicAuth, error)

	// List will attempt to find multiple resources.
	List(namespace string, opts metav1.ListOptions) (*v3.NamespacedBasicAuthList, error)

	// Watch will start watching resources.
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)

	// Patch will patch the resource with the matching name.
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.NamespacedBasicAuth, err error)
}

// NamespacedBasicAuthCache interface for retrieving NamespacedBasicAuth resources in memory.
type NamespacedBasicAuthCache interface {
	// Get returns the resources with the specified name from the cache.
	Get(namespace, name string) (*v3.NamespacedBasicAuth, error)

	// List will attempt to find resources from the Cache.
	List(namespace string, selector labels.Selector) ([]*v3.NamespacedBasicAuth, error)

	// AddIndexer adds  a new Indexer to the cache with the provided name.
	// If you call this after you already have data in the store, the results are undefined.
	AddIndexer(indexName string, indexer NamespacedBasicAuthIndexer)

	// GetByIndex returns the stored objects whose set of indexed values
	// for the named index includes the given indexed value.
	GetByIndex(indexName, key string) ([]*v3.NamespacedBasicAuth, error)
}

// NamespacedBasicAuthHandler is function for performing any potential modifications to a NamespacedBasicAuth resource.
type NamespacedBasicAuthHandler func(string, *v3.NamespacedBasicAuth) (*v3.NamespacedBasicAuth, error)

// NamespacedBasicAuthIndexer computes a set of indexed values for the provided object.
type NamespacedBasicAuthIndexer func(obj *v3.NamespacedBasicAuth) ([]string, error)

// NamespacedBasicAuthGenericController wraps wrangler/pkg/generic.Controller so that the function definitions adhere to NamespacedBasicAuthController interface.
type NamespacedBasicAuthGenericController struct {
	generic.ControllerInterface[*v3.NamespacedBasicAuth, *v3.NamespacedBasicAuthList]
}

// OnChange runs the given resource handler when the controller detects a resource was changed.
func (c *NamespacedBasicAuthGenericController) OnChange(ctx context.Context, name string, sync NamespacedBasicAuthHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.NamespacedBasicAuth](sync))
}

// OnRemove runs the given object handler when the controller detects a resource was changed.
func (c *NamespacedBasicAuthGenericController) OnRemove(ctx context.Context, name string, sync NamespacedBasicAuthHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.NamespacedBasicAuth](sync))
}

// Cache returns a cache of resources in memory.
func (c *NamespacedBasicAuthGenericController) Cache() NamespacedBasicAuthCache {
	return &NamespacedBasicAuthGenericCache{
		c.ControllerInterface.Cache(),
	}
}

// NamespacedBasicAuthGenericCache wraps wrangler/pkg/generic.Cache so the function definitions adhere to NamespacedBasicAuthCache interface.
type NamespacedBasicAuthGenericCache struct {
	generic.CacheInterface[*v3.NamespacedBasicAuth]
}

// AddIndexer adds  a new Indexer to the cache with the provided name.
// If you call this after you already have data in the store, the results are undefined.
func (c NamespacedBasicAuthGenericCache) AddIndexer(indexName string, indexer NamespacedBasicAuthIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.NamespacedBasicAuth](indexer))
}
