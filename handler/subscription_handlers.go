package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rivermq/rivermq/model"
)

// CreateSubscriptionHandler does that
func CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var sub model.Subscription
	json.NewDecoder(r.Body).Decode(&sub)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resultSub, err := model.SaveSubscription(sub)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	} else {
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
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
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
	}
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
		w.WriteHeader(http.StatusNotFound)
	} else {
		deleteErr := model.DeleteSubscriptionByID(subID)
		if deleteErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
