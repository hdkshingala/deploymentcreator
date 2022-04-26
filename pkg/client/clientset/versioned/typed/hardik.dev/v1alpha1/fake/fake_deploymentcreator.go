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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/hdkshingala/deploymentcreator/pkg/apis/hardik.dev/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDeploymentCreators implements DeploymentCreatorInterface
type FakeDeploymentCreators struct {
	Fake *FakeHardikV1alpha1
	ns   string
}

var deploymentcreatorsResource = schema.GroupVersionResource{Group: "hardik.dev", Version: "v1alpha1", Resource: "deploymentcreators"}

var deploymentcreatorsKind = schema.GroupVersionKind{Group: "hardik.dev", Version: "v1alpha1", Kind: "DeploymentCreator"}

// Get takes name of the deploymentCreator, and returns the corresponding deploymentCreator object, and an error if there is any.
func (c *FakeDeploymentCreators) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.DeploymentCreator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(deploymentcreatorsResource, c.ns, name), &v1alpha1.DeploymentCreator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeploymentCreator), err
}

// List takes label and field selectors, and returns the list of DeploymentCreators that match those selectors.
func (c *FakeDeploymentCreators) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.DeploymentCreatorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(deploymentcreatorsResource, deploymentcreatorsKind, c.ns, opts), &v1alpha1.DeploymentCreatorList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.DeploymentCreatorList{ListMeta: obj.(*v1alpha1.DeploymentCreatorList).ListMeta}
	for _, item := range obj.(*v1alpha1.DeploymentCreatorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested deploymentCreators.
func (c *FakeDeploymentCreators) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(deploymentcreatorsResource, c.ns, opts))

}

// Create takes the representation of a deploymentCreator and creates it.  Returns the server's representation of the deploymentCreator, and an error, if there is any.
func (c *FakeDeploymentCreators) Create(ctx context.Context, deploymentCreator *v1alpha1.DeploymentCreator, opts v1.CreateOptions) (result *v1alpha1.DeploymentCreator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(deploymentcreatorsResource, c.ns, deploymentCreator), &v1alpha1.DeploymentCreator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeploymentCreator), err
}

// Update takes the representation of a deploymentCreator and updates it. Returns the server's representation of the deploymentCreator, and an error, if there is any.
func (c *FakeDeploymentCreators) Update(ctx context.Context, deploymentCreator *v1alpha1.DeploymentCreator, opts v1.UpdateOptions) (result *v1alpha1.DeploymentCreator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(deploymentcreatorsResource, c.ns, deploymentCreator), &v1alpha1.DeploymentCreator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeploymentCreator), err
}

// Delete takes name of the deploymentCreator and deletes it. Returns an error if one occurs.
func (c *FakeDeploymentCreators) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(deploymentcreatorsResource, c.ns, name, opts), &v1alpha1.DeploymentCreator{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDeploymentCreators) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(deploymentcreatorsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.DeploymentCreatorList{})
	return err
}

// Patch applies the patch and returns the patched deploymentCreator.
func (c *FakeDeploymentCreators) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.DeploymentCreator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(deploymentcreatorsResource, c.ns, name, pt, data, subresources...), &v1alpha1.DeploymentCreator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DeploymentCreator), err
}
