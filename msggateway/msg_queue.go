package msggateway

import (
	"context"
	"fmt"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
)

const (
	privateExchange           = "privateExchange"
	groupExchange             = "groupExchange"
	broadcastExchange         = "broadcastExchange"
	deadLetterPrivateExchange = "deadLetterPrivateExchange"
	deadLetterGroupExchange   = "deadLetterGroupExchange"
)

type MsgQueue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

var MQ *MsgQueue

func NewMsgQueue() (*MsgQueue, error) {
	if MQ != nil {
		return MQ, nil
	}
	url := fmt.Sprintf("amqp://%s:%s@%s", config.Config.RabbitMq.Username, config.Config.RabbitMq.Password, config.Config.RabbitMq.Address)
	fmt.Println("rabbitmq dial：", config.Config.RabbitMq.Username, config.Config.RabbitMq.Password, config.Config.RabbitMq.Address)
	conn, err := amqp.Dial(url)
	if err != nil {
		return &MsgQueue{}, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return &MsgQueue{}, err
	}

	//defer conn.Close()
	MQ = &MsgQueue{
		conn: conn,
		ch:   ch,
	}
	return MQ, nil
}
func (mq *MsgQueue) InitMQ() error {
	if err := mq.createExchange(); err != nil {
		return err
	}
	//if err := mq.createDeadLetterExchange(); err != nil {
	//	return err
	//}
	return nil
}
func (mq *MsgQueue) createExchange() error {
	if err := mq.ch.ExchangeDeclare(broadcastExchange, "fanout", true, false, false, false, nil); err != nil {
		return err
	}
	if err := mq.ch.ExchangeDeclare(privateExchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	return nil
}

// 创建组的时候创建一个groupExchange
func (mq *MsgQueue) CreateGroupExchange(groupID uint64) error {
	gid := strconv.Itoa(int(groupID))
	if err := mq.ch.ExchangeDeclare(groupExchange+"_"+gid, "fanout", true, false, false, false, nil); err != nil {
		return err
	}
	return nil
}

//重新考虑死信
//func (mq *MsgQueue) createDeadLetterExchange() error {
//	if err := mq.ch.ExchangeDeclare(deadLetterPrivateExchange, "direct", true, false, false, false, nil); err != nil {
//		return err
//	}
//	if err := mq.ch.ExchangeDeclare(deadLetterGroupExchange, "direct", true, false, false, false, nil); err != nil {
//		return err
//	}
//	return nil
//}

// 每注册一个用户就创建一个mqmember
func (mq *MsgQueue) CreateMqMember(msgTTL uint64, client *Client) error {
	fmt.Println("createMqMember：", client.UserID)
	//创建队列
	args := map[string]interface{}{
		"x-message-ttl": msgTTL,
	}
	UID := strconv.Itoa(int(client.UserID))
	//根据uid创建一个用户队列
	q, err := mq.ch.QueueDeclare(
		UID,
		true,
		false,
		false,
		false,
		args)
	if err != nil {
		return err
	}
	//绑定privateExchange
	if err := mq.ch.QueueBind(q.Name,
		UID,
		privateExchange,
		false,
		nil); err != nil {
		return err
	}

	//绑定消费者
	msgs, err := mq.ch.Consume(q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			//client.writeMessage()
			fmt.Println("rabbitmq consumer get a msg：", msg)
		}
	}()

	return nil
}

// 加入群聊，添加绑定
func (mq *MsgQueue) JoinGroup(userID, groupID uint64) error {
	UID := strconv.Itoa(int(userID))
	GID := strconv.Itoa(int(groupID))
	//绑定groupExchange
	if err := mq.ch.QueueBind(UID,
		"",
		groupExchange+"_"+GID,
		false,
		nil); err != nil {
		return err
	}
	return nil
}

// 离开群聊，解除绑定
func (mq *MsgQueue) LeaveGroup(userID, groupID uint64) error {
	UID := strconv.Itoa(int(userID))
	GID := strconv.Itoa(int(groupID))
	//绑定groupExchange
	if err := mq.ch.QueueUnbind(UID,
		"",
		groupExchange+"_"+GID,
		nil); err != nil {
		return err
	}
	return nil
}

func (mq *MsgQueue) NotificationPrivateMsg(ctx context.Context, client *Client, data []byte) error {
	fmt.Println("(消息提醒)发送消息给：", client.UserID)
	UID := strconv.Itoa(int(client.UserID))
	return mq.ch.PublishWithContext(ctx,
		privateExchange,
		UID,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
}
