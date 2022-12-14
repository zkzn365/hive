package repo

import (
	"answer/internal/base/data"
	"answer/internal/repo/activity"
	"answer/internal/repo/activity_common"
	"answer/internal/repo/answer"
	"answer/internal/repo/auth"
	"answer/internal/repo/captcha"
	"answer/internal/repo/collection"
	"answer/internal/repo/comment"
	"answer/internal/repo/common"
	"answer/internal/repo/config"
	"answer/internal/repo/export"
	"answer/internal/repo/meta"
	"answer/internal/repo/notification"
	"answer/internal/repo/question"
	"answer/internal/repo/rank"
	"answer/internal/repo/reason"
	"answer/internal/repo/report"
	"answer/internal/repo/revision"
	"answer/internal/repo/search_common"
	"answer/internal/repo/site_info"
	"answer/internal/repo/tag"
	"answer/internal/repo/tag_common"
	"answer/internal/repo/unique"
	"answer/internal/repo/user"

	"github.com/google/wire"
)

// ProviderSetRepo is data providers.
var ProviderSetRepo = wire.NewSet(
	common.NewCommonRepo,
	data.NewData,
	data.NewDB,
	data.NewCache,
	comment.NewCommentRepo,
	comment.NewCommentCommonRepo,
	captcha.NewCaptchaRepo,
	unique.NewUniqueIDRepo,
	report.NewReportRepo,
	activity_common.NewFollowRepo,
	activity_common.NewVoteRepo,
	config.NewConfigRepo,
	user.NewUserRepo,
	user.NewUserBackyardRepo,
	rank.NewUserRankRepo,
	question.NewQuestionRepo,
	answer.NewAnswerRepo,
	activity_common.NewActivityRepo,
	activity.NewVoteRepo,
	activity.NewFollowRepo,
	activity.NewAnswerActivityRepo,
	activity.NewQuestionActivityRepo,
	activity.NewUserActiveActivityRepo,
	activity.NewActivityRepo,
	tag.NewTagRepo,
	tag_common.NewTagCommonRepo,
	tag.NewTagRelRepo,
	collection.NewCollectionRepo,
	collection.NewCollectionGroupRepo,
	auth.NewAuthRepo,
	revision.NewRevisionRepo,
	search_common.NewSearchRepo,
	meta.NewMetaRepo,
	export.NewEmailRepo,
	reason.NewReasonRepo,
	site_info.NewSiteInfo,
	notification.NewNotificationRepo,
)
