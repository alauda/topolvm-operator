/*
Copyright 2021.

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
// Code generated by lister-gen. DO NOT EDIT.

package v2

import (
	v2 "github.com/alauda/topolvm-operator/apis/topolvm/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TopolvmClusterLister helps list TopolvmClusters.
// All objects returned here must be treated as read-only.
type TopolvmClusterLister interface {
	// List lists all TopolvmClusters in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.TopolvmCluster, err error)
	// TopolvmClusters returns an object that can list and get TopolvmClusters.
	TopolvmClusters(namespace string) TopolvmClusterNamespaceLister
	TopolvmClusterListerExpansion
}

// topolvmClusterLister implements the TopolvmClusterLister interface.
type topolvmClusterLister struct {
	indexer cache.Indexer
}

// NewTopolvmClusterLister returns a new TopolvmClusterLister.
func NewTopolvmClusterLister(indexer cache.Indexer) TopolvmClusterLister {
	return &topolvmClusterLister{indexer: indexer}
}

// List lists all TopolvmClusters in the indexer.
func (s *topolvmClusterLister) List(selector labels.Selector) (ret []*v2.TopolvmCluster, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v2.TopolvmCluster))
	})
	return ret, err
}

// TopolvmClusters returns an object that can list and get TopolvmClusters.
func (s *topolvmClusterLister) TopolvmClusters(namespace string) TopolvmClusterNamespaceLister {
	return topolvmClusterNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TopolvmClusterNamespaceLister helps list and get TopolvmClusters.
// All objects returned here must be treated as read-only.
type TopolvmClusterNamespaceLister interface {
	// List lists all TopolvmClusters in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.TopolvmCluster, err error)
	// Get retrieves the TopolvmCluster from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v2.TopolvmCluster, error)
	TopolvmClusterNamespaceListerExpansion
}

// topolvmClusterNamespaceLister implements the TopolvmClusterNamespaceLister
// interface.
type topolvmClusterNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TopolvmClusters in the indexer for a given namespace.
func (s topolvmClusterNamespaceLister) List(selector labels.Selector) (ret []*v2.TopolvmCluster, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v2.TopolvmCluster))
	})
	return ret, err
}

// Get retrieves the TopolvmCluster from the indexer for a given namespace and name.
func (s topolvmClusterNamespaceLister) Get(name string) (*v2.TopolvmCluster, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v2.Resource("topolvmcluster"), name)
	}
	return obj.(*v2.TopolvmCluster), nil
}