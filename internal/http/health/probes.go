package health

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/helpers"
)

// ConditionCheck defines the API for functions that check the condition of the system
type ConditionCheck func() Condition

// Condition contains the information regarding the status of a piece of the System
type Condition struct {
	Ready   bool   `json:"ready"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name"`
}

type status struct {
	App      string      `json:"app"`
	Instance string      `json:"instance"`
	Status   []Condition `json:"status,omitempty"`
}

func alwaysGood() Condition {
	return Condition{
		Ready:   true,
		Message: "",
		Name:    "Default",
	}
}

// NewAPI creates a health check API
func NewAPI(app string, instance string) *Check {
	return &Check{
		app:              app,
		instance:         instance,
		healthConditions: []ConditionCheck{alwaysGood},
		readyConditions:  []ConditionCheck{alwaysGood},
	}
}

//Check implements the Check interface
type Check struct {
	app              string
	instance         string
	healthConditions []ConditionCheck
	readyConditions  []ConditionCheck
}

// RegisterHealthCondition registers functions that determin whether the systeam is healthy
func (h *Check) RegisterHealthCondition(condition ConditionCheck) {
	h.healthConditions = append(h.healthConditions, condition)
}

// RegisterReadynessCondition registers functions that determin whether the systeam is ready to be used
func (h *Check) RegisterReadynessCondition(condition ConditionCheck) {
	h.readyConditions = append(h.readyConditions, condition)
}

// AddHandlersTo add the liveness, readiness and info handlers to the router
func (h *Check) AddHandlersTo(router *mux.Router) {
	r := router.PathPrefix("/info").Subrouter()
	r.HandleFunc("/alive", h.livenessHandler)
	r.HandleFunc("/ready", h.readinessHandler)
	r.HandleFunc("/", h.aboutHandler)
}

func (h *Check) livenessHandler(w http.ResponseWriter, r *http.Request) {
	h.conditionHandler(w, h.healthConditions)

}

func (h *Check) readinessHandler(w http.ResponseWriter, r *http.Request) {
	h.conditionHandler(w, h.readyConditions)
}

func (h *Check) aboutHandler(w http.ResponseWriter, r *http.Request) {
	helpers.FormatResponse(w, &status{App: h.app, Instance: h.instance}, http.StatusOK)
}

func (h *Check) conditionHandler(w http.ResponseWriter, conditions []ConditionCheck) {
	ready := true
	appStatus := status{
		App:      h.app,
		Status:   []Condition{},
		Instance: h.instance,
	}
	for _, v := range conditions {
		condition := v()
		appStatus.Status = append(appStatus.Status, condition)
		if !condition.Ready {
			ready = false
		}
	}
	code := http.StatusOK
	if !ready {
		code = http.StatusBadRequest
	}
	helpers.FormatResponse(w, appStatus, code)
}
