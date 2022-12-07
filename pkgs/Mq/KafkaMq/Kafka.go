package KafkaMq

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"time"
)

var m *MsgQueue

type MsgQueue struct {
	Service ConsumerI
	Config  *Config
}

func init() {
	m = New()
}

func New() *MsgQueue { return m.New() }
func (m *MsgQueue) New() *MsgQueue {
	v := new(MsgQueue)
	return v
}

func AddConfig(topic, host []string, group string) { m.AddConfig(topic, host, group) }
func (m *MsgQueue) AddConfig(topic, host []string, group string) {
	c := &Config{
		Topic: topic,
		Host:  host,
		Group: group,
	}
	m.Config = c
}

func AddConsumer(service ConsumerI) { m.AddConsumer(service) }
func (m *MsgQueue) AddConsumer(service ConsumerI) {
	m.Service = service
}

func (m *MsgQueue) ConsumeLoop() {
	for {
		m.Consumer()
		time.Sleep(3 * time.Second)
		fmt.Println("reconnect kafka publish...")
	}
}

func (m *MsgQueue) Consumer() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(m.Config.Host, m.Config.Group, m.Config.Topic, config)
	if err != nil {
		fmt.Printf("Failed to start consumer: %s", err)
		return
	}
	defer consumer.Close()

	// Consume all channels, wait for signal to exit
	for {
		select {
		case msg, more := <-consumer.Messages():
			var mqMsg Msg
			mqMsg.Topic = msg.Topic
			mqMsg.Msg = string(msg.Value)
			m.Service.Consume(mqMsg)

			if more {
				//fmt.Printf("%s/%d/%d\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
				consumer.MarkOffset(msg, "")
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				fmt.Printf("Rebalanced: %+v\n", ntf)
			}
		case err, more := <-consumer.Errors():
			if more {
				fmt.Printf("Error: %s\n", err.Error())
			}
			break
		}
	}
}

func (m *MsgQueue) Producer(topic, msg string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	pMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	producer, err := sarama.NewSyncProducer(m.Config.Host, config)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return err
	}
	defer producer.Close()
	_, _, err = producer.SendMessage(pMsg)
	if err != nil {
		return err
	}
	return nil
}