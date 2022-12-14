package notice_queue

import (
	"answer/internal/schema"
)

var (
	NotificationQueue = make(chan *schema.NotificationMsg, 128)
)

func AddNotification(msg *schema.NotificationMsg) {
	NotificationQueue <- msg
}
