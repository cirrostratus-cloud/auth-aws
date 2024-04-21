package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/cirrostratus-cloud/common/event"
	user_event "github.com/cirrostratus-cloud/oauth2/event"
	log "github.com/sirupsen/logrus"
)

type SNSEventBus struct {
	snsClient      *sns.Client
	topicArnPrefix string
	subscribers    map[event.EventName][]func(event event.Event) error
}

func NewSNSEventBus(snsClient *sns.Client, topicArnPrefix string) *SNSEventBus {
	return &SNSEventBus{
		snsClient:      snsClient,
		topicArnPrefix: topicArnPrefix,
		subscribers:    make(map[event.EventName][]func(event event.Event) error),
	}
}

func (e *SNSEventBus) getTopicArn(eventName event.EventName) string {
	return e.topicArnPrefix + strings.ReplaceAll(string(eventName), "/", "_")
}

func (e *SNSEventBus) Publish(eventName event.EventName, event event.Event) error {
	data, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}
	message := string(data)
	topicArn := e.getTopicArn(eventName)
	log.WithFields(log.Fields{
		"topicArn": topicArn,
		"message":  message,
	}).Info("Publishing event")
	_, err = e.snsClient.Publish(
		context.TODO(),
		&sns.PublishInput{
			Message:  &message,
			TopicArn: &topicArn,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *SNSEventBus) Subscribe(eventName event.EventName, suscriber func(event event.Event) error) error {
	if e.subscribers[eventName] == nil {
		e.subscribers[eventName] = make([]func(event event.Event) error, 0)
	}
	e.subscribers[eventName] = append(e.subscribers[eventName], suscriber)
	return nil
}

func (e *SNSEventBus) triggerAll(eventName event.EventName, event event.Event) error {
	for _, subscriber := range e.subscribers[eventName] {
		err := subscriber(event)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error": err,
				},
			).Error("Error in subscriber")
			return err
		}
	}
	return nil
}

func (e *SNSEventBus) triggerEvent(eventName event.EventName, payload string) error {
	log.WithFields(log.Fields{
		"EventName": eventName,
		"Payload":   payload,
	}).Info("Triggering event")

	switch eventName {
	case user_event.UserCreatedEventName:
		var userCreatedEvent user_event.UserCreatedEvent
		err := json.Unmarshal([]byte(payload), &userCreatedEvent)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error": err,
				},
			).Error("Error unmarshalling user created event")
			return err
		}
		return e.triggerAll(eventName, userCreatedEvent)

	case user_event.UserPasswordChangedEventName:
		var passwordChangedEvent user_event.PasswordChangedEvent
		err := json.Unmarshal([]byte(payload), &passwordChangedEvent)
		if err != nil {
			return err
		}
		return e.triggerAll(eventName, passwordChangedEvent)
	case user_event.UserPasswordRecoveredEventName:
		var userPasswordRecoveredEvent user_event.UserPasswordRecoveredEvent
		err := json.Unmarshal([]byte(payload), &userPasswordRecoveredEvent)
		if err != nil {
			return err
		}
		return e.triggerAll(eventName, userPasswordRecoveredEvent)
	case user_event.UserEmailConfirmedEventName:
		var userEmailConfirmedEvent user_event.UserEmailConfirmedEvent
		err := json.Unmarshal([]byte(payload), &userEmailConfirmedEvent)
		if err != nil {
			return err
		}
		return e.triggerAll(eventName, userEmailConfirmedEvent)
	default:
		return errors.New("event not found")
	}
}

func (e *SNSEventBus) Trigger(eventName event.EventName, payload string) error {
	log.WithFields(log.Fields{
		"EventName": eventName,
		"Payload":   payload,
	}).Debug("Triggering event")
	return e.triggerEvent(eventName, payload)
}
