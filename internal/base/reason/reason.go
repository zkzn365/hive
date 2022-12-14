package reason

const (
	// Success .
	Success = "base.success"
	// UnknownError unknown error
	UnknownError = "base.unknown"
	// RequestFormatError request format error
	RequestFormatError = "base.request_format_error"
	// UnauthorizedError unauthorized error
	UnauthorizedError = "base.unauthorized_error"
	// DatabaseError database error
	DatabaseError = "base.database_error"
)

const (
	EmailOrPasswordWrong             = "error.object.email_or_password_incorrect"
	CommentNotFound                  = "error.comment.not_found"
	QuestionNotFound                 = "error.question.not_found"
	QuestionCannotDeleted            = "error.question.cannot_deleted"
	QuestionCannotClose              = "error.question.cannot_close"
	QuestionCannotUpdate             = "error.question.cannot_update"
	AnswerNotFound                   = "error.answer.not_found"
	AnswerCannotDeleted              = "error.answer.cannot_deleted"
	AnswerCannotUpdate               = "error.answer.cannot_update"
	CommentEditWithoutPermission     = "error.comment.edit_without_permission"
	DisallowVote                     = "error.object.disallow_vote"
	DisallowFollow                   = "error.object.disallow_follow"
	DisallowVoteYourSelf             = "error.object.disallow_vote_your_self"
	CaptchaVerificationFailed        = "error.object.captcha_verification_failed"
	OldPasswordVerificationFailed    = "error.object.old_password_verification_failed"
	NewPasswordSameAsPreviousSetting = "error.object.new_password_same_as_previous_setting"
	UserNotFound                     = "error.user.not_found"
	UsernameInvalid                  = "error.user.username_invalid"
	UsernameDuplicate                = "error.user.username_duplicate"
	UserSetAvatar                    = "error.user.set_avatar"
	EmailDuplicate                   = "error.email.duplicate"
	EmailVerifyURLExpired            = "error.email.verify_url_expired"
	EmailNeedToBeVerified            = "error.email.need_to_be_verified"
	UserSuspended                    = "error.user.suspended"
	ObjectNotFound                   = "error.object.not_found"
	TagNotFound                      = "error.tag.not_found"
	TagNotContainSynonym             = "error.tag.not_contain_synonym_tags"
	TagCannotUpdate                  = "error.tag.cannot_update"
	RankFailToMeetTheCondition       = "error.rank.fail_to_meet_the_condition"
	ThemeNotFound                    = "error.theme.not_found"
	LangNotFound                     = "error.lang.not_found"
	ReportHandleFailed               = "error.report.handle_failed"
	ReportNotFound                   = "error.report.not_found"
	ReadConfigFailed                 = "error.config.read_config_failed"
	DatabaseConnectionFailed         = "error.database.connection_failed"
	InstallCreateTableFailed         = "error.database.create_table_failed"
	InstallConfigFailed              = "error.install.create_config_failed"
	SiteInfoNotFound                 = "error.site_info.not_found"
	UploadFileSourceUnsupported      = "error.upload.source_unsupported"
	RecommendTagNotExist             = "error.tag.recommend_tag_not_found"
	RecommendTagEnter                = "error.tag.recommend_tag_enter"
	RevisionReviewUnderway           = "error.revision.review_underway"
	RevisionNoPermission             = "error.revision.no_permission"
)
