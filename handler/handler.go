package handler

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rivermq/rivermq/metric"
	"github.com/rivermq/rivermq/model"
	"github.com/zeromq/goczmq"
)

var (
	handlerPushSocket *goczmq.Sock
	encoder           *gob.Encoder
)

func init() {
	handlerPushSocket, _ = goczmq.NewPush("inproc://handler")
	encoder = gob.NewEncoder(handlerPushSocket)
}

// CreateMessageHander accepts a message or delivery
func CreateMessageHander(w http.ResponseWriter, r *http.Request) {
	var msg model.Message
	json.NewDecoder(r.Body).Decode(&msg)
	msg.Status = model.StatusAccepted
	metric.MessageCtr.With(prometheus.Labels{"status": model.StatusAccepted}).Inc()
	resultMsg, err := model.SaveMessage(msg)

	go func(msg model.Message) {
		encoder.Encode(msg)
	}(resultMsg)

	if err != nil {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	} else {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusAccepted)}).Inc()
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(resultMsg); err != nil {
			metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		}
	}
}

// CreateSubscriptionHandler allows clients to create a Subscription in order
// to receive Message events
func CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var sub model.Subscription
	json.NewDecoder(r.Body).Decode(&sub)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resultSub, err := model.SaveSubscription(sub)
	if err != nil {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	} else {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusCreated)}).Inc()
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resultSub); err != nil {
			panic(err)
		}
	}
}

// FindAllSubscriptionsHandler returns a collection of all existing
// Subscriptions.  A `type` parameter allows the collection to be filtered
// by Message type
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
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusOK)}).Inc()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
}

// FindSubscriptionByIDHandler returns a Subscription matching the supplied
// ID
func FindSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID := vars["subID"]
	sub, err := model.FindSubscriptionByID(subID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusNotFound)}).Inc()
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
	}
	metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusOK)}).Inc()
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		panic(err)
	}
}

// DeleteSubscriptionByIDHandler allows the deletion of a Subscription with
// the supplied ID
func DeleteSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subID := vars["subID"]
	_, err := model.FindSubscriptionByID(subID)
	if err != nil {
		metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusNotFound)}).Inc()
		w.WriteHeader(http.StatusNotFound)
	} else {
		deleteErr := model.DeleteSubscriptionByID(subID)
		if deleteErr != nil {
			metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusInternalServerError)}).Inc()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		} else {
			metric.HTTPResponseCtr.With(prometheus.Labels{"code": string(http.StatusOK)}).Inc()
			w.WriteHeader(http.StatusOK)
		}
	}
}
