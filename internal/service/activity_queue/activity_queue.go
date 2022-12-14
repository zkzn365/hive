package activity_queue

import (
	"answer/internal/schema"
)

var (
	ActivityQueue = make(chan *schema.ActivityMsg, 128)
)

// AddActivity add new activity
func AddActivity(msg *schema.ActivityMsg) {
	ActivityQueue <- msg
}
