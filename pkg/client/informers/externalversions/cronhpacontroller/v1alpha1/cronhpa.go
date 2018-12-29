/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	cronhpacontrollerv1alpha1 "github.com/hex108/cron-hpa-controller/pkg/apis/cronhpacontroller/v1alpha1"
	versioned "github.com/hex108/cron-hpa-controller/pkg/client/clientset/versioned"
	internalinterfaces "github.com/hex108/cron-hpa-controller/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/hex108/cron-hpa-controller/pkg/client/listers/cronhpacontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CronHPAInformer provides access to a shared informer and lister for
// CronHPAs.
type CronHPAInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.CronHPALister
}

type cronHPAInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCronHPAInformer constructs a new informer for CronHPA type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCronHPAInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCronHPAInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCronHPAInformer constructs a new informer for CronHPA type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCronHPAInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CronhpacontrollerV1alpha1().CronHPAs(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CronhpacontrollerV1alpha1().CronHPAs(namespace).Watch(options)
			},
		},
		&cronhpacontrollerv1alpha1.CronHPA{},
		resyncPeriod,
		indexers,
	)
}

func (f *cronHPAInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCronHPAInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cronHPAInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cronhpacontrollerv1alpha1.CronHPA{}, f.defaultInformer)
}

func (f *cronHPAInformer) Lister() v1alpha1.CronHPALister {
	return v1alpha1.NewCronHPALister(f.Informer().GetIndexer())
}
