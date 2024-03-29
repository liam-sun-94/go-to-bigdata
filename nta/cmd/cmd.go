package cmd

import (
	"flag"
	"path/filepath"
	"strconv"
	// "fmt"
	// "log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	// "google.golang.org/grpc"

	// "k8s.io/kubernetes/pkg/util/rand"

	"github.com/tangfeixiong/go-to-bigdata/nta/pkg/hbase"
	"github.com/tangfeixiong/go-to-bigdata/nta/pkg/server"
	"github.com/tangfeixiong/go-to-bigdata/pkg/util/homedir"
)

func RootCommandFor(name string) *cobra.Command {
	var config server.Config
	// in, out, errout := os.Stdin, os.Stdout, os.Stderr
	cfg := &config.HBase

	root := &cobra.Command{
		Use:   name,
		Short: "Collector server for NTA with gRPC & ReST API",
		Long: `
        Collector server for NTA
        
        This is a ..., and acting as a client of Apache HBase.
        
        It is inspired by some github projects.
        `,
	}
	root.AddCommand(createCollectorCommand(&config))
	root.AddCommand(createAgentCommand(cfg))

	root.PersistentFlags().StringVar(&cfg.Kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file. it means running out of cluster if supplied")
	if home := homedir.HomeDir(); home != "" {
		root.PersistentFlags().Lookup("kubeconfig").NoOptDefVal = filepath.Join(home, ".kube", "config")
	}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	return root
}

func createCollectorCommand(config *server.Config) *cobra.Command {

	command := &cobra.Command{
		Use:   "collect",
		Short: "Serving with gRPC and a gRPC-Gateway",
		Run: func(cmd *cobra.Command, args []string) {
			// pflag.Parse()
			flag.Set("v", strconv.Itoa(config.LogLevel))
			flag.Parse()
			server.Start(config)
		},
	}

	command.Flags().StringVar(&config.SecureAddress, "grpc_addr", "0.0.0.0:10001", "IP:port format")
	command.Flags().StringVar(&config.InsecureAddress, "http_addr", "0.0.0.0:10002", "IP:port format. Serve HTTP, or No HTTP if empty")
	command.Flags().BoolVar(&config.SecureHTTP, "secure_http", false, "Currently not used, if both HTTP address and HTTPS flag not set, just gRPC noly")
	command.Flags().IntVar(&config.LogLevel, "log_level", 2, "for glog")
	// command.Flags().AddGoFlagSet(flag.CommandLine)

	return command
}

func createAgentCommand(config *hbase.Config) *cobra.Command {

	command := &cobra.Command{
		Use:   "agent",
		Short: "init configurtions of etc/hadoop/",
		Run: func(cmd *cobra.Command, args []string) {
			// pflag.Parse()
			flag.Set("v", strconv.Itoa(config.LogLevel))
			flag.Parse()

			server.StartAgent(config)
		},
	}

	command.PersistentFlags().StringVar(&config.Name, "name", "", "StatefulSet name, or lookup value via label <crd group>/go-to-kubernetes")
	command.PersistentFlags().StringVar(&config.ServiceName, "service", "", "Kubernetes Service object name")
	command.PersistentFlags().StringVar(&config.Namespace, "namespace", "", "Kubernetes namespace, or lookup value from env name POD_NAMESPACE, otherwise 'default'")
	command.PersistentFlags().StringVar(&config.BaseDomain, "domain", "cluster.local", "Domain of K8s DNS")
	command.PersistentFlags().StringVar(&config.CustomResourceName, "custom_resource", "", "custom resource name")
	command.PersistentFlags().StringVar(&config.ClusterID, "hdfs_cluster_id", "", "HDFS cluster name, auto-creating by default")
	command.PersistentFlags().StringVar(&config.NodeType, "hdfs_node_type", "namenode", "Or datanode")
	command.PersistentFlags().StringVar(&config.Dir, "hadoop_dir", "/hadoop-3.0.0", "Directory of etc")
	command.PersistentFlags().IntVar(&config.Port, "service_port", 9000, "Service port")
	return command
}
