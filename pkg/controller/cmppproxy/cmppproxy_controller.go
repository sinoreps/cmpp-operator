package cmppproxy

import (
	"context"
	"reflect"

	cmppv1alpha1 "github.com/sinoreps/cmpp-operator/pkg/apis/cmpp/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_cmppproxy")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CMPPProxy Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCMPPProxy{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cmppproxy-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CMPPProxy
	err = c.Watch(&source.Kind{Type: &cmppv1alpha1.CMPPProxy{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CMPPProxy
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cmppv1alpha1.CMPPProxy{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCMPPProxy{}

// ReconcileCMPPProxy reconciles a CMPPProxy object
type ReconcileCMPPProxy struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CMPPProxy object and makes changes based on the state read
// and what is in the CMPPProxy.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCMPPProxy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CMPPProxy")

	// Fetch the CMPPProxy instance
	instance := &cmppv1alpha1.CMPPProxy{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForCMPPProxy(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Define a new Service
		svc := r.serviceForCMPPProxy(instance)
		reqLogger.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return reconcile.Result{}, err
		}

		// Resources created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	pendingUpdates := false
	// Ensure the deployment size is the same as the spec
	size := instance.Spec.NumConnections
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		pendingUpdates = true
	}
	if found.Spec.Template.Spec.Containers[0].Image != instance.Spec.Image {
		found.Spec.Template.Spec.Containers[0].Image = instance.Spec.Image
		pendingUpdates = true
	}
	serverAddr := instance.Spec.ServerAddr
	account := instance.Spec.Account
	password := instance.Spec.Password
	enterpriseCode := instance.Spec.EnterpriseCode
	serviceCode := instance.Spec.ServiceCode
	envs := found.Spec.Template.Spec.Containers[0].Env
	pendingUpdates = setEnvVarIfNeeded(&envs, "CMPP_SERVER_ADDR", serverAddr) || pendingUpdates
	pendingUpdates = setEnvVarIfNeeded(&envs, "CMPP_ACCOUNT", account) || pendingUpdates
	pendingUpdates = setEnvVarIfNeeded(&envs, "CMPP_PASSWORD", password) || pendingUpdates
	pendingUpdates = setEnvVarIfNeeded(&envs, "CMPP_ENTERPRISE_CODE", enterpriseCode) || pendingUpdates
	pendingUpdates = setEnvVarIfNeeded(&envs, "CMPP_SERVICE_CODE", serviceCode) || pendingUpdates
	if pendingUpdates {
		found.Spec.Template.Spec.Containers[0].Env = envs
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the CMPPProxy status with the pod names
	// List the pods for this CMPPProxy's deployment
	podList := &corev1.PodList{}
	listOpts := client.InNamespace(instance.Namespace).MatchingLabels(labelsForCMPPProxy(instance.Name))

	if err = r.client.List(context.TODO(), listOpts, podList); err != nil {
		reqLogger.Error(err, "Failed to list pods", "CMPPProxy.Namespace", instance.Namespace, "CMPPProxy.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Pods) {
		instance.Status.Pods = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update CMPPProxy status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// deploymentForCMPPProxy returns an cmppproxy Deployment object
func (r *ReconcileCMPPProxy) deploymentForCMPPProxy(m *cmppv1alpha1.CMPPProxy) *appsv1.Deployment {
	ls := labelsForCMPPProxy(m.Name)
	replicas := m.Spec.NumConnections
	image := m.Spec.Image
	serverAddr := m.Spec.ServerAddr
	account := m.Spec.Account
	password := m.Spec.Password
	enterpriseCode := m.Spec.EnterpriseCode
	serviceCode := m.Spec.ServiceCode

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  "cmppproxy",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
						}},
						Env: []corev1.EnvVar{
							corev1.EnvVar{
								Name:  "CMPP_SERVER_ADDR",
								Value: serverAddr,
							},
							corev1.EnvVar{
								Name:  "CMPP_ACCOUNT",
								Value: account,
							},
							corev1.EnvVar{
								Name:  "CMPP_PASSWORD",
								Value: password,
							},
							corev1.EnvVar{
								Name:  "CMPP_ENTERPRISE_CODE",
								Value: enterpriseCode,
							},
							corev1.EnvVar{
								Name:  "CMPP_SERVICE_CODE",
								Value: serviceCode,
							}},
					}},
				},
			},
		},
	}
	// Set CMPPProxy instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// serviceForCMPPProxy returns an cmppproxy Service object
func (r *ReconcileCMPPProxy) serviceForCMPPProxy(m *cmppv1alpha1.CMPPProxy) *corev1.Service {
	ls := labelsForCMPPProxy(m.Name)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}
	// Set CMPPProxy instance as the owner and controller
	controllerutil.SetControllerReference(m, svc, r.scheme)
	return svc
}

// labelsForCMPPProxy returns the labels for selecting the resources
// belonging to the given CMPPProxy CR name.
func labelsForCMPPProxy(name string) map[string]string {
	return map[string]string{"app": "cmppproxy", "cmppproxy_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func getEnvVar(envs *[]corev1.EnvVar, key string) *corev1.EnvVar {
	for i, env := range *envs {
		if env.Name == key {
			return &(*envs)[i]
		}
	}
	return nil
}

func setEnvVarIfNeeded(envs *[]corev1.EnvVar, key string, value string) bool {
	existing := getEnvVar(envs, key)
	if existing != nil {
		if existing.Value == value {
			return false
		}
		(*existing).Value = value
	} else {
		env := corev1.EnvVar{
			Name:  key,
			Value: value,
		}
		*envs = append(*envs, env)
	}

	return true
}
