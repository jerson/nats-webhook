package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"nats_webhook/entities"
	"nats_webhook/modules/config"
	"nats_webhook/modules/nats"
	"net/http"
	"strings"
)

func init() {
	err := config.ReadDefault()
	if err != nil {
		panic(err)
	}

	if config.Vars.Debug {
		logrus.SetLevel(logrus.TraceLevel)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/event/{source}/{subject}", eventHandler)
	r.Use(authMiddleware)

	addr := fmt.Sprintf(":%d", config.Vars.Server.Port)
	logrus.Info("running on: ", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if config.Vars.API.Key == "" {
			next.ServeHTTP(w, r)
			return
		}

		authorization := r.Header.Get("Authorization")
		splitToken := strings.Split(authorization, "Bearer ")
		if len(splitToken) > 1 && strings.TrimSpace(splitToken[1]) == config.Vars.API.Key {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", errors.New("not allowed"))
	})
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	err = event(vars["source"], vars["subject"], body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", "ok")
}

func event(source, subject string, body []byte) error {
	guid := xid.New()
	log := logrus.WithFields(logrus.Fields{
		"id":      guid.String(),
		"source":  source,
		"subject": subject,
	})
	log.Debug("New event", string(body))
	cn, err := nats.Connect()
	if err != nil {
		log.Error(fmt.Errorf("connect error: %w", err))
		return err
	}
	defer cn.Close()

	payload := entities.Payload{
		ID:      guid.String(),
		Source:  source,
		Subject: subject,
		Body:    body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Error(fmt.Errorf("payload error: %w", err))
		return err
	}
	err = cn.Publish(subject, data)
	if err != nil {
		log.Warn(fmt.Errorf("publish error: %w", err))
		return err
	}

	log.Debug("Published")

	return nil
}
