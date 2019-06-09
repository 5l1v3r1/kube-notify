package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/masahiro331/kube-notify/pkg/slack"
	"golang.org/x/xerrors"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informers "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	listers "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Controller struct {
	clientset  kubernetes.Interface
	jobsLister listers.JobLister
	jobsSynced cache.InformerSynced
	workqueue  workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface, jobInformer informers.JobInformer) *Controller {
	utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
	controller := &Controller{
		clientset:  clientset,
		jobsLister: jobInformer.Lister(),
		jobsSynced: jobInformer.Informer().HasSynced,
		workqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Notifyj"),
	}
	jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, new interface{}) {
			newJob := new.(*batchv1.Job)
			oldJob := old.(*batchv1.Job)

			//TODO: Ignore Logic ########
			fmt.Printf("%+v\n", oldJob)
			fmt.Printf("%+v\n", newJob)
			//###########################

			controller.handleObject(new)
		},
	})
	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	if ok := cache.WaitForCacheSync(stopCh, c.jobsSynced); !ok {
		return xerrors.New("failed to wait for caches to sync")
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	c.enqueue(object)
}

func MetaNamespaceKeyFunc(obj interface{}) (string, error) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return "", xerrors.Errorf("failed to MetaNamespaceKeyFunc: %v", err)
	}
	return getType(obj)[1:] + "/" + key, nil
}

func getType(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func SplitMetaNamespaceKey(key string) (kind, namespace, name string) {
	parts := strings.Split(key, "/")

	return parts[0], parts[1], parts[2]
}

func (c *Controller) enqueue(obj interface{}) {
	var key string
	var err error
	if key, err = MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)

		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return xerrors.Errorf("failed to error syncing: %s", err)
		}

		c.workqueue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	s := slack.SlackWriter{}
	kind, namespace, name := SplitMetaNamespaceKey(key)

	//TODO: Post Slack pkg/slack/slack.go ########
	fmt.Print(s)
	fmt.Print(namespace)
	fmt.Print(name)
	fmt.Print(kind)
	//############################################

	switch kind {
	case "":
		return nil

	}
	return nil
}
