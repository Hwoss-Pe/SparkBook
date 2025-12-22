package events

import (
	commentevents "Webook/comment/events"
	followevents "Webook/follow/events"
	"Webook/interactive/repository/dao"
	"Webook/pkg/logger"
	"Webook/pkg/saramax"
	"time"

	"github.com/IBM/sarama"
	xcontext "golang.org/x/net/context"
)

type NotificationEventConsumer struct {
	client       sarama.Client
	l            logger.Logger
	notifDAO     dao.NotificationDAO
	resolveOwner func(ctx xcontext.Context, biz string, bizId int64) (int64, error)
}

func NewNotificationEventConsumer(client sarama.Client, l logger.Logger, notifDAO dao.NotificationDAO, resolveOwner func(ctx xcontext.Context, biz string, bizId int64) (int64, error)) *NotificationEventConsumer {
	return &NotificationEventConsumer{client: client, l: l, notifDAO: notifDAO, resolveOwner: resolveOwner}
}

func (c *NotificationEventConsumer) Start() error {
	cgLike, err := sarama.NewConsumerGroupFromClient("notification_like", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cgLike.Consume(xcontext.Background(), []string{TopicLike}, saramax.NewHandler[LikeEvent](c.l, c.consumeLike))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", TopicLike))
			}
		}
	}()

	cgCollect, err := sarama.NewConsumerGroupFromClient("notification_collect", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cgCollect.Consume(xcontext.Background(), []string{TopicCollect}, saramax.NewHandler[CollectEvent](c.l, c.consumeCollect))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", TopicCollect))
			}
		}
	}()

	cgFollow, err := sarama.NewConsumerGroupFromClient("notification_follow", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cgFollow.Consume(xcontext.Background(), []string{followevents.TopicFollow}, saramax.NewHandler[followevents.FollowEvent](c.l, c.consumeFollow))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", followevents.TopicFollow))
			}
		}
	}()

	cgComment, err := sarama.NewConsumerGroupFromClient("notification_comment", c.client)
	if err != nil {
		return err
	}
	go func() {
		for {
			err2 := cgComment.Consume(xcontext.Background(), []string{commentevents.TopicCommentCreate}, saramax.NewHandler[commentevents.CommentCreateEvent](c.l, c.consumeCommentCreate))
			if err2 != nil {
				c.l.Error("退出了消费循环异常", logger.Error(err2))
				time.Sleep(time.Second * 5)
			} else {
				c.l.Info("消费循环正常退出，可能是发生了Rebalance，正在重试", logger.String("topic", commentevents.TopicCommentCreate))
			}
		}
	}()
	return nil
}

func (c *NotificationEventConsumer) consumeLike(_ *sarama.ConsumerMessage, evt LikeEvent) error {
	ctx, cancel := xcontext.WithTimeout(xcontext.Background(), time.Second*3)
	defer cancel()
	recv := int64(0)
	if c.resolveOwner != nil {
		owner, err := c.resolveOwner(ctx, evt.Biz, evt.BizId)
		if err == nil {
			recv = owner
		}
	}
	if recv == 0 || recv == evt.Uid {
		return nil
	}
	n := dao.Notification{
		ReceiverId: recv,
		SenderId:   evt.Uid,
		Type:       "interaction",
		BizType:    evt.Biz,
		BizId:      evt.BizId,
		Content:    "点赞了你的内容",
		Status:     0,
		Ctime:      time.Now().Unix(),
		Utime:      time.Now().Unix(),
	}
	return c.notifDAO.Insert(ctx, n)
}

func (c *NotificationEventConsumer) consumeCollect(_ *sarama.ConsumerMessage, evt CollectEvent) error {
	ctx, cancel := xcontext.WithTimeout(xcontext.Background(), time.Second*3)
	defer cancel()
	recv := int64(0)
	if c.resolveOwner != nil {
		owner, err := c.resolveOwner(ctx, evt.Biz, evt.BizId)
		if err == nil {
			recv = owner
		}
	}
	if recv == 0 || recv == evt.Uid {
		return nil
	}
	n := dao.Notification{
		ReceiverId: recv,
		SenderId:   evt.Uid,
		Type:       "interaction",
		BizType:    evt.Biz,
		BizId:      evt.BizId,
		Content:    "收藏了你的内容",
		Status:     0,
		Ctime:      time.Now().Unix(),
		Utime:      time.Now().Unix(),
	}
	return c.notifDAO.Insert(ctx, n)
}

func (c *NotificationEventConsumer) consumeFollow(_ *sarama.ConsumerMessage, evt followevents.FollowEvent) error {
	ctx, cancel := xcontext.WithTimeout(xcontext.Background(), time.Second*3)
	defer cancel()
	if evt.Follower == evt.Followee {
		return nil
	}
	n := dao.Notification{
		ReceiverId: evt.Followee,
		SenderId:   evt.Follower,
		Type:       "follow",
		BizType:    "user",
		BizId:      evt.Follower,
		Content:    "关注了你",
		Status:     0,
		Ctime:      time.Now().Unix(),
		Utime:      time.Now().Unix(),
	}
	return c.notifDAO.Insert(ctx, n)
}

func (c *NotificationEventConsumer) consumeCommentCreate(_ *sarama.ConsumerMessage, evt commentevents.CommentCreateEvent) error {
	ctx, cancel := xcontext.WithTimeout(xcontext.Background(), time.Second*3)
	defer cancel()
	recv := int64(0)
	if c.resolveOwner != nil {
		owner, err := c.resolveOwner(ctx, evt.Biz, evt.BizId)
		if err == nil {
			recv = owner
		}
	}
	if recv == 0 || recv == evt.Uid {
		return nil
	}
	n := dao.Notification{
		ReceiverId: recv,
		SenderId:   evt.Uid,
		Type:       "interaction",
		BizType:    evt.Biz,
		BizId:      evt.BizId,
		Content:    "评论了你的内容",
		Status:     0,
		Ctime:      time.Now().Unix(),
		Utime:      time.Now().Unix(),
	}
	return c.notifDAO.Insert(ctx, n)
}
