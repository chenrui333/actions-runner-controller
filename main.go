/*
Copyright 2020 The actions-runner-controller authors.

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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
	actionsv1alpha1 "github.com/summerwind/actions-runner-controller/api/v1alpha1"
	"github.com/summerwind/actions-runner-controller/controllers"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

const (
	defaultRunnerImage = "chenrui/actions-runner:v2.165.2"
	defaultDockerImage = "docker:19.03.6-dind"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = actionsv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool

		runnerImage string
		dockerImage string

		ghToken             string
		ghAppID             int64
		ghAppInstallationID int64
		ghAppPrivateKey     string

		ghClient *github.Client
	)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&runnerImage, "runner-image", defaultRunnerImage, "The image name of self-hosted runner container.")
	flag.StringVar(&dockerImage, "docker-image", defaultDockerImage, "The image name of docker sidecar container.")
	flag.StringVar(&ghToken, "github-token", "", "The personal access token of GitHub.")
	flag.Int64Var(&ghAppID, "github-app-id", 0, "The application ID of GitHub App.")
	flag.Int64Var(&ghAppInstallationID, "github-app-installation-id", 0, "The installation ID of GitHub App.")
	flag.StringVar(&ghAppPrivateKey, "github-app-private-key", "", "The path of a private key file to authenticate as a GitHub App")
	flag.Parse()

	if ghToken == "" {
		ghToken = os.Getenv("GITHUB_TOKEN")
	}
	if ghAppID == 0 {
		appID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
		if err == nil {
			ghAppID = appID
		}
	}
	if ghAppInstallationID == 0 {
		appInstallationID, err := strconv.ParseInt(os.Getenv("GITHUB_APP_INSTALLATION_ID"), 10, 64)
		if err == nil {
			ghAppInstallationID = appInstallationID
		}
	}
	if ghAppPrivateKey == "" {
		ghAppPrivateKey = os.Getenv("GITHUB_APP_PRIVATE_KEY")
	}

	if ghAppID != 0 {
		if ghAppInstallationID == 0 {
			fmt.Fprintln(os.Stderr, "Error: The installation ID must be specified.")
			os.Exit(1)
		}

		if ghAppPrivateKey == "" {
			fmt.Fprintln(os.Stderr, "Error: The path of a private key file must be specified.")
			os.Exit(1)
		}

		tr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, ghAppID, ghAppInstallationID, ghAppPrivateKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid GitHub App credentials: %v\n", err)
			os.Exit(1)
		}
		ghClient = github.NewClient(&http.Client{Transport: tr})
	} else if ghToken != "" {
		tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghToken},
		))
		ghClient = github.NewClient(tc)
	} else {
		fmt.Fprintln(os.Stderr, "Error: GitHub App credentials or personal access token must be specified.")
		os.Exit(1)
	}

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	runnerReconciler := &controllers.RunnerReconciler{
		Client:       mgr.GetClient(),
		Log:          ctrl.Log.WithName("controllers").WithName("Runner"),
		Scheme:       mgr.GetScheme(),
		GitHubClient: ghClient,
		RunnerImage:  runnerImage,
		DockerImage:  dockerImage,
	}

	if err = runnerReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Runner")
		os.Exit(1)
	}

	runnerSetReconciler := &controllers.RunnerReplicaSetReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("RunnerReplicaSet"),
		Scheme: mgr.GetScheme(),
	}

	if err = runnerSetReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RunnerReplicaSet")
		os.Exit(1)
	}

	runnerDeploymentReconciler := &controllers.RunnerDeploymentReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("RunnerDeployment"),
		Scheme: mgr.GetScheme(),
	}

	if err = runnerDeploymentReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "RunnerDeployment")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
