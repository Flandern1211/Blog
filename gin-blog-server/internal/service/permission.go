package service

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"
	g2 "gin-blog/pkg/errors"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionService interface {
	// Role
	GetRoleList(c *gin.Context, query request.PageQuery, keyword string) ([]response.RoleVO, int64, error)
	GetRoleOption(c *gin.Context) ([]response.OptionVO, error)
	SaveOrUpdateRole(c *gin.Context, req request.AddOrEditRoleReq) error
	DeleteRoles(c *gin.Context, ids []int) error

	// Menu
	GetUserMenu(c *gin.Context, authId int, isSuper bool) ([]response.MenuTreeVO, error)
	GetMenuTreeList(c *gin.Context, keyword string) ([]response.MenuTreeVO, error)
	GetMenuOption(c *gin.Context) ([]response.TreeOptionVO, error)
	SaveOrUpdateMenu(c *gin.Context, req request.AddOrEditMenuReq) error
	DeleteMenu(c *gin.Context, id int) error

	// Resource
	GetResourceTreeList(c *gin.Context, keyword string) ([]response.ResourceTreeVO, error)
	GetResourceOption(c *gin.Context) ([]response.TreeOptionVO, error)
	SaveOrUpdateResource(c *gin.Context, req request.AddOrEditResourceReq) error
	DeleteResource(c *gin.Context, id int) error
	UpdateResourceAnonymous(c *gin.Context, req request.EditAnonymousReq) error
}

type permissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(repo repository.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

// Role implementations
func (s *permissionService) GetRoleList(c *gin.Context, query request.PageQuery, keyword string) ([]response.RoleVO, int64, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetRoleList(db, query.GetPage(), query.GetSize(), keyword)
	if err != nil {
		return nil, 0, err
	}

	var res []response.RoleVO
	for _, role := range list {
		rVO := response.RoleVO{
			ID:        role.ID,
			Name:      role.Name,
			Label:     role.Label,
			IsDisable: role.IsDisable,
			CreatedAt: role.CreatedAt,
		}
		rVO.ResourceIds, _ = s.repo.GetResourceIdsByRoleId(db, role.ID)
		rVO.MenuIds, _ = s.repo.GetMenuIdsByRoleId(db, role.ID)
		res = append(res, rVO)
	}
	return res, total, nil
}

func (s *permissionService) GetRoleOption(c *gin.Context) ([]response.OptionVO, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	list, err := s.repo.GetRoleOption(db)
	if err != nil {
		return nil, err
	}
	var res []response.OptionVO
	for _, role := range list {
		res = append(res, response.OptionVO{ID: role.ID, Label: role.Label})
	}
	return res, nil
}

func (s *permissionService) SaveOrUpdateRole(c *gin.Context, req request.AddOrEditRoleReq) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	if req.ID == 0 {
		return s.repo.SaveRole(db, req.Name, req.Label)
	}
	return s.repo.UpdateRole(db, req.ID, req.Name, req.Label, req.IsDisable, req.ResourceIds, req.MenuIds)
}

func (s *permissionService) DeleteRoles(c *gin.Context, ids []int) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	return s.repo.DeleteRoles(db, ids)
}

// Menu implementations
func (s *permissionService) GetUserMenu(c *gin.Context, authId int, isSuper bool) ([]response.MenuTreeVO, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	var menus []entity.Menu
	var err error
	if isSuper {
		menus, err = s.repo.GetAllMenuList(db)
	} else {
		menus, err = s.repo.GetMenuListByUserId(db, authId)
	}
	if err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus, 0), nil
}

func (s *permissionService) GetMenuTreeList(c *gin.Context, keyword string) ([]response.MenuTreeVO, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	menus, err := s.repo.GetMenuList(db, keyword)
	if err != nil {
		return nil, err
	}
	return s.buildMenuTree(menus, 0), nil
}

func (s *permissionService) GetMenuOption(c *gin.Context) ([]response.TreeOptionVO, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	menus, err := s.repo.GetMenuList(db, "")
	if err != nil {
		return nil, err
	}
	return s.buildMenuTreeOption(menus, 0), nil
}

