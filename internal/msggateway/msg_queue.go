package msggateway

import (
	"context"
	"fmt"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	"github.com/chdlvy/go-chatFrame/pkg/common/model"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
	"time"
)

const (
	privateExchange   = "privateExchange"
	groupExchange     = "groupExchange"
	broadcastExchange = "broadcastExchange"
	DLPrivateExchange = "dlPrivateExchange"
	DLGroupExchange   = "deadLetterGroupExchange"
)
const (
	//队列名称和绑定的前缀规则
	privateQueuePre      = "uqueue_"
	privateBindPre       = "uid_"
	privateConsumerPre   = "consumer_"
	DLprivateConsumerPre = "dlConsumer_"
	PrivateDLQueuePre    = "dlqueue_"
	DLprivateBindPre     = "DLkey_"
)

const (
	msgSeq = "seq_"
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
	return nil
}
func (mq *MsgQueue) createExchange() error {
	if err := mq.ch.ExchangeDeclare(broadcastExchange, "fanout", true, false, false, false, nil); err != nil {
		return err
	}
	if err := mq.ch.ExchangeDeclare(privateExchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if err := mq.ch.ExchangeDeclare(DLPrivateExchange, "direct", true, false, false, false, nil); err != nil {
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

func (mq *MsgQueue) NotificationPrivateMsg(ctx context.Context, UserID uint64, data []byte) error {
	//fmt.Println("(消息提醒)发送消息给：", client.UserID)
	UBind := privateBindPre + strconv.Itoa(int(UserID))
	return mq.ch.PublishWithContext(ctx,
		privateExchange,
		UBind,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
			Timestamp:    time.Now(),
		})
}

// 每注册一个用户就创建一个mqmember
// queueName: uqueue_xxx
// bind：uid_xxx
func (mq *MsgQueue) CreateMqMember(ctx context.Context, msgTTL int64, client *Client) error {
	//创建队列
	UID := strconv.Itoa(int(client.UserID))
	dlname := DLprivateBindPre + UID
	args := map[string]interface{}{
		"x-message-ttl":             msgTTL,
		"x-dead-letter-exchange":    DLPrivateExchange,
		"x-dead-letter-routing-key": dlname,
	}
	UQueue := privateQueuePre + UID
	//根据uid创建一个用户队列
	q, err := mq.ch.QueueDeclare(
		UQueue,
		true,
		false,
		false,
		false,
		args)
	if err != nil {
		return err
	}
	//绑定privateExchange
	UBind := privateBindPre + UID
	if err = mq.ch.QueueBind(q.Name,
		UBind,
		privateExchange,
		false,
		nil); err != nil {
		return err
	}

	mq.StartNotification(ctx, client)

	return nil
}

func (mq *MsgQueue) CreateDLQueueByMember(client *Client) error {

	UID := strconv.Itoa(int(client.UserID))
	qname := PrivateDLQueuePre + UID

	_, err := mq.ch.QueueDeclare(
		qname,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	//将dead_exchange与dead_queue进行绑定
	UBind := DLprivateBindPre + UID
	if err = mq.ch.QueueBind(qname, UBind, DLPrivateExchange, false, nil); err != nil {
		return err
	}
	return nil
}

func (mq *MsgQueue) GetOfflineMsg(ctx context.Context, client *Client) error {
	go func(client *Client) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ctx Done,client disconnect")
				return
			default:
				UID := strconv.Itoa(int(client.UserID))

				dlname := PrivateDLQueuePre + UID
				dlmsg, err := mq.ch.QueueDeclarePassive(dlname, true, false, false, false, nil)
				if err != nil {
					log.Fatal("GetOfflineMsg err：", err)
					//return err
				}
				//说明死信队列中有消息
				if dlmsg.Messages > 0 {
					//获取消息
					dlchan, err := mq.ch.Consume(dlname, DLprivateConsumerPre+UID, false, false, false, false, nil)
					if err != nil {
						log.Fatal(err)
					}
					for dlm := range dlchan {
						content := dlm.Body
						var data *model.MsgData
						json.Unmarshal(content, &data)
						client.writeMessage(1, []byte("接收旧消息----"))
						client.writeMessage(int(data.ContentType), content)
						dlm.Ack(false)
					}
				}
			}

		}
	}(client)
	return nil
}

func (mq *MsgQueue) StartNotification(ctx context.Context, client *Client) error {
	UID := strconv.Itoa(int(client.UserID))
	qname := privateQueuePre + UID
	//消费常规队列的消费者
	msgs, err := mq.ch.Consume(qname,
		privateConsumerPre+UID,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	//消费死信队列的消费者
	dlqname := PrivateDLQueuePre + UID
	dlmsg, err := mq.ch.Consume(dlqname, DLprivateConsumerPre+UID, false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			//如果用户断开连接，关闭协程
			case <-ctx.Done():
				fmt.Println("ctx Done,client disconnect")
				return
			case msg := <-msgs:
				//如果客户端离线则不ack
				client.writeMessage(1, msg.Body)
				msg.Ack(false)

			case dlm := <-dlmsg:
				content := dlm.Body
				fmt.Println(string(content))
				var data *model.MsgData
				json.Unmarshal(content, &data)
				client.writeMessage(1, []byte("接收旧消息----"))
				client.writeMessage(int(data.ContentType), content)
				dlm.Ack(false)
				//msg.Ack(false)
				//case dlm := <-dlmsg:
				//	content := dlm.Body
				//	fmt.Println(string(content))
				//	var data *model.MsgData
				//	json.Unmarshal(content, &data)
				//	client.writeMessage(1, []byte("接收旧消息----"))
				//	client.writeMessage(int(data.ContentType), content)
				//	dlm.Ack(false)
			}

		}
	}()
	return nil
}
func handleMessage(ctx context.Context, msg *amqp.Delivery, client *Client) {
	UID := strconv.Itoa(int(client.UserID))
	content := msg.Body
	var data *model.MsgData
	json.Unmarshal(content, &data)
	curseq := ctx.Value(msgSeq + UID).(int64)
	if data.Seq == curseq+1 {

	}

}

//func handleMessage(msg amqp.Delivery) {
//
//}

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
