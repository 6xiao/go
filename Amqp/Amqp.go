package Amqp

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
)

var (
	flgPrefetch = flag.Int("prefetch", 64, "prefetch message from mq")
)

func ConnectMq(url string) (conn *amqp.Connection, channel *amqp.Channel, err error) {
	conn, err = amqp.Dial(url)
	if err == nil {
		channel, err = conn.Channel()
	}

	return
}

func NewMqExchange(channel *amqp.Channel, name, _type string, durable bool) error {
	return channel.ExchangeDeclare(
		name,    // name
		_type,   // type
		durable, // durable
		false,    // auto-delete
		false,   // internal
		false,   // nowait
		nil,     // args
	)
}

func Publish(channel *amqp.Channel, exchange, rkey string, msg []byte) error {
	if channel == nil {
		return fmt.Errorf("channel is nil")
	}

	return channel.Publish(exchange, rkey, false, false,
		amqp.Publishing{ContentType: "application/octet-stream", Body: msg})
}

func newMqConsumer(url, exchange, queue, rkey, ctag string, ack, durable, exclusive bool) (
	conn *amqp.Connection, channel *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {

	conn, channel, err = ConnectMq(url)
	if err != nil {
		return
	}

	channel.QueueDeclare(
		queue,     // name of the queue
		durable,   // durable
		exclusive, // delete when usused
		exclusive, // exclusive
		false,     // noWait
		nil,       // arguments
	)

	if err = channel.QueueBind(
		queue,    // name of the queue
		rkey,     // bindingKey
		exchange, // sourceExchange
		false,    // noWait
		nil,      // arguments
	); err != nil {
		conn.Close()
		return
	}

	deliveries, err = channel.Consume(
		queue,     // name
		ctag,      // consumerTag,
		!ack,      // noAck
		exclusive, // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		conn.Close()
		return
	}

	if ack {
		channel.Qos(*flgPrefetch, 0, true)
	}
	return
}

func NewMqConsumer(url, exchange, queue, rkey, ctag string, ack, durable bool) (
	conn *amqp.Connection, channel *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {

	return newMqConsumer(url, exchange, queue, rkey, ctag, ack, durable, false)
}

func NewExclusiveMqConsumer(url, exchange, queue, rkey, ctag string) (
	conn *amqp.Connection, channel *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {

	return newMqConsumer(url, exchange, queue, rkey, ctag, false, false, true)
}
