/*
Copyright 2022 The NanoBus Authors.

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

package dapr

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dapr/dapr/pkg/acl"
	global_config "github.com/dapr/dapr/pkg/config"
	env "github.com/dapr/dapr/pkg/config/env"
	"github.com/dapr/dapr/pkg/cors"
	"github.com/dapr/dapr/pkg/grpc"
	"github.com/dapr/dapr/pkg/metrics"
	"github.com/dapr/dapr/pkg/modes"
	"github.com/dapr/dapr/pkg/operator/client"
	"github.com/dapr/dapr/pkg/resiliency"
	"github.com/dapr/dapr/pkg/runtime"
	"github.com/dapr/dapr/pkg/runtime/security"
	"github.com/dapr/dapr/pkg/version"
	"github.com/dapr/kit/logger"
)

type Runtime struct {
	mode                     string
	daprHTTPPort             string
	daprAPIListenAddresses   string
	daprInternalGRPCPort     string
	profilePort              string
	componentsPath           string
	config                   string
	appID                    string
	controlPlaneAddress      string
	sentryAddress            string
	placementServiceHostAddr string
	allowedOrigins           string
	enableProfiling          bool
	runtimeVersion           bool
	buildInfo                bool
	appMaxConcurrency        int
	enableMTLS               bool
	appSSL                   bool
	daprHTTPMaxRequestSize   int

	loggerOptions   logger.Options
	metricsExporter metrics.Exporter
	options         []runtime.Option
	runtime         *runtime.DaprRuntime
}

var log = logger.NewLogger("dapr.runtime")

func New() *Runtime {
	return &Runtime{}
}

func (r *Runtime) AttachFlags() {
	// Dapr flags
	flag.StringVar(&r.mode, "mode", string(modes.StandaloneMode), "Runtime mode for Dapr")
	flag.StringVar(&r.daprHTTPPort, "dapr-http-port", "0", "HTTP port for Dapr API to listen on")
	flag.StringVar(&r.daprAPIListenAddresses, "dapr-listen-addresses", "127.0.0.1", "One or more addresses for the Dapr API to listen on, CSV limited")
	//daprAPIGRPCPort := flag.String("dapr-grpc-port", fmt.Sprintf("%v", DefaultDaprAPIGRPCPort), "gRPC port for the Dapr API to listen on")
	flag.StringVar(&r.daprInternalGRPCPort, "dapr-internal-grpc-port", "", "gRPC port for the Dapr Internal API to listen on")
	//appPort := flag.String("app-port", "", "The port the application is listening on")
	flag.StringVar(&r.profilePort, "profile-port", fmt.Sprintf("%v", runtime.DefaultProfilePort), "The port for the profile server")
	//appProtocol := flag.String("app-protocol", string(HTTPProtocol), "Protocol for the application: grpc or http")
	flag.StringVar(&r.componentsPath, "components-path", "", "Path for components directory. If empty, components will not be loaded. Self-hosted mode only")
	flag.StringVar(&r.config, "config", "", "Path to config file, or name of a configuration object")
	flag.StringVar(&r.appID, "app-id", "", "A unique ID for Dapr. Used for Service Discovery and state")
	flag.StringVar(&r.controlPlaneAddress, "control-plane-address", "", "Address for a Dapr control plane")
	flag.StringVar(&r.sentryAddress, "sentry-address", "", "Address for the Sentry CA service")
	flag.StringVar(&r.placementServiceHostAddr, "placement-host-address", "", "Addresses for Dapr Actor Placement servers")
	flag.StringVar(&r.allowedOrigins, "allowed-origins", cors.DefaultAllowedOrigins, "Allowed HTTP origins")
	flag.BoolVar(&r.enableProfiling, "enable-profiling", false, "Enable profiling")
	flag.BoolVar(&r.runtimeVersion, "version", false, "Prints the runtime version")
	flag.BoolVar(&r.buildInfo, "build-info", false, "Prints the build info")
	//waitCommand := flag.Bool("wait", false, "wait for Dapr outbound ready")
	flag.IntVar(&r.appMaxConcurrency, "app-max-concurrency", -1, "Controls the concurrency level when forwarding requests to user code")
	flag.BoolVar(&r.enableMTLS, "enable-mtls", false, "Enables automatic mTLS for daprd to daprd communication channels")
	flag.BoolVar(&r.appSSL, "app-ssl", false, "Sets the URI scheme of the app to https and attempts an SSL connection")
	flag.IntVar(&r.daprHTTPMaxRequestSize, "dapr-http-max-request-size", -1, "Increasing max size of request body in MB to handle uploading of big files. By default 4 MB.")

	r.loggerOptions = logger.DefaultOptions()
	r.loggerOptions.AttachCmdFlags(flag.StringVar, flag.BoolVar)

	r.metricsExporter = metrics.NewExporter(metrics.DefaultMetricNamespace)
	r.metricsExporter.Options().AttachCmdFlags(flag.StringVar, flag.BoolVar)
}

func (r *Runtime) AddOptions(options ...runtime.Option) {
	r.options = append(r.options, options...)
}

func (r *Runtime) Initialize() error {
	if r.runtimeVersion {
		fmt.Println(version.Version())
		os.Exit(0)
	}

	if r.buildInfo {
		fmt.Printf("Version: %s\nGit Commit: %s\nGit Version: %s\n", version.Version(), version.Commit(), version.GitVersion())
		os.Exit(0)
	}

	if r.appID == "" {
		return errors.New("app-id parameter cannot be empty")
	}

	// Apply options to all loggers
	r.loggerOptions.SetAppID(r.appID)
	if err := logger.ApplyOptionsToLoggers(&r.loggerOptions); err != nil {
		return err
	}

	log.Infof("starting Dapr Runtime -- version %s -- commit %s", version.Version(), version.Commit())
	log.Infof("log level set to: %s", r.loggerOptions.OutputLevel)

	// Initialize dapr metrics exporter
	if err := r.metricsExporter.Init(); err != nil {
		return err
	}

	profPort, err := strconv.Atoi(r.profilePort)
	if err != nil {
		return fmt.Errorf("error parsing profile-port flag: %w", err)
	}

	var daprInternalGRPC int
	if r.daprInternalGRPCPort != "" {
		daprInternalGRPC, err = strconv.Atoi(r.daprInternalGRPCPort)
		if err != nil {
			log.Fatalf("error parsing dapr-internal-grpc-port: %w", err)
		}
	} else {
		daprInternalGRPC, err = grpc.GetFreePort()
		if err != nil {
			log.Fatalf("failed to get free port for internal grpc server: %w", err)
		}
	}

	var maxRequestBodySize int
	if r.daprHTTPMaxRequestSize != -1 {
		maxRequestBodySize = r.daprHTTPMaxRequestSize
	} else {
		maxRequestBodySize = runtime.DefaultMaxRequestBodySize
	}

	placementAddresses := []string{}
	if r.placementServiceHostAddr != "" {
		placementAddresses = parsePlacementAddr(r.placementServiceHostAddr)
	}

	var concurrency int
	if r.appMaxConcurrency != -1 {
		concurrency = r.appMaxConcurrency
	}

	daprAPIHTTPPort, _ := strconv.Atoi(r.daprHTTPPort)
	daprAPIListenAddressList := strings.Split(r.daprAPIListenAddresses, ",")
	if len(daprAPIListenAddressList) == 0 && daprAPIHTTPPort != 0 {
		daprAPIListenAddressList = []string{"127.0.0.1"}
		//[]string{runtime.DefaultAPIListenAddress}
	}

	runtimeConfig := runtime.NewRuntimeConfig(
		r.appID,
		placementAddresses,
		r.controlPlaneAddress,
		r.allowedOrigins,
		r.config,
		r.componentsPath,
		string(runtime.EmbeddedProtocol),
		r.mode,
		daprAPIHTTPPort,
		daprInternalGRPC,
		0,
		daprAPIListenAddressList,
		nil,
		0,
		profPort,
		r.enableProfiling,
		concurrency,
		r.enableMTLS,
		r.sentryAddress,
		r.appSSL,
		maxRequestBodySize,
		"",
		4,
		false,
		5*time.Second,
		false,
	)

	// // set environment variables
	// // TODO - consider adding host address to runtime config and/or caching result in utils package
	// host, err := utils.GetHostAddress()
	// if err != nil {
	// 	log.Warnf("failed to get host address, env variable %s will not be set", env.HostAddress)
	// }

	variables := map[string]string{
		env.AppID: r.appID,
		//env.AppPort:         *appPort,
		//env.HostAddress:     host,
		env.DaprPort: strconv.Itoa(daprInternalGRPC),
		// env.DaprGRPCPort:    *daprAPIGRPCPort,
		env.DaprHTTPPort:    r.daprHTTPPort,
		env.DaprMetricsPort: r.metricsExporter.Options().Port, // TODO - consider adding to runtime config
		env.DaprProfilePort: r.profilePort,
	}

	if err = setEnvVariables(variables); err != nil {
		log.Fatal(err)
	}

	var globalConfig *global_config.Configuration
	var configErr error

	if r.enableMTLS {
		runtimeConfig.CertChain, err = security.GetCertChain()
		if err != nil {
			return err
		}
	}

	var accessControlList *global_config.AccessControlList
	var namespace string
	var podName string

	if r.config != "" {
		switch modes.DaprMode(r.mode) {
		case modes.KubernetesMode:
			client, conn, clientErr := client.GetOperatorClient(r.controlPlaneAddress, security.TLSServerName, runtimeConfig.CertChain)
			if clientErr != nil {
				return err
			}
			defer conn.Close()
			namespace = os.Getenv("NAMESPACE")
			podName = os.Getenv("POD_NAME")
			globalConfig, configErr = global_config.LoadKubernetesConfiguration(r.config, namespace, podName, client)
		case modes.StandaloneMode:
			globalConfig, _, configErr = global_config.LoadStandaloneConfiguration(r.config)
		}

		if configErr != nil {
			log.Debugf("Config error: %v", configErr)
		}
	}

	if configErr != nil {
		return fmt.Errorf("error loading configuration: %w", configErr)
	}
	if globalConfig == nil {
		log.Info("loading default configuration")
		globalConfig = global_config.LoadDefaultConfiguration()
	}

	accessControlList, err = acl.ParseAccessControlSpec(globalConfig.Spec.AccessControlSpec, string(runtimeConfig.ApplicationProtocol))
	if err != nil {
		return err
	}
	resiliencyProvider := resiliency.NoOp{}
	r.runtime, err = runtime.NewDaprRuntime(runtimeConfig, globalConfig, accessControlList, &resiliencyProvider), nil
	return err
}

func (r *Runtime) Run() error {
	if err := r.runtime.Run(r.options...); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) Shutdown(d time.Duration) {
	r.runtime.Shutdown(d)
}

func (r *Runtime) WaitUntilShutdown() error {
	return r.runtime.WaitUntilShutdown()
}

func setEnvVariables(variables map[string]string) error {
	for key, value := range variables {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func parsePlacementAddr(val string) []string {
	parsed := []string{}
	p := strings.Split(val, ",")
	for _, addr := range p {
		parsed = append(parsed, strings.TrimSpace(addr))
	}
	return parsed
}
