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

// DockerCredentialController interface for managing DockerCredential resources.
type DockerCredentialController interface {
	generic.ControllerMeta
	DockerCredentialClient

	// OnChange runs the given handler when the controller detects a resource was changed.
	OnChange(ctx context.Context, name string, sync DockerCredentialHandler)

	// OnRemove runs the given handler when the controller detects a resource was changed.
	OnRemove(ctx context.Context, name string, sync DockerCredentialHandler)

	// Enqueue adds the resource with the given name to the worker queue of the controller.
	Enqueue(namespace, name string)

	// EnqueueAfter runs Enqueue after the provided duration.
	EnqueueAfter(namespace, name string, duration time.Duration)

	// Cache returns a cache for the resource type T.
	Cache() DockerCredentialCache
}

// DockerCredentialClient interface for managing DockerCredential resources in Kubernetes.
type DockerCredentialClient interface {
	// Create creates a new object and return the newly created Object or an error.
	Create(*v3.DockerCredential) (*v3.DockerCredential, error)

	// Update updates the object and return the newly updated Object or an error.
	Update(*v3.DockerCredential) (*v3.DockerCredential, error)

	// Delete deletes the Object in the given name.
	Delete(namespace, name string, options *metav1.DeleteOptions) error

	// Get will attempt to retrieve the resource with the specified name.
	Get(namespace, name string, options metav1.GetOptions) (*v3.DockerCredential, error)

	// List will attempt to find multiple resources.
	List(namespace string, opts metav1.ListOptions) (*v3.DockerCredentialList, error)

	// Watch will start watching resources.
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)

	// Patch will patch the resource with the matching name.
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.DockerCredential, err error)
}

// DockerCredentialCache interface for retrieving DockerCredential resources in memory.
type DockerCredentialCache interface {
	// Get returns the resources with the specified name from the cache.
	Get(namespace, name string) (*v3.DockerCredential, error)

	// List will attempt to find resources from the Cache.
	List(namespace string, selector labels.Selector) ([]*v3.DockerCredential, error)

	// AddIndexer adds  a new Indexer to the cache with the provided name.
	// If you call this after you already have data in the store, the results are undefined.
	AddIndexer(indexName string, indexer DockerCredentialIndexer)

	// GetByIndex returns the stored objects whose set of indexed values
	// for the named index includes the given indexed value.
	GetByIndex(indexName, key string) ([]*v3.DockerCredential, error)
}

// DockerCredentialHandler is function for performing any potential modifications to a DockerCredential resource.
type DockerCredentialHandler func(string, *v3.DockerCredential) (*v3.DockerCredential, error)

// DockerCredentialIndexer computes a set of indexed values for the provided object.
type DockerCredentialIndexer func(obj *v3.DockerCredential) ([]string, error)

// DockerCredentialGenericController wraps wrangler/pkg/generic.Controller so that the function definitions adhere to DockerCredentialController interface.
type DockerCredentialGenericController struct {
	generic.ControllerInterface[*v3.DockerCredential, *v3.DockerCredentialList]
}

// OnChange runs the given resource handler when the controller detects a resource was changed.
func (c *DockerCredentialGenericController) OnChange(ctx context.Context, name string, sync DockerCredentialHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.DockerCredential](sync))
}

// OnRemove runs the given object handler when the controller detects a resource was changed.
func (c *DockerCredentialGenericController) OnRemove(ctx context.Context, name string, sync DockerCredentialHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.DockerCredential](sync))
}

// Cache returns a cache of resources in memory.
func (c *DockerCredentialGenericController) Cache() DockerCredentialCache {
	return &DockerCredentialGenericCache{
		c.ControllerInterface.Cache(),
	}
}

// DockerCredentialGenericCache wraps wrangler/pkg/generic.Cache so the function definitions adhere to DockerCredentialCache interface.
type DockerCredentialGenericCache struct {
	generic.CacheInterface[*v3.DockerCredential]
}

// AddIndexer adds  a new Indexer to the cache with the provided name.
// If you call this after you already have data in the store, the results are undefined.
func (c DockerCredentialGenericCache) AddIndexer(indexName string, indexer DockerCredentialIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.DockerCredential](indexer))
}
