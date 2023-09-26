package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL            string
	ExchangeName   string
	ExchangeSuffix string
	ExchangePrefix string
	AcK            bool
}

type IntegrationEvent struct {
	Data interface{} `json:"data"`
}

type RabbitMQ struct {
	conn   *amqp.Connection
	memory map[string]*amqp.Channel
	config RabbitMQConfig
	mutex  sync.Locker
}

type IRabbitMQ interface {
	Publish()
	AddSubscription()
	RemoveSubscription()
}

func NewRabbitMQ(config RabbitMQConfig) (*RabbitMQ, error) {

	// Context oluşturun ve bir süre sınırlayın (10 saniye)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Context kullanım sona erdiğinde iptal et

	// RabbitMQ bağlantısı kurma işlemi
	var conn *amqp.Connection
	var err error

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Bağlantı süresi aşıldı.")
			return nil, errors.New("Connection Timeout.")
		default:
			conn, err = amqp.Dial(config.URL)
			if err == nil {
				fmt.Println(fmt.Sprintf("RabbitMQ bağlantısı başarılı.State:%v", !conn.IsClosed()))
				return &RabbitMQ{
					conn:   conn,
					memory: make(map[string]*amqp.Channel),
					config: config,
				}, nil
			}
			fmt.Printf("Bağlantı kurma hatası: %v. Yeniden denenecek...\n", err)
			time.Sleep(1 * time.Second) // Bağlantı hatası durumunda 1 saniye bekleme
		}
	}

}

func (a *RabbitMQ) Close() {

	err := a.conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for _, channel := range a.memory {
		err = channel.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

}

func (a *RabbitMQ) AddSubscription(test interface{}, handler interface{}) (any, error) {
	exchange := a.config.ExchangeName
	queue := reflect.TypeOf(test).Name()
	queue = a.config.ExchangeSuffix + "." + strings.Replace(queue, a.config.ExchangePrefix, "", -1)

	_, ok := a.memory[queue]
	if ok {
		return nil, errors.New("Queue is Already Declared.")
	}

	channel, err := a.createChannel(10)

	if err != nil {
		fmt.Printf("State Channel:%v", !channel.IsClosed())
		return nil, err
	}

	err = a.declareExchange(channel, exchange, "direct")
	if err != nil {
		return nil, err
	}

	_, err = a.declareQueue(channel, queue)
	if err != nil {
		return nil, err
	}

	_, err = a.declareDeadLetterQueue(channel, queue)
	if err != nil {
		return nil, err
	}

	err = a.bindQueue(channel, queue, exchange)
	if err != nil {
		return nil, err
	}

	_, err = a.initializeCostumer(channel, queue, test, handler)
	if err != nil {
		return nil, err
	}
	a.memory[queue] = channel
	fmt.Println("Subscribed:", queue)
	return nil, nil
}

func (a *RabbitMQ) RemoveSubscription(queue string) {

}

func (a *RabbitMQ) Publish(exchange string, routingKey string, data interface{}) (any, error) {
	channel, err := a.createChannel(10)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	encodedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        encodedData,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (a *RabbitMQ) createChannel(prefetchCount int) (*amqp.Channel, error) {
	channel, err := a.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (a *RabbitMQ) declareDeadLetterQueue(channel *amqp.Channel, queueName string) (*amqp.Queue, error) {

	queue, err := channel.QueueDeclare(queueName+".deadLetter", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &queue, err
}

func (a *RabbitMQ) declareQueue(channel *amqp.Channel, queueName string) (*amqp.Queue, error) {
	arg := amqp.Table{"x-dead-letter-exchange": a.config.ExchangeName,
		"x-dead-letter-routing-key": queueName + ".deadLetter"}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, arg)
	if err != nil {
		return nil, err
	}

	return &queue, nil

}

func (a *RabbitMQ) declareExchange(channel *amqp.Channel, name string, kind string) error {
	err := channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
	return err
}

func (a *RabbitMQ) bindQueue(channel *amqp.Channel, queue string, exchange string) error {
	err := channel.QueueBind(queue, queue, a.config.ExchangeName, false, nil)
	return err
}

func (a *RabbitMQ) initializeCostumer(channel *amqp.Channel, queue string, param interface{}, handler interface{}) (any, error) {
	deliveries, err := channel.Consume(queue, queue, a.config.AcK, true, false, false, nil)
	if err != nil {
		return nil, err
	}
	go func(delivery <-chan amqp.Delivery) {
		for d := range delivery {
			// TODO: Make Parser

			handler.(func(amqp.Delivery))(d)
			if !a.config.AcK {
				d.Ack(false)
			}
		}
	}(deliveries)
	return true, nil
}
