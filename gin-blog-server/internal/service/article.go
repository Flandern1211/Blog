package service

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"
	"strconv"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	bizErr "gin-blog/pkg/errors"
)

type ArticleService interface {
	// Article
	GetList(c *gin.Context, query request.ArticleQuery) ([]response.ArticleVO, int64, error)
	GetById(c *gin.Context, id int) (*response.ArticleVO, error)
	SaveOrUpdate(c *gin.Context, authId int, req request.AddOrEditArticleReq) error
	UpdateTop(c *gin.Context, req request.UpdateArticleTopReq) error
	SoftDelete(c *gin.Context, req request.SoftDeleteReq) error
	Delete(c *gin.Context, ids []int) error

	// Front-end specific Article methods
	GetBlogArticleList(c *gin.Context, query request.FArticleQuery) ([]response.ArticleVO, int64, error)
	GetBlogArticle(c *gin.Context, id int) (*response.BlogArticleVO, error)

	// Category
	GetCategoryList(c *gin.Context, query request.CategoryQuery) ([]response.CategoryVO, int64, error)
	SaveOrUpdateCategory(c *gin.Context, req request.AddOrEditCategoryReq) error
	DeleteCategories(c *gin.Context, ids []int) error
	GetCategoryOption(c *gin.Context) ([]response.OptionVO, error)

	// Tag
	GetTagList(c *gin.Context, query request.TagQuery) ([]response.TagVO, int64, error)
	SaveOrUpdateTag(c *gin.Context, req request.AddOrEditTagReq) error
	DeleteTags(c *gin.Context, ids []int) error
	GetTagOption(c *gin.Context) ([]response.OptionVO, error)
}

type articleService struct {
	repo         repository.ArticleRepository
	interactRepo repository.InteractionRepository
}

func NewArticleService(repo repository.ArticleRepository, interactRepo repository.InteractionRepository) ArticleService {
	return &articleService{
		repo:         repo,
		interactRepo: interactRepo,
	}
}

// Article implementations
func (s *articleService) GetList(c *gin.Context, query request.ArticleQuery) ([]response.ArticleVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetList(db, query.GetPage(), query.GetSize(), query.Title, query.CategoryId, query.TagId, query.Type, query.Status, query.IsDelete)
	if err != nil {
		return nil, 0, err
	}

	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	var res []response.ArticleVO
	for _, art := range list {
		vo := response.ArticleVO{Article: art}
		// Get stats from Redis
		likeCount, _ := rdb.HGet(rctx, global.ARTICLE_LIKE_COUNT, strconv.Itoa(art.ID)).Int()
		viewCount, _ := rdb.ZScore(rctx, global.ARTICLE_VIEW_COUNT, strconv.Itoa(art.ID)).Result()
		vo.LikeCount = likeCount
		vo.ViewCount = int(viewCount)
		res = append(res, vo)
	}
	return res, total, nil
}

func (s *articleService) GetById(c *gin.Context, id int) (*response.ArticleVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	art, err := s.repo.GetById(db, id)
	if err != nil {
		return nil, err
	}
	vo := &response.ArticleVO{Article: *art}
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()
	likeCount, _ := rdb.HGet(rctx, global.ARTICLE_LIKE_COUNT, strconv.Itoa(art.ID)).Int()
	viewCount, _ := rdb.ZScore(rctx, global.ARTICLE_VIEW_COUNT, strconv.Itoa(art.ID)).Result()
	vo.LikeCount = likeCount
	vo.ViewCount = int(viewCount)
	return vo, nil
}

func (s *articleService) SaveOrUpdate(c *gin.Context, authId int, req request.AddOrEditArticleReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	article := &entity.Article{
		Model:       entity.Model{ID: req.ID},
		Title:       req.Title,
		Desc:        req.Desc,
		Content:     req.Content,
		Img:         req.Img,
		Type:        req.Type,
		Status:      req.Status,
		IsTop:       req.IsTop,
		OriginalUrl: req.OriginalUrl,
		UserId:      authId,
	}
	return s.repo.SaveOrUpdate(db, article, req.CategoryName, req.TagNames)
}

func (s *articleService) UpdateTop(c *gin.Context, req request.UpdateArticleTopReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateTop(db, req.ID, req.IsTop)
}

func (s *articleService) SoftDelete(c *gin.Context, req request.SoftDeleteReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.SoftDelete(db, req.Ids, req.IsDelete)
}

func (s *articleService) Delete(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.Delete(db, ids)
}

func (s *articleService) entityToVO(article entity.Article) response.ArticleVO {
	return response.ArticleVO{
		Article: article,
	}
}

