/*
Copyright 2021 The Topolvm-Operator Authors. All rights reserved.

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

package operator

import (
	"flag"
	"fmt"
	"github.com/coreos/pkg/capnslog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	topolvmv1 "topolvm-operator/api/v1"
	"topolvm-operator/cmd/topolvm"
	"topolvm-operator/controllers"
	"topolvm-operator/pkg/cluster"
	"topolvm-operator/pkg/operator/csidriver"
	"topolvm-operator/pkg/operator/psp"
)

var OperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Check Disk and Create Volume group",
}

var (
	scheme = runtime.NewScheme()
	logger = capnslog.NewPackageLogger("topolvm/operator", "topolvm-cluster")
)

func init() {
	OperatorCmd.RunE = startOperator
	addScheme()
}

func addScheme() {

	_ = clientgoscheme.AddToScheme(scheme)

	_ = topolvmv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func startOperator(cmd *cobra.Command, args []string) error {

	cluster.SetLogLevel()
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "355331c5.cybozu.com",
	})

	if err != nil {
		logger.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := cluster.NewContext()
	ctx.Client = mgr.GetClient()

	cluster.NameSpace = os.Getenv(cluster.PodNameSpaceEnv)
	if cluster.NameSpace == "" {
		logger.Errorf("unable get env %s ", cluster.PodNameSpaceEnv)
		return fmt.Errorf("get env:%s failed ", cluster.PodNameSpaceEnv)
	}

	err = checkAndCreatePsp(ctx.Clientset)
	if err != nil {
		logger.Errorf("checkAndCreatePsp failed err %v", err)
		return err
	}

	err = csidriver.CheckTopolvmCsiDriverExisting(ctx.Clientset)
	if err != nil {
		logger.Errorf("CheckTopolvmCsiDriverExisting failed err %v", err)
		return err
	}

	operatorImage := topolvm.GetOperatorImage(ctx.Clientset, "")
	c := controllers.NewTopolvmClusterReconciler(mgr.GetScheme(), ctx, operatorImage)
	if err := c.SetupWithManager(mgr); err != nil {
		logger.Error(err, "unable to create controller", "controller", "TopolvmCluster")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	logger.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "problem running manager")
		os.Exit(1)
	}

	return nil
}

func checkAndCreatePsp(clientset kubernetes.Interface) error {

	existing, err := psp.CheckPspExisting(clientset, cluster.TopolvmNodePsp)
	if err != nil {
		return errors.Wrapf(err, "check psp %s failed", cluster.TopolvmNodePsp)
	}

	if !existing {
		err = psp.CreateTopolvmNodePsp(clientset)
		if err != nil {
			return errors.Wrapf(err, "create psp %s failed", cluster.TopolvmNodePsp)
		}
	} else {
		logger.Infof("psp %s existing", cluster.TopolvmNodePsp)
	}

	existing, err = psp.CheckPspExisting(clientset, cluster.TopolvmPrepareVgPsp)
	if err != nil {
		return errors.Wrapf(err, "check psp %s failed", cluster.TopolvmPrepareVgPsp)
	}

	if !existing {
		err = psp.CreateTopolvmPrepareVgPsp(clientset)
		if err != nil {
			return errors.Wrapf(err, "create psp %s failed", cluster.TopolvmPrepareVgPsp)
		}
	} else {
		logger.Infof("psp %s existing", cluster.TopolvmPrepareVgPsp)
	}

	return nil
}