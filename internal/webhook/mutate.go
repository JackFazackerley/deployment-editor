package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	v1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	mutateLog = ctrl.Log.WithName("mutate")
)

//+kubebuilder:webhook:webhookVersions={v1beta1},path=/mutate,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps,resources=deployments,verbs=create;update,versions=v1,name=jackfazackerley.hello.com,admissionReviewVersions={v1,v1beta1}
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

type Mutate struct {
	Client   client.Client
	Replicas int32
	decoder  *admission.Decoder
}

func (m *Mutate) Handle(ctx context.Context, req admission.Request) admission.Response {
	deployment := &v1.Deployment{}
	err := m.decoder.Decode(req, deployment)
	if err != nil {
		mutateLog.Error(err, "decoding admission request")
		return admission.Errored(http.StatusBadRequest, err)
	}

	if deployment.Annotations != nil {
		value, ok := deployment.Annotations["jackfazackerley.com/should-edit-replicas"]
		if ok {
			should, err := strconv.ParseBool(value)
			if err != nil {
				should = false
			}

			if should {
				deployment.Spec.Replicas = &m.Replicas
			}
		}
	}

	marshalledDeployment, err := json.Marshal(deployment)
	if err != nil {
		mutateLog.Error(err, "marshalling deployment")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshalledDeployment)
}

func (m *Mutate) InjectDecoder(d *admission.Decoder) error {
	m.decoder = d
	return nil
}
