package controller

import (
	"context"
	"log"
	"time"

	"github.com/hdkshingala/deploymentcreator/pkg/apis/hardik.dev/v1alpha1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	dcClientset "github.com/hdkshingala/deploymentcreator/pkg/client/clientset/versioned"
	dcInformer "github.com/hdkshingala/deploymentcreator/pkg/client/informers/externalversions/hardik.dev/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	client       kubernetes.Interface
	dcClient     dcClientset.Interface
	deploySynced cache.InformerSynced
	wq           workqueue.RateLimitingInterface
	timeStarted  time.Time
}

func NewController(client kubernetes.Interface, klient dcClientset.Interface, dcInformer dcInformer.DeploymentCreatorInformer) *Controller {

	cont := &Controller{
		client:       client,
		dcClient:     klient,
		deploySynced: dcInformer.Informer().HasSynced,
		wq:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kluster"),
		timeStarted:  time.Now().UTC(),
	}

	dcInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    cont.handleAdd,
			UpdateFunc: cont.handleUpdate,
			DeleteFunc: cont.handleDelete,
		},
	)

	return cont
}

func (cont *Controller) handleAdd(new interface{}) {
	dc, ok := new.(*v1alpha1.DeploymentCreator)
	if !ok {
		log.Printf("Error converting object '%+v' to DeploymentCreator when Add was called.\n", new)
	}

	if cont.timeStarted.Sub(dc.CreationTimestamp.Time.UTC()) > 0*time.Second {
		return
	} else {
		cont.wq.Add(dc)
	}
}

func (cont *Controller) handleUpdate(old, new interface{}) {
	newDc, ok1 := new.(*v1alpha1.DeploymentCreator)
	oldDc, ok2 := old.(*v1alpha1.DeploymentCreator)
	if !ok1 || !ok2 {
		log.Printf("Error converting objects '%+v' or '%+v' to DeploymentCreator when update was called.\n", old, new)
	}

	if oldDc.Spec.Image == newDc.Spec.Image && oldDc.Spec.Replicas == newDc.Spec.Replicas {
		return
	} else {
		cont.wq.Add(newDc)
	}
}

func (cont *Controller) handleDelete(old interface{}) {
	dc, ok := old.(*v1alpha1.DeploymentCreator)
	if !ok {
		log.Printf("Error converting object '%+v' to DeploymentCreator when Delete was called.\n", old)
	}

	if cont.timeStarted.Sub(dc.CreationTimestamp.Time.UTC()) > 0*time.Second {
		return
	} else {
		labels := dc.Labels
		if labels == nil {
			dc.Labels = map[string]string{}
		}
		dc.Labels["hardik.dev/isDeleted"] = "deleted"
		cont.wq.Add(dc)
	}
}

func (cont *Controller) Run(ch chan struct{}) error {
	if bool := cache.WaitForCacheSync(ch, cont.deploySynced); !bool {
		log.Println("cache was not synced.")
	}

	go wait.Until(cont.worker, time.Second, ch)
	<-ch
	return nil
}

func (cont *Controller) worker() {
	for cont.process() {

	}
}

func (cont *Controller) process() bool {
	item, shutDown := cont.wq.Get()
	if shutDown {
		log.Println("Cache is closed")
		return false
	}
	defer cont.wq.Forget(item)

	dc := item.(*v1alpha1.DeploymentCreator)
	ctx := context.Background()
	if dc.Labels["hardik.dev/isDeleted"] == "deleted" {
		log.Printf("Deleting deployment having name '%s'.\n", dc.Name)
		err := cont.client.AppsV1().Deployments(dc.Namespace).Delete(ctx, dc.Name+"-deployment", metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete deployment having name '%s'. Error: '%s'.\n", dc.Name, err.Error())
			return false
		}
		return true
	}

	labels := map[string]string{
		"k8s.io/app": dc.Name,
		"controller": dc.Name,
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dc.Name + "-deployment",
			Namespace: dc.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &dc.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  dc.Name + "-container",
							Image: dc.Spec.Image,
						},
					},
				},
			},
		},
	}

	log.Printf("Creating deployment having name '%s'.\n", dc.Name)
	deploy, err := cont.client.AppsV1().Deployments(dc.Namespace).Create(ctx, deploy, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			deploy, _ := cont.client.AppsV1().Deployments(dc.Namespace).Get(ctx, dc.Name, metav1.GetOptions{})
			deploy = deploy.DeepCopy()
			deploy.Spec.Replicas = &dc.Spec.Replicas
			deploy.Spec.Template.Spec.Containers[0].Image = dc.Spec.Image
			log.Printf("Updating deployment having name '%s'.\n", dc.Name)
			deploy, err := cont.client.AppsV1().Deployments(dc.Namespace).Update(ctx, deploy, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("Error received which Updating deployment, %s", err.Error())
				return false
			}
		} else {
			log.Printf("Error received which creating deployment, %s", err.Error())
			return false
		}
	}

	return true
}