func (s *permissionService) buildMenuTreeOption(menus []entity.Menu, parentId int) []response.TreeOptionVO {
	var tree []response.TreeOptionVO
	for _, m := range menus {
		if m.ParentId == parentId {
			tree = append(tree, response.TreeOptionVO{
				ID:       m.ID,
				Label:    m.Name,
				Children: s.buildMenuTreeOption(menus, m.ID),
			})
		}
	}
	return tree
}

func (s *permissionService) buildMenuTree(menus []entity.Menu, parentId int) []response.MenuTreeVO {
	var tree []response.MenuTreeVO
	for _, m := range menus {
		if m.ParentId == parentId {
			tree = append(tree, response.MenuTreeVO{
				Menu:     m,
				Children: s.buildMenuTree(menus, m.ID),
			})
		}
	}
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].OrderNum < tree[j].OrderNum
	})
	return tree
}

func (s *permissionService) SaveOrUpdateMenu(c *gin.Context, req request.AddOrEditMenuReq) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	menu := &entity.Menu{
		Model:        entity.Model{ID: req.ID},
		ParentId:     req.ParentId,
		Name:         req.Name,
		Path:         req.Path,
		Component:    req.Component,
		Icon:         req.Icon,
		OrderNum:     req.OrderNum,
		Redirect:     req.Redirect,
		Catalogue:    req.Catalogue,
		Hidden:       req.Hidden,
		KeepAlive:    req.KeepAlive,
		External:     req.External,
		ExternalLink: req.ExternalLink,
	}
	return s.repo.SaveOrUpdateMenu(db, menu)
}

func (s *permissionService) DeleteMenu(c *gin.Context, id int) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	inUse, _ := s.repo.CheckMenuInUse(db, id)
	if inUse {
		return g2.NewDefault(g2.CodeMenuUsedByRole)
	}
	hasChild, _ := s.repo.CheckMenuHasChild(db, id)
	if hasChild {
		return g2.NewDefault(g2.CodeMenuHasChildren)
	}
	return s.repo.DeleteMenu(db, id)
}

// Resource implementations
func (s *permissionService) GetResourceTreeList(c *gin.Context, keyword string) ([]response.ResourceTreeVO, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	resources, err := s.repo.GetResourceList(db, keyword)
	if err != nil {
		return nil, err
	}
	return s.buildResourceTree(resources, 0), nil
}

func (s *permissionService) buildResourceTree(resources []entity.Resource, parentId int) []response.ResourceTreeVO {
	var tree []response.ResourceTreeVO
	for _, r := range resources {
		if r.ParentId == parentId {
			tree = append(tree, response.ResourceTreeVO{
				ID:        r.ID,
				CreatedAt: r.CreatedAt,
				Name:      r.Name,
				Url:       r.Url,
				Method:    r.Method,
				Anonymous: r.Anonymous,
				Children:  s.buildResourceTree(resources, r.ID),
			})
		}
	}
	return tree
}

func (s *permissionService) GetResourceOption(c *gin.Context) ([]response.TreeOptionVO, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	resources, err := s.repo.GetResourceList(db, "")
	if err != nil {
		return nil, err
	}
	return s.buildResourceOptionTree(resources, 0), nil
}

func (s *permissionService) buildResourceOptionTree(resources []entity.Resource, parentId int) []response.TreeOptionVO {
	var tree []response.TreeOptionVO
	for _, r := range resources {
		if r.ParentId == parentId {
			tree = append(tree, response.TreeOptionVO{
				ID:       r.ID,
				Label:    r.Name,
				Children: s.buildResourceOptionTree(resources, r.ID),
			})
		}
	}
	return tree
}

func (s *permissionService) SaveOrUpdateResource(c *gin.Context, req request.AddOrEditResourceReq) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	return s.repo.SaveOrUpdateResource(db, req.ID, req.ParentId, req.Name, req.Url, req.Method)
}

func (s *permissionService) DeleteResource(c *gin.Context, id int) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	return s.repo.DeleteResource(db, id)
}

func (s *permissionService) UpdateResourceAnonymous(c *gin.Context, req request.EditAnonymousReq) error {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	return s.repo.UpdateResourceAnonymous(db, req.ID, req.Anonymous)
}
