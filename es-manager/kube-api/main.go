package kubeapi
 
import (
	"context"
	"fmt"
	"log"
    "os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	vault "github.com/hashicorp/vault/api"
    "k8s.io/client-go/restmapper"
    "k8s.io/cli-runtime/pkg/resource"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"

)

func createObject(kubeClientset kubernetes.Interface, restConfig rest.Config, obj runtime.Object) (runtime.Object, error) {
    // Create a REST mapper that tracks information about the available resources in the cluster.
    groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
    if err != nil {
        return nil, err
    }
    rm := restmapper.NewDiscoveryRESTMapper(groupResources)

    // Get some metadata needed to make the REST request.
    gvk := obj.GetObjectKind().GroupVersionKind()
    gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
    mapping, err := rm.RESTMapping(gk, gvk.Version)
    if err != nil {
        return nil, err
    }

    namespace, err := meta.NewAccessor().Namespace(obj)
    if err != nil {
        return nil, err
    }

    // Create a client specifically for creating the object.
    restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
    if err != nil {
        return nil, err
    }

    // Use the REST helper to create the object in the "default" namespace.
    restHelper := resource.NewHelper(restClient, mapping)
    return restHelper.Create(namespace, false, obj)
}

func newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
    restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
    restConfig.GroupVersion = &gv
    if len(gv.Group) == 0 {
        restConfig.APIPath = "/api"
    } else {
        restConfig.APIPath = "/apis"
    }

    return rest.RESTClientFor(&restConfig)
}
func parseK8sYaml(fileR []byte) []runtime.Object {

	acceptedK8sTypes := regexp.MustCompile(`(Role|ClusterRole|RoleBinding|ClusterRoleBinding|ServiceAccount)`)
	fileAsString := string(fileR[:])
	sepYamlfiles := strings.Split(fileAsString, "---")
	retVal := make([]runtime.Object, 0, len(sepYamlfiles))
	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, groupVersionKind, err := decode([]byte(f), nil, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			continue
		}

		if !acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
			log.Printf("The custom-roles configMap contained K8s object types which are not supported! Skipping object with type: %s", groupVersionKind.Kind)
		} else {
			retVal = append(retVal, obj)
		}

	}
	return retVal
}
func kubernetesApply() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
        manifestPath := "vault-es/hello-world/one-secret"

	    // Apply the ExternalSecret manifest
	    err = esapply(clientset, manifestPath)
	    if err != nil {
		    panic(err.Error())
	    }

		time.Sleep(10 * time.Second)
	}
    return nil
}

