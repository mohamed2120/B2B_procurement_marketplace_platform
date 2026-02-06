package events

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// RedisEventBus implements EventBus using Redis pub/sub
type RedisEventBus struct {
	client  *redis.Client
	channel string
	pubsub  *redis.PubSub
}

// NewRedisEventBus creates a new Redis-based event bus
func NewRedisEventBus(redisClient *redis.Client) *RedisEventBus {
	channel := os.Getenv("REDIS_PUBSUB_CHANNEL")
	if channel == "" {
		channel = "b2b_events"
	}

	return &RedisEventBus{
		client:  redisClient,
		channel: channel,
	}
}

// Publish publishes an event to the event bus
func (r *RedisEventBus) Publish(ctx interface{}, event *EventEnvelope) error {
	ctxVal, ok := ctx.(context.Context)
	if !ok {
		ctxVal = context.Background()
	}

	data, err := event.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Publish to Redis channel
	err = r.client.Publish(ctxVal, r.channel, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Subscribe subscribes to events of a specific type
func (r *RedisEventBus) Subscribe(ctx interface{}, eventType EventType, handler func(*EventEnvelope) error) error {
	ctxVal, ok := ctx.(context.Context)
	if !ok {
		ctxVal = context.Background()
	}

	// Subscribe to the channel
	pubsub := r.client.Subscribe(ctxVal, r.channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctxVal.Done():
			return ctxVal.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			envelope, err := DeserializeEventEnvelope([]byte(msg.Payload))
			if err != nil {
				// Log error but continue processing
				fmt.Printf("Error deserializing event: %v\n", err)
				continue
			}

			// Only process events of the requested type
			if envelope.Type == eventType {
				if err := handler(envelope); err != nil {
					fmt.Printf("Error handling event: %v\n", err)
				}
			}
		}
	}
}

// SubscribeAll subscribes to all events
func (r *RedisEventBus) SubscribeAll(ctx interface{}, handler func(*EventEnvelope) error) error {
	ctxVal, ok := ctx.(context.Context)
	if !ok {
		ctxVal = context.Background()
	}

	pubsub := r.client.Subscribe(ctxVal, r.channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctxVal.Done():
			return ctxVal.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			envelope, err := DeserializeEventEnvelope([]byte(msg.Payload))
			if err != nil {
				fmt.Printf("Error deserializing event: %v\n", err)
				continue
			}

			if err := handler(envelope); err != nil {
				fmt.Printf("Error handling event: %v\n", err)
			}
		}
	}
}
