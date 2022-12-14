package service

import (
	"answer/internal/service/action"
	"answer/internal/service/activity"
	"answer/internal/service/activity_common"
	answercommon "answer/internal/service/answer_common"
	"answer/internal/service/auth"
	collectioncommon "answer/internal/service/collection_common"
	"answer/internal/service/comment"
	"answer/internal/service/comment_common"
	"answer/internal/service/dashboard"
	"answer/internal/service/export"
	"answer/internal/service/follow"
	"answer/internal/service/meta"
	"answer/internal/service/notification"
	notficationcommon "answer/internal/service/notification_common"
	"answer/internal/service/object_info"
	questioncommon "answer/internal/service/question_common"
	"answer/internal/service/rank"
	"answer/internal/service/reason"
	"answer/internal/service/report"
	"answer/internal/service/report_backyard"
	"answer/internal/service/report_handle_backyard"
	"answer/internal/service/revision_common"
	"answer/internal/service/search_parser"
	"answer/internal/service/siteinfo"
	"answer/internal/service/siteinfo_common"
	"answer/internal/service/tag"
	tagcommon "answer/internal/service/tag_common"
	"answer/internal/service/uploader"
	"answer/internal/service/user_backyard"
	usercommon "answer/internal/service/user_common"

	"github.com/google/wire"
)

// ProviderSetService is providers.
var ProviderSetService = wire.NewSet(
	comment.NewCommentService,
	comment_common.NewCommentCommonService,
	report.NewReportService,
	NewVoteService,
	tag.NewTagService,
	follow.NewFollowService,
	NewCollectionGroupService,
	NewCollectionService,
	action.NewCaptchaService,
	auth.NewAuthService,
	NewUserService,
	NewQuestionService,
	NewAnswerService,
	export.NewEmailService,
	tagcommon.NewTagCommonService,
	usercommon.NewUserCommon,
	questioncommon.NewQuestionCommon,
	answercommon.NewAnswerCommon,
	uploader.NewUploaderService,
	collectioncommon.NewCollectionCommon,
	revision_common.NewRevisionService,
	NewRevisionService,
	rank.NewRankService,
	search_parser.NewSearchParser,
	NewSearchService,
	meta.NewMetaService,
	object_info.NewObjService,
	report_handle_backyard.NewReportHandle,
	report_backyard.NewReportBackyardService,
	user_backyard.NewUserBackyardService,
	reason.NewReasonService,
	siteinfo_common.NewSiteInfoCommonService,
	siteinfo.NewSiteInfoService,
	notficationcommon.NewNotificationCommon,
	notification.NewNotificationService,
	activity.NewAnswerActivityService,
	dashboard.NewDashboardService,
	activity_common.NewActivityCommon,
	activity.NewActivityService,
)
