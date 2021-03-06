package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"github.com/endocode/kelefstis/check"
	"text/template"
	"io/ioutil"
)

// simple k8s client that lists all available pods
// it gets config from $HOME/.kube/config
func main()  {
	usage:=`kelefstis.

		Usage:
	kelefstis <check> [--kubeconfig <config>]
	kelefstis [ -h | --help ]

Options:
	-h --help             Show this screen.
		check                 Template with the checks to run
	--kubeconfig <config>
`
	arguments,_ := docopt.Parse(usage,nil,false,"kelefstis 0.1", false)

	if arguments["--help"].(bool) {
		fmt.Printf(usage)
		return
	}

	fmt.Println(arguments)
	checkfile:=(arguments["<check>"].(string))

	kubeconfig, ok := arguments["--kubeconfig"].(string)
	if ! ok {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	checktemplate, err := ioutil.ReadFile(checkfile)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	chk := check.Check{true, clientset.CoreV1()}
	listNodes(clientset)
	listPods(clientset)
	tmpl, err := template.New("test").Parse(string(checktemplate))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, &chk)
	if err != nil {
		panic(err.Error())
	}
}
func listPods(clientset *kubernetes.Clientset) (*apiv1.PodList, error) {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if len(pods.Items) > 0 {
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for _, pod := range pods.Items {
			fmt.Printf("  Pod %-36s -%48s\n", pod.Name, pod.Labels)
		}

	}
	return pods, err
}
func listNodes(clientset *kubernetes.Clientset) (*apiv1.NodeList, error) {
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if len(nodes.Items) > 0 {
		fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))
		for _, node := range nodes.Items {
			fmt.Printf("  Node %-36s\n", node.Name)
		}

	}
	return nodes, err
}
