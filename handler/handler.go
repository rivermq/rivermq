package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rivermq/rivermq/inspect"
	"github.com/rivermq/rivermq/model"
)

var (
	httpResponseCtr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_resp_total",
			Help: "Number of http responses",
		},
		[]string{"code"},
	)

	// MessageCtr allows the counting of Messages and their status
	MessageCtr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_total",
			Help: "Number of messages received",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(httpResponseCtr)
	prometheus.MustRegister(MessageCtr)
}

// CreateMessageHander accepts a message or delivery
func CreateMessageHander(w http.ResponseWriter, r *http.Request) {
	var msg model.Message
	json.NewDecoder(r.Body).Decode(&msg)
	msg.Status = model.StatusAccepted
	MessageCtr.With(prometheus.Labels{"status": model.StatusAccepted}).Inc()
	resultMsg, err := model.SaveMessage(msg)
	go func(msg model.Message) {
		inspect.HandleMessage(msg)
	}(resultMsg)
	if err != nil {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	} else {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusAccepted)}).Inc()
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(resultMsg); err != nil {
			panic(err)
		}
	}
}

// CreateSubscriptionHandler does that
func CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var sub model.Subscription
	json.NewDecoder(r.Body).Decode(&sub)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resultSub, err := model.SaveSubscription(sub)
	if err != nil {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	} else {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusCreated)}).Inc()
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resultSub); err != nil {
			panic(err)
		}
	}
}

// FindAllSubscriptionsHandler does that
func FindAllSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check for type query parameter
	msgType := r.FormValue("type")
	var subs []model.Subscription
	var err error
	if msgType != "" {
		subs, err = model.FindAllSubscriptionsByType(msgType)
	} else {
		subs, err = model.FindAllSubscriptions()
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusOK)}).Inc()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		panic(err)
	}
}

// FindSubscriptionByIDHandler does that
func FindSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID := vars["subID"]
	sub, err := model.FindSubscriptionByID(subID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusNotFound)}).Inc()
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
	}
	httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusOK)}).Inc()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		panic(err)
	}
}

// DeleteSubscriptionByIDHandler does that
func DeleteSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID := vars["subID"]
	_, err := model.FindSubscriptionByID(subID)
	if err != nil {
		httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusNotFound)}).Inc()
		w.WriteHeader(http.StatusNotFound)
	} else {
		deleteErr := model.DeleteSubscriptionByID(subID)
		if deleteErr != nil {
			httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		} else {
			httpResponseCtr.With(prometheus.Labels{"code": string(http.StatusOK)}).Inc()
			w.WriteHeader(http.StatusOK)
		}
	}
}
