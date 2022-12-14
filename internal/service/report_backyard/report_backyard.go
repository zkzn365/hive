package report_backyard

import (
	"answer/internal/service/config"
	"answer/pkg/htmltext"
	"context"

	"answer/internal/base/pager"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/repo/common"
	"answer/internal/schema"
	answercommon "answer/internal/service/answer_common"
	"answer/internal/service/comment_common"
	questioncommon "answer/internal/service/question_common"
	"answer/internal/service/report_common"
	"answer/internal/service/report_handle_backyard"
	usercommon "answer/internal/service/user_common"

	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
)

// ReportBackyardService user service
type ReportBackyardService struct {
	reportRepo        report_common.ReportRepo
	commonUser        *usercommon.UserCommon
	commonRepo        *common.CommonRepo
	answerRepo        answercommon.AnswerRepo
	questionRepo      questioncommon.QuestionRepo
	commentCommonRepo comment_common.CommentCommonRepo
	reportHandle      *report_handle_backyard.ReportHandle
	configRepo        config.ConfigRepo
}

// NewReportBackyardService new report service
func NewReportBackyardService(
	reportRepo report_common.ReportRepo,
	commonUser *usercommon.UserCommon,
	commonRepo *common.CommonRepo,
	answerRepo answercommon.AnswerRepo,
	questionRepo questioncommon.QuestionRepo,
	commentCommonRepo comment_common.CommentCommonRepo,
	reportHandle *report_handle_backyard.ReportHandle,
	configRepo config.ConfigRepo) *ReportBackyardService {
	return &ReportBackyardService{
		reportRepo:        reportRepo,
		commonUser:        commonUser,
		commonRepo:        commonRepo,
		answerRepo:        answerRepo,
		questionRepo:      questionRepo,
		commentCommonRepo: commentCommonRepo,
		reportHandle:      reportHandle,
		configRepo:        configRepo,
	}
}

// ListReportPage list report pages
func (rs *ReportBackyardService) ListReportPage(ctx context.Context, dto schema.GetReportListPageDTO) (pageModel *pager.PageModel, err error) {
	var (
		resp  []*schema.GetReportListPageResp
		flags []entity.Report
		total int64

		flaggedUserIds,
		userIds []string

		flaggedUsers,
		users map[string]*schema.UserBasicInfo
	)

	pageModel = &pager.PageModel{}

	flags, total, err = rs.reportRepo.GetReportListPage(ctx, dto)
	if err != nil {
		return
	}

	_ = copier.Copy(&resp, flags)
	for _, r := range resp {
		flaggedUserIds = append(flaggedUserIds, r.ReportedUserID)
		userIds = append(userIds, r.UserID)
		r.Format()
	}

	// flagged users
	flaggedUsers, err = rs.commonUser.BatchUserBasicInfoByID(ctx, flaggedUserIds)

	// flag users
	users, err = rs.commonUser.BatchUserBasicInfoByID(ctx, userIds)
	for _, r := range resp {
		r.ReportedUser = flaggedUsers[r.ReportedUserID]
		r.ReportUser = users[r.UserID]
	}

	rs.parseObject(ctx, &resp)
	return pager.NewPageModel(total, resp), nil
}

// HandleReported handle the reported object
func (rs *ReportBackyardService) HandleReported(ctx context.Context, req schema.ReportHandleReq) (err error) {
	var (
		reported   = entity.Report{}
		handleData = entity.Report{
			FlaggedContent: req.FlaggedContent,
			FlaggedType:    req.FlaggedType,
			Status:         entity.ReportStatusCompleted,
		}
		exist = false
	)

	reported, exist, err = rs.reportRepo.GetByID(ctx, req.ID)
	if err != nil {
		err = errors.BadRequest(reason.ReportHandleFailed).WithError(err).WithStack()
		return
	}
	if !exist {
		err = errors.NotFound(reason.ReportNotFound)
		return
	}

	// check if handle or not
	if reported.Status != entity.ReportStatusPending {
		return
	}

	if err = rs.reportHandle.HandleObject(ctx, reported, req); err != nil {
		return
	}

	err = rs.reportRepo.UpdateByID(ctx, reported.ID, handleData)
	return
}

func (rs *ReportBackyardService) parseObject(ctx context.Context, resp *[]*schema.GetReportListPageResp) {
	var (
		res = *resp
	)

	for i, r := range res {
		var (
			objIds map[string]string
			exists,
			ok bool
			err error
			questionId,
			answerId,
			commentId string
			question *entity.Question
			answer   *entity.Answer
			cmt      *entity.Comment
		)

		objIds, err = rs.commonRepo.GetObjectIDMap(r.ObjectID)
		if err != nil {
			continue
		}

		questionId, ok = objIds["question"]
		if !ok {
			continue
		}

		question, exists, err = rs.questionRepo.GetQuestion(ctx, questionId)
		if err != nil || !exists {
			continue
		}

		answerId, ok = objIds["answer"]
		if ok {
			answer, _, err = rs.answerRepo.GetAnswer(ctx, answerId)
		}

		commentId, ok = objIds["comment"]
		if ok {
			cmt, _, err = rs.commentCommonRepo.GetComment(ctx, commentId)
		}

		switch r.OType {
		case "question":
			r.QuestionID = questionId
			r.Title = question.Title
			r.Excerpt = htmltext.FetchExcerpt(question.ParsedText, "...", 240)

		case "answer":
			r.QuestionID = questionId
			r.AnswerID = answerId
			r.Title = question.Title
			r.Excerpt = htmltext.FetchExcerpt(answer.ParsedText, "...", 240)

		case "comment":
			r.QuestionID = questionId
			r.AnswerID = answerId
			r.CommentID = commentId
			r.Title = question.Title
			r.Excerpt = htmltext.FetchExcerpt(cmt.ParsedText, "...", 240)
		}

		// parse reason
		if r.ReportType > 0 {
			r.Reason = &schema.ReasonItem{
				ReasonType: r.ReportType,
			}
			err = rs.configRepo.GetJsonConfigByIDAndSetToObject(r.ReportType, r.Reason)
		}
		if r.FlaggedType > 0 {
			r.FlaggedReason = &schema.ReasonItem{
				ReasonType: r.FlaggedType,
			}
			_ = rs.configRepo.GetJsonConfigByIDAndSetToObject(r.FlaggedType, r.FlaggedReason)
		}

		res[i] = r
	}
	resp = &res
}
