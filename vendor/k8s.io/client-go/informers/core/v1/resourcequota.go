/*
Copyright 2017 The Kubernetes Authors.

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

// This file was automatically generated by informer-gen

package v1

import (
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// ResourceQuotaInformer provides access to a shared informer and lister for
// ResourceQuotas.
type ResourceQuotaInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ResourceQuotaLister
}

type resourceQuotaInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

func newResourceQuotaInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	sharedIndexInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().ResourceQuotas(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().ResourceQuotas(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&core_v1.ResourceQuota{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	return sharedIndexInformer
}

func (f *resourceQuotaInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&core_v1.ResourceQuota{}, newResourceQuotaInformer)
}

func (f *resourceQuotaInformer) Lister() v1.ResourceQuotaLister {
	return v1.NewResourceQuotaLister(f.Informer().GetIndexer())
}