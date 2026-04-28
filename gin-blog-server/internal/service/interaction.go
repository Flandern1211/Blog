package service

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"
	"gin-blog/internal/utils"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InteractionService interface {
	// Message
	GetMessageList(c *gin.Context, query request.MessageQuery) ([]entity.Message, int64, error)
	DeleteMessages(c *gin.Context, ids []int) error
	UpdateMessagesReview(c *gin.Context, req request.UpdateReviewReq) error

	// Comment
	GetCommentList(c *gin.Context, query request.CommentQuery) ([]entity.Comment, int64, error)
	DeleteComments(c *gin.Context, ids []int) error
	UpdateCommentsReview(c *gin.Context, req request.UpdateReviewReq) error

	// Front
	AddMessage(c *gin.Context, authId int, req request.FAddMessageReq) error
	GetFrontCommentList(c *gin.Context, query request.FCommentQuery) ([]response.CommentVO, int64, error)
	GetCommentReplyList(c *gin.Context, id int, page, size int) ([]response.CommentVO, error)
	AddComment(c *gin.Context, authId int, req request.FAddCommentReq) error
}

type interactionService struct {
	repo         repository.InteractionRepository
	blogInfoRepo repository.BlogInfoRepository
}

func NewInteractionService(repo repository.InteractionRepository, blogInfoRepo repository.BlogInfoRepository) InteractionService {
	return &interactionService{
		repo:         repo,
		blogInfoRepo: blogInfoRepo,
	}
}

// Message implementations
func (s *interactionService) GetMessageList(c *gin.Context, query request.MessageQuery) ([]entity.Message, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetMessageList(db, query.GetPage(), query.GetSize(), query.Nickname, query.IsReview)
}

func (s *interactionService) DeleteMessages(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeleteMessages(db, ids)
}

func (s *interactionService) UpdateMessagesReview(c *gin.Context, req request.UpdateReviewReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateMessagesReview(db, req.Ids, req.IsReview)
}

// Comment implementations
func (s *interactionService) GetCommentList(c *gin.Context, query request.CommentQuery) ([]entity.Comment, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetCommentList(db, query.GetPage(), query.GetSize(), query.Type, query.IsReview, query.Nickname)
}

func (s *interactionService) DeleteComments(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeleteComments(db, ids)
}

func (s *interactionService) UpdateCommentsReview(c *gin.Context, req request.UpdateReviewReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateCommentsReview(db, req.Ids, req.IsReview)
}

// Front implementations

func (s *interactionService) AddMessage(c *gin.Context, authId int, req request.FAddMessageReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)

	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSource(ipAddress)
	isReview := s.blogInfoRepo.GetConfigBool(db, global.CONFIG_IS_COMMENT_REVIEW)

	// 获取用户信息
	authRepo := repository.NewAuthRepository()
	user, err := authRepo.GetUserAuthInfoById(db, authId)
	if err != nil {
		return err
	}

	message := &entity.Message{
		Nickname:  user.UserInfo.Nickname,
		Avatar:    user.UserInfo.Avatar,
		Content:   template.HTMLEscapeString(req.Content),
		IpAddress: ipAddress,
		IpSource:  ipSource,
		Speed:     req.Speed,
		IsReview:  isReview,
	}

	return s.repo.SaveMessage(db, message)
}

func (s *interactionService) GetFrontCommentList(c *gin.Context, query request.FCommentQuery) ([]response.CommentVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	comments, replyMap, total, err := s.repo.GetFrontCommentList(db, query.GetPage(), query.GetSize(), query.TopicId, query.Type)
	if err != nil {
		return nil, 0, err
	}

	// 获取点赞数据
	likeCountMap := rdb.HGetAll(rctx, global.COMMENT_LIKE_COUNT).Val()

	var res []response.CommentVO
	for _, comment := range comments {
		vo := response.CommentVO{
			Comment:    comment,
			LikeCount:  0,
			ReplyCount: len(replyMap[comment.ID]),
		}
		if count, ok := likeCountMap[strconv.Itoa(comment.ID)]; ok {
			vo.LikeCount, _ = strconv.Atoi(count)
		}

		// 处理回复
		var replies []response.CommentVO
		for _, reply := range replyMap[comment.ID] {
			replyVO := response.CommentVO{
				Comment:   reply,
				LikeCount: 0,
			}
			if count, ok := likeCountMap[strconv.Itoa(reply.ID)]; ok {
				replyVO.LikeCount, _ = strconv.Atoi(count)
			}
			replies = append(replies, replyVO)
		}
		vo.ReplyList = replies
		res = append(res, vo)
	}

	return res, total, nil
}

func (s *interactionService) GetCommentReplyList(c *gin.Context, id int, page, size int) ([]response.CommentVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	replies, err := s.repo.GetCommentReplyList(db, id, page, size)
	if err != nil {
		return nil, err
	}

	likeCountMap := rdb.HGetAll(rctx, global.COMMENT_LIKE_COUNT).Val()

	var res []response.CommentVO
	for _, reply := range replies {
		vo := response.CommentVO{
			Comment:   reply,
			LikeCount: 0,
		}
		if count, ok := likeCountMap[strconv.Itoa(reply.ID)]; ok {
			vo.LikeCount, _ = strconv.Atoi(count)
		}
		res = append(res, vo)
	}

	return res, nil
}

func (s *interactionService) AddComment(c *gin.Context, authId int, req request.FAddCommentReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	isReview := s.blogInfoRepo.GetConfigBool(db, global.CONFIG_IS_COMMENT_REVIEW)

	comment := &entity.Comment{
		UserId:      authId,
		ReplyUserId: req.ReplyUserId,
		TopicId:     req.TopicId,
		ParentId:    req.ParentId,
		Content:     template.HTMLEscapeString(req.Content),
		Type:        req.Type,
		IsReview:    isReview,
	}

	return s.repo.AddComment(db, comment)
}
