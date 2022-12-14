package notificationcommon

import (
	"context"
	"fmt"
	"time"

	"answer/internal/base/constant"
	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/schema"
	"answer/internal/service/activity_common"
	"answer/internal/service/notice_queue"
	"answer/internal/service/object_info"
	usercommon "answer/internal/service/user_common"

	"github.com/goccy/go-json"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

type NotificationRepo interface {
	AddNotification(ctx context.Context, notification *entity.Notification) (err error)
	GetNotificationPage(ctx context.Context, search *schema.NotificationSearch) ([]*entity.Notification, int64, error)
	ClearUnRead(ctx context.Context, userID string, notificationType int) (err error)
	ClearIDUnRead(ctx context.Context, userID string, id string) (err error)
	GetByUserIdObjectIdTypeId(ctx context.Context, userID, objectID string, notificationType int) (*entity.Notification, bool, error)
	UpdateNotificationContent(ctx context.Context, notification *entity.Notification) (err error)
	GetById(ctx context.Context, id string) (*entity.Notification, bool, error)
}

type NotificationCommon struct {
	data              *data.Data
	notificationRepo  NotificationRepo
	activityRepo      activity_common.ActivityRepo
	followRepo        activity_common.FollowRepo
	userCommon        *usercommon.UserCommon
	objectInfoService *object_info.ObjService
}

func NewNotificationCommon(
	data *data.Data,
	notificationRepo NotificationRepo,
	userCommon *usercommon.UserCommon,
	activityRepo activity_common.ActivityRepo,
	followRepo activity_common.FollowRepo,
	objectInfoService *object_info.ObjService,
) *NotificationCommon {
	notification := &NotificationCommon{
		data:              data,
		notificationRepo:  notificationRepo,
		activityRepo:      activityRepo,
		followRepo:        followRepo,
		userCommon:        userCommon,
		objectInfoService: objectInfoService,
	}
	notification.HandleNotification()
	return notification
}

func (ns *NotificationCommon) HandleNotification() {
	go func() {
		for msg := range notice_queue.NotificationQueue {
			log.Debugf("received notification %+v", msg)
			err := ns.AddNotification(context.TODO(), msg)
			if err != nil {
				log.Error(err)
			}
		}
	}()
}

// AddNotification
// need set
// UserID
// Type  1 inbox 2 achievement
// [inbox] Activity
// [achievement] Rank
// ObjectInfo.Title
// ObjectInfo.ObjectID
// ObjectInfo.ObjectType
func (ns *NotificationCommon) AddNotification(ctx context.Context, msg *schema.NotificationMsg) error {
	req := &schema.NotificationContent{
		TriggerUserID:  msg.TriggerUserID,
		ReceiverUserID: msg.ReceiverUserID,
		ObjectInfo: schema.ObjectInfo{
			Title:      msg.Title,
			ObjectID:   msg.ObjectID,
			ObjectType: msg.ObjectType,
		},
		NotificationAction: msg.NotificationAction,
		Type:               msg.Type,
	}
	var questionID string // just for notify all followers
	objInfo, err := ns.objectInfoService.GetInfo(ctx, req.ObjectInfo.ObjectID)
	if err != nil {
		log.Error(err)
	} else {
		req.ObjectInfo.Title = objInfo.Title
		questionID = objInfo.QuestionID
		objectMap := make(map[string]string)
		objectMap["question"] = objInfo.QuestionID
		objectMap["answer"] = objInfo.AnswerID
		objectMap["comment"] = objInfo.CommentID
		req.ObjectInfo.ObjectMap = objectMap
	}

	if msg.Type == schema.NotificationTypeAchievement {
		notificationInfo, exist, err := ns.notificationRepo.GetByUserIdObjectIdTypeId(ctx, req.ReceiverUserID, req.ObjectInfo.ObjectID, req.Type)
		if err != nil {
			return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
		}
		rank, err := ns.activityRepo.GetUserIDObjectIDActivitySum(ctx, req.ReceiverUserID, req.ObjectInfo.ObjectID)
		if err != nil {
			return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
		}
		req.Rank = rank
		if exist {
			//modify notification
			updateContent := &schema.NotificationContent{}
			err := json.Unmarshal([]byte(notificationInfo.Content), updateContent)
			if err != nil {
				return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
			}
			updateContent.Rank = rank
			content, err := json.Marshal(updateContent)
			if err != nil {
				return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
			}
			notificationInfo.Content = string(content)
			err = ns.notificationRepo.UpdateNotificationContent(ctx, notificationInfo)
			if err != nil {
				return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
			}
			return nil
		}
	}

	info := &entity.Notification{}
	now := time.Now()
	info.UserID = req.ReceiverUserID
	info.Type = req.Type
	info.IsRead = schema.NotificationNotRead
	info.Status = schema.NotificationStatusNormal
	info.CreatedAt = now
	info.UpdatedAt = now
	info.ObjectID = req.ObjectInfo.ObjectID

	userBasicInfo, exist, err := ns.userCommon.GetUserBasicInfoByID(ctx, req.TriggerUserID)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	if !exist {
		return errors.InternalServer(reason.UserNotFound).WithError(err).WithStack()
	}
	req.UserInfo = userBasicInfo
	content, err := json.Marshal(req)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	info.Content = string(content)
	err = ns.notificationRepo.AddNotification(ctx, info)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	err = ns.addRedDot(ctx, info.UserID, info.Type)
	if err != nil {
		log.Error("addRedDot Error", err.Error())
	}

	go ns.SendNotificationToAllFollower(context.Background(), msg, questionID)
	return nil
}

func (ns *NotificationCommon) addRedDot(ctx context.Context, userID string, botType int) error {
	key := fmt.Sprintf("answer_RedDot_%d_%s", botType, userID)
	err := ns.data.Cache.SetInt64(ctx, key, 1, 30*24*time.Hour) //Expiration time is one month.
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return nil
}

// SendNotificationToAllFollower send notification to all followers
func (ns *NotificationCommon) SendNotificationToAllFollower(ctx context.Context, msg *schema.NotificationMsg,
	questionID string) {
	if msg.NoNeedPushAllFollow {
		return
	}
	if msg.NotificationAction != constant.UpdateQuestion &&
		msg.NotificationAction != constant.AnswerTheQuestion &&
		msg.NotificationAction != constant.UpdateAnswer &&
		msg.NotificationAction != constant.AdoptAnswer {
		return
	}
	condObjectID := msg.ObjectID
	if len(questionID) > 0 {
		condObjectID = questionID
	}
	userIDs, err := ns.followRepo.GetFollowUserIDs(ctx, condObjectID)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("send notification to all followers: %s %d", condObjectID, len(userIDs))
	for _, userID := range userIDs {
		t := &schema.NotificationMsg{}
		_ = copier.Copy(t, msg)
		t.ReceiverUserID = userID
		t.TriggerUserID = msg.TriggerUserID
		t.NoNeedPushAllFollow = true
		notice_queue.AddNotification(t)
	}
}
