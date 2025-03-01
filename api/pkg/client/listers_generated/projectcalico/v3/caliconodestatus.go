// Copyright (c) 2021 Tigera, Inc. All rights reserved.

// Code generated by lister-gen. DO NOT EDIT.

package v3

import (
	v3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CalicoNodeStatusLister helps list CalicoNodeStatuses.
// All objects returned here must be treated as read-only.
type CalicoNodeStatusLister interface {
	// List lists all CalicoNodeStatuses in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v3.CalicoNodeStatus, err error)
	// Get retrieves the CalicoNodeStatus from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v3.CalicoNodeStatus, error)
	CalicoNodeStatusListerExpansion
}

// calicoNodeStatusLister implements the CalicoNodeStatusLister interface.
type calicoNodeStatusLister struct {
	indexer cache.Indexer
}

// NewCalicoNodeStatusLister returns a new CalicoNodeStatusLister.
func NewCalicoNodeStatusLister(indexer cache.Indexer) CalicoNodeStatusLister {
	return &calicoNodeStatusLister{indexer: indexer}
}

// List lists all CalicoNodeStatuses in the indexer.
func (s *calicoNodeStatusLister) List(selector labels.Selector) (ret []*v3.CalicoNodeStatus, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.CalicoNodeStatus))
	})
	return ret, err
}

// Get retrieves the CalicoNodeStatus from the index for a given name.
func (s *calicoNodeStatusLister) Get(name string) (*v3.CalicoNodeStatus, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v3.Resource("caliconodestatus"), name)
	}
	return obj.(*v3.CalicoNodeStatus), nil
}
