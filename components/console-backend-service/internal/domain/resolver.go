package domain

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"math/rand"
	"time"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/eventing"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/apigateway"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/serverless"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/servicecatalogaddons"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/ui"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/experimental"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/application"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/k8s"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/rafter"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/servicecatalog"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/resource"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	ui  *ui.Resolver
	k8s *k8s.Resolver

	sc         *servicecatalog.PluggableContainer
	sca        *servicecatalogaddons.PluggableContainer
	app        *application.PluggableContainer
	rafter     *rafter.PluggableContainer
	ag         *apigateway.Resolver
	serverless *serverless.PluggableContainer
	eventing   *eventing.PluggableContainer
}

func GetRandomNumber() time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Intn(120)-60) * time.Second
}

func New(restConfig *rest.Config, appCfg application.Config, rafterCfg rafter.Config, serverlessCfg serverless.Config, informerResyncPeriod time.Duration, _ experimental.FeatureToggles, systemNamespaces []string) (*Resolver, error) {
	serviceFactory, err := resource.NewServiceFactoryForConfig(restConfig, informerResyncPeriod+GetRandomNumber())
	if err != nil {
		return nil, errors.Wrap(err, "while initializing service factory")
	}
	genericServiceFactory, err := resource.NewGenericServiceFactoryForConfig(restConfig, informerResyncPeriod+GetRandomNumber())
	if err != nil {
		return nil, errors.Wrap(err, "while initializing generic service factory")
	}

	uiContainer, err := ui.New(restConfig, informerResyncPeriod+GetRandomNumber())
	if err != nil {
		return nil, errors.Wrap(err, "while initializing UI resolver")
	}
	makePluggable := module.MakePluggableFunc(uiContainer.BackendModuleInformer)

	rafterContainer, err := rafter.New(serviceFactory, rafterCfg)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing Rafter resolver")
	}
	makePluggable(rafterContainer)

	scContainer, err := servicecatalog.New(restConfig, informerResyncPeriod+GetRandomNumber(), rafterContainer.Retriever)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing ServiceCatalog container")
	}
	makePluggable(scContainer)

	scaContainer, err := servicecatalogaddons.New(restConfig, informerResyncPeriod+GetRandomNumber(), scContainer.ServiceCatalogRetriever)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing ServiceCatalogAddons container")
	}
	makePluggable(scaContainer)

	appContainer, err := application.New(restConfig, appCfg, informerResyncPeriod+GetRandomNumber(), rafterContainer.Retriever)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing Application resolver")
	}
	makePluggable(appContainer)

	k8sResolver, err := k8s.New(restConfig, informerResyncPeriod+GetRandomNumber(), appContainer.ApplicationRetriever, scContainer.ServiceCatalogRetriever, scaContainer.ServiceCatalogAddonsRetriever, systemNamespaces)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing K8S resolver")
	}

	agResolver := apigateway.New(genericServiceFactory)
	makePluggable(agResolver)

	serverlessResolver, err := serverless.New(serviceFactory, serverlessCfg, scaContainer.ServiceCatalogAddonsRetriever)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing serverless resolver")
	}
	makePluggable(serverlessResolver)

	eventingResolver, err := eventing.New(serviceFactory)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing eventing resolver")
	}
	makePluggable(eventingResolver)

	return &Resolver{
		k8s:        k8sResolver,
		ui:         uiContainer.Resolver,
		sc:         scContainer,
		sca:        scaContainer,
		app:        appContainer,
		rafter:     rafterContainer,
		ag:         agResolver,
		serverless: serverlessResolver,
		eventing:   eventingResolver,
	}, nil
}

// WaitForCacheSync waits for caches to populate. This is blocking operation.
func (r *Resolver) WaitForCacheSync(stopCh <-chan struct{}) {
	// Not pluggable modules
	r.k8s.WaitForCacheSync(stopCh)
	r.ui.WaitForCacheSync(stopCh)

	// Pluggable modules
	r.sc.StopCacheSyncOnClose(stopCh)
	r.sca.StopCacheSyncOnClose(stopCh)
	r.app.StopCacheSyncOnClose(stopCh)
	r.rafter.StopCacheSyncOnClose(stopCh)
	r.ag.StopCacheSyncOnClose(stopCh)
	r.eventing.StopCacheSyncOnClose(stopCh)
	r.serverless.StopCacheSyncOnClose(stopCh)
}