func (s *articleService) GetBlogArticleList(c *gin.Context, query request.FArticleQuery) ([]response.ArticleVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetBlogArticleList(db, query.GetPage(), query.GetSize(), query.CategoryId, query.TagId)
	if err != nil {
		return nil, 0, err
	}
	var voList []response.ArticleVO
	for _, article := range list {
		voList = append(voList, s.entityToVO(article))
	}
	return voList, total, nil
}

func (s *articleService) GetBlogArticle(c *gin.Context, id int) (*response.BlogArticleVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	article, err := s.repo.GetBlogArticle(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizErr.ErrNotFound
		}
		return nil, err
	}

	vo := &response.BlogArticleVO{
		Article:           *article,
		RecommendArticles: make([]response.RecommendArticleVO, 0),
		NewestArticles:    make([]response.RecommendArticleVO, 0),
	}

	recommendList, _ := s.repo.GetRecommendList(db, id, 6)
	for _, v := range recommendList {
		vo.RecommendArticles = append(vo.RecommendArticles, response.RecommendArticleVO{
			ID:        v.ID,
			Img:       v.Img,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
		})
	}

	newestList, _ := s.repo.GetNewestList(db, 5)
	for _, v := range newestList {
		vo.NewestArticles = append(vo.NewestArticles, response.RecommendArticleVO{
			ID:        v.ID,
			Img:       v.Img,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
		})
	}

	lastArt, _ := s.repo.GetLastArticle(db, id)
	vo.LastArticle = response.ArticlePaginationVO{
		ID:    lastArt.ID,
		Img:   lastArt.Img,
		Title: lastArt.Title,
	}

	nextArt, _ := s.repo.GetNextArticle(db, id)
	vo.NextArticle = response.ArticlePaginationVO{
		ID:    nextArt.ID,
		Img:   nextArt.Img,
		Title: nextArt.Title,
	}

	rdb.ZIncrBy(rctx, global.ARTICLE_VIEW_COUNT, 1, strconv.Itoa(id))

	vo.ViewCount = int64(rdb.ZScore(rctx, global.ARTICLE_VIEW_COUNT, strconv.Itoa(id)).Val())
	likeCount, _ := strconv.Atoi(rdb.HGet(rctx, global.ARTICLE_LIKE_COUNT, strconv.Itoa(id)).Val())
	vo.LikeCount = int64(likeCount)
	vo.CommentCount, _ = s.interactRepo.GetArticleCommentCount(db, id)

	return vo, nil
}

// Category implementations
func (s *articleService) GetCategoryList(c *gin.Context, query request.CategoryQuery) ([]response.CategoryVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetCategoryList(db, query.GetPage(), query.GetSize(), query.Keyword)
	if err != nil {
		return nil, 0, err
	}
	var res []response.CategoryVO
	for _, cat := range list {
		res = append(res, response.CategoryVO{
			Category:     cat.Category,
			ArticleCount: cat.ArticleCount,
		})
	}
	return res, total, nil
}

func (s *articleService) SaveOrUpdateCategory(c *gin.Context, req request.AddOrEditCategoryReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.SaveOrUpdateCategory(db, req.ID, req.Name)
}

func (s *articleService) DeleteCategories(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeleteCategories(db, ids)
}

func (s *articleService) GetCategoryOption(c *gin.Context) ([]response.OptionVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, err := s.repo.GetCategoryOption(db)
	if err != nil {
		return nil, err
	}
	var res []response.OptionVO
	for _, cat := range list {
		res = append(res, response.OptionVO{ID: cat.ID, Label: cat.Name})
	}
	return res, nil
}

// Tag implementations
func (s *articleService) GetTagList(c *gin.Context, query request.TagQuery) ([]response.TagVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetTagList(db, query.GetPage(), query.GetSize(), query.Keyword)
	if err != nil {
		return nil, 0, err
	}
	var res []response.TagVO
	for _, tag := range list {
		res = append(res, response.TagVO{
			Tag:          tag.Tag,
			ArticleCount: tag.ArticleCount,
		})
	}
	return res, total, nil
}

func (s *articleService) SaveOrUpdateTag(c *gin.Context, req request.AddOrEditTagReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.SaveOrUpdateTag(db, req.ID, req.Name)
}

func (s *articleService) DeleteTags(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeleteTags(db, ids)
}

func (s *articleService) GetTagOption(c *gin.Context) ([]response.OptionVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, err := s.repo.GetTagOption(db)
	if err != nil {
		return nil, err
	}
	var res []response.OptionVO
	for _, tag := range list {
		res = append(res, response.OptionVO{ID: tag.ID, Label: tag.Name})
	}
	return res, nil
}
