package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const MQURL = "amqp://test:test@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// bind key
	Key string
	// 连接信息
	Mqurl string
}

// 新建RabbitMQ实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key}
}

// 断开连接
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rmq := NewRabbitMQ(queueName, "", "")
	var err error

	// get connect
	rmq.conn, err = amqp.Dial(rmq.Mqurl)
	rmq.failOnErr(err, "Failed to connect to rabbitmq")

	// get channel
	rmq.channel, err = rmq.conn.Channel()
	rmq.failOnErr(err, "Failed to open a channel")

	return rmq
}

func (r *RabbitMQ) PublishSimple(message string) {
	// 1. 申请队列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 调用channel, 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// if true, 根据exchange,routekey规则无法找到queue时, 把消息返回给发送者
		false,
		// if true, 当queue上没有消费者时, 把消息返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) ConsumeSimple() {
	// 1. 申请队列
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 接受消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Recived a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for message. To exit press CTRL+C")
	<-forever
}

// 订阅模式下的publish
func (r *RabbitMQ) PublishPub(message string) {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true, false, false, false, nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	//2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 订阅模式下的consume
func (r *RabbitMQ) ReciveSub() {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true, false,
		// YES表示该exchange不能被client用来推送消息, 仅用来进行exchange之间的绑定
		false,
		false, nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	//2. 尝试创建队列, 这里队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	//3. 绑定到exchagne中
	err = r.channel.QueueBind(
		q.Name,
		// 在发布/订阅模式下, 这里的key要为空
		"",
		r.Exchange,
		false,
		nil,
	)

	//4. 消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range message {
			log.Printf("Recived a message: %s", d.Body)
		}
	}()
	<-forever
}

// 路由模式
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	rmq := NewRabbitMQ("", exchangeName, routingKey)

	var err error
	rmq.conn, err = amqp.Dial(rmq.Mqurl)
	rmq.failOnErr(err, "Failed to connect rabbitmq")

	rmq.channel, err = rmq.conn.Channel()
	rmq.failOnErr(err, "Failed to open a channel")
	return rmq
}
