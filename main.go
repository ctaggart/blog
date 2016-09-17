package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/rest"
)

func main() {
	// os.Setenv("HTTPS_PROXY", "http://127.0.0.1:8080") // enable proxy for debugging

	project := "ctaggartcom"
	zone := "us-central1-f"
	cluster := "demo3"

	gce, err := newContainerClient()
	check(err)

	err = createCluster(gce, project, zone, cluster)
	check(err)

	clstr, err := waitForClusterProvisioning(gce, project, zone, cluster)
	check(err)

	kbrnts, err := newKubernetesClient(clstr)
	check(err)

	listNodes(kbrnts)
}

// check simply panics if there is an error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// newContainerClient creates a new client for the Google Container Engine API.
func newContainerClient() (*container.Service, error) {
	ctx := context.Background()
	o := []option.ClientOption{
		option.WithEndpoint("https://container.googleapis.com/"),
		// option.WithScopes(container.CloudPlatformScope),
	}
	httpClient, endpoint, err := transport.NewHTTPClient(ctx, o...)
	if err != nil {
		return nil, err
	}
	client, err := container.New(httpClient)
	if err != nil {
		return nil, err
	}
	client.BasePath = endpoint
	return client, nil
}

func createCluster(gce *container.Service, project string, zone string, cluster string) error {
	req := &container.CreateClusterRequest{}
	req.Cluster = &container.Cluster{Name: cluster, InitialNodeCount: 3}
	_, err := gce.Projects.Zones.Clusters.Create(project, zone, req).Do()
	return err
}

func getCluster(gce *container.Service, project string, zone string, cluster string) (*container.Cluster, error) {
	return gce.Projects.Zones.Clusters.Get(project, zone, cluster).Do()
}

func waitForClusterProvisioning(gce *container.Service, project string, zone string, cluster string) (*container.Cluster, error) {
	for {
		clstr, err := getCluster(gce, project, zone, cluster)
		if err != nil {
			return nil, err
		}
		switch clstr.Status {
		case "PROVISIONING":
			time.Sleep(5 * time.Second)
		case "RUNNING":
			return clstr, nil
		default:
			return nil, fmt.Errorf("invalid cluster status: %s", clstr.Status)
		}
	}
}

func newKubernetesClient(clstr *container.Cluster) (*kubernetes.Clientset, error) {
	cert, err := base64.StdEncoding.DecodeString(clstr.MasterAuth.ClientCertificate)
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(clstr.MasterAuth.ClientKey)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(clstr.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}
	config := &rest.Config{
		Host:            clstr.Endpoint,
		TLSClientConfig: rest.TLSClientConfig{CertData: cert, KeyData: key, CAData: ca},
		Username:        clstr.MasterAuth.Username,
		Password:        clstr.MasterAuth.Password,
		// Insecure:        true,
	}
	kbrnts, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kbrnts, nil
}

func listNodes(kbrnts *kubernetes.Clientset) {
	nodes, err := kbrnts.Core().Nodes().List(api.ListOptions{})
	check(err)
	fmt.Printf("There are %d nodes in the cluster:\n", len(nodes.Items))
	for i, node := range nodes.Items {
		fmt.Printf("%d %s\n", i, node.Name)
	}
}

func listPods(kbrnts *kubernetes.Clientset) {
	pods, err := kbrnts.Core().Pods("").List(api.ListOptions{})
	check(err)
	fmt.Printf("There are %d pods in the cluster:\n", len(pods.Items))
	for i, pod := range pods.Items {
		fmt.Printf("%d %s %s\n", i, pod.Namespace, pod.Name)
	}
}

func deleteCluster(gce *container.Service, project string, zone string, cluster string) error {
	_, err := gce.Projects.Zones.Clusters.Delete(project, zone, cluster).Do()
	return err
}
