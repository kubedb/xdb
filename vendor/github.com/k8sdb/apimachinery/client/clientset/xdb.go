package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	"k8s.io/kubernetes/pkg/api"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/watch"
)

type XdbNamespacer interface {
	Xdbes(namespace string) XdbInterface
}

type XdbInterface interface {
	List(opts api.ListOptions) (*aci.XdbList, error)
	Get(name string) (*aci.Xdb, error)
	Create(postgres *aci.Xdb) (*aci.Xdb, error)
	Update(postgres *aci.Xdb) (*aci.Xdb, error)
	Delete(name string) error
	Watch(opts api.ListOptions) (watch.Interface, error)
	UpdateStatus(postgres *aci.Xdb) (*aci.Xdb, error)
}

type XdbImpl struct {
	r  rest.Interface
	ns string
}

func newXdb(c *ExtensionsClient, namespace string) *XdbImpl {
	return &XdbImpl{c.restClient, namespace}
}

func (c *XdbImpl) List(opts api.ListOptions) (result *aci.XdbList, err error) {
	result = &aci.XdbList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *XdbImpl) Get(name string) (result *aci.Xdb, err error) {
	result = &aci.Xdb{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *XdbImpl) Create(postgres *aci.Xdb) (result *aci.Xdb, err error) {
	result = &aci.Xdb{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		Body(postgres).
		Do().
		Into(result)
	return
}

func (c *XdbImpl) Update(postgres *aci.Xdb) (result *aci.Xdb, err error) {
	result = &aci.Xdb{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		Name(postgres.Name).
		Body(postgres).
		Do().
		Into(result)
	return
}

func (c *XdbImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		Name(name).
		Do().
		Error()
}

func (c *XdbImpl) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *XdbImpl) UpdateStatus(postgres *aci.Xdb) (result *aci.Xdb, err error) {
	result = &aci.Xdb{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeXdb).
		Name(postgres.Name).
		SubResource("status").
		Body(postgres).
		Do().
		Into(result)
	return
}
