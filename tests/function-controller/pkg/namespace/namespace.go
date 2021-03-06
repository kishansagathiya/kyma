package namespace

import (
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	k8sretry "k8s.io/client-go/util/retry"
	eventingv1alpha1 "knative.dev/eventing/pkg/apis/eventing/v1alpha1"

	"github.com/kyma-project/kyma/tests/function-controller/pkg/retry"
	"github.com/kyma-project/kyma/tests/function-controller/pkg/shared"
)

const (
	TestNamespaceLabelKey   = "created-by"
	TestNamespaceLabelValue = "serverless-controller-manager-test"
)

type Namespace struct {
	coreCli typedcorev1.CoreV1Interface
	name    string
	log     shared.Logger
	verbose bool
}

func New(name string, coreCli typedcorev1.CoreV1Interface, container shared.Container) *Namespace {
	return &Namespace{coreCli: coreCli, name: name, log: container.Log, verbose: container.Verbose}
}

func (n Namespace) GetName() string {
	return n.name
}

func (n *Namespace) Create() (string, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: n.name,
			Labels: map[string]string{
				eventingv1alpha1.InjectionAnnotation: "enabled",               // https://knative.dev/v0.12-docs/eventing/broker-trigger/#annotation
				TestNamespaceLabelKey:                TestNamespaceLabelValue, // convenience for cleaning up stale namespaces during development
			},
		},
	}

	backoff := wait.Backoff{
		Duration: 500 * time.Millisecond,
		Factor:   2,
		Jitter:   0.1,
		Steps:    4,
	}

	err := k8sretry.OnError(backoff, func(err error) bool {
		return true
	}, func() error {
		_, err := n.coreCli.Namespaces().Create(ns)
		if apierrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	})
	if err != nil {
		return n.name, errors.Wrapf(err, "while creating namespace %s", n.name)
	}

	n.log.Logf("CREATE: namespace %s", n.name)
	if n.verbose {
		n.log.Logf("%+v", ns)
	}
	return n.name, nil
}

func (n *Namespace) Delete() error {
	err := retry.WithIgnoreOnNotFound(retry.DefaultBackoff, func() error {
		if n.verbose {
			n.log.Logf("DELETE: namespace: %s", n.name)
		}
		return n.coreCli.Namespaces().Delete(n.name, &metav1.DeleteOptions{})
	}, n.log)
	if err != nil {
		return errors.Wrapf(err, "while deleting namespace %s", n.name)
	}
	return nil
}
