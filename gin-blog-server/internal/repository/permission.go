package repository

import (
	"gin-blog/internal/model/entity"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	// Role
	GetRoleList(db *gorm.DB, page, size int, keyword string) ([]entity.Role, int64, error)
	GetRoleOption(db *gorm.DB) ([]entity.Role, error)
	SaveRole(db *gorm.DB, name, label string) error
	UpdateRole(db *gorm.DB, id int, name, label string, isDisable bool, resourceIds, menuIds []int) error
	DeleteRoles(db *gorm.DB, ids []int) error
	GetResourceIdsByRoleId(db *gorm.DB, roleId int) ([]int, error)
	GetMenuIdsByRoleId(db *gorm.DB, roleId int) ([]int, error)

	// Menu
	GetMenuList(db *gorm.DB, keyword string) ([]entity.Menu, error)
	GetMenuListByUserId(db *gorm.DB, userId int) ([]entity.Menu, error)
	GetAllMenuList(db *gorm.DB) ([]entity.Menu, error)
	SaveOrUpdateMenu(db *gorm.DB, menu *entity.Menu) error
	DeleteMenu(db *gorm.DB, id int) error
	GetMenuById(db *gorm.DB, id int) (*entity.Menu, error)
	CheckMenuInUse(db *gorm.DB, id int) (bool, error)
	CheckMenuHasChild(db *gorm.DB, id int) (bool, error)

	// Resource
	GetResourceList(db *gorm.DB, keyword string) ([]entity.Resource, error)
	SaveOrUpdateResource(db *gorm.DB, id, parentId int, name, url, method string) error
	DeleteResource(db *gorm.DB, id int) error
	UpdateResourceAnonymous(db *gorm.DB, id int, anonymous bool) error
}

type permissionRepository struct{}

func NewPermissionRepository() PermissionRepository {
	return &permissionRepository{}
}

// Role implementations
func (r *permissionRepository) GetRoleList(db *gorm.DB, page, size int, keyword string) ([]entity.Role, int64, error) {
	var list []entity.Role
	var total int64
	db = db.Model(&entity.Role{})
	if keyword != "" {
		db = db.Where("name LIKE ? OR label LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := db.Count(&total).Scopes(Paginate(page, size)).Find(&list).Error
	return list, total, err
}

func (r *permissionRepository) GetRoleOption(db *gorm.DB) ([]entity.Role, error) {
	var list []entity.Role
	err := db.Model(&entity.Role{}).Select("id", "label").Find(&list).Error
	return list, err
}

func (r *permissionRepository) SaveRole(db *gorm.DB, name, label string) error {
	return db.Create(&entity.Role{Name: name, Label: label}).Error
}

func (r *permissionRepository) UpdateRole(db *gorm.DB, id int, name, label string, isDisable bool, resourceIds, menuIds []int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Role{Model: entity.Model{ID: id}}).Updates(entity.Role{
			Name:      name,
			Label:     label,
			IsDisable: isDisable,
		}).Error; err != nil {
			return err
		}

		// Update role_resource
		if err := tx.Delete(&entity.RoleResource{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		if len(resourceIds) > 0 {
			var roleResources []entity.RoleResource
			for _, rid := range resourceIds {
				roleResources = append(roleResources, entity.RoleResource{RoleId: id, ResourceId: rid})
			}
			if err := tx.Create(&roleResources).Error; err != nil {
				return err
			}
		}

		// Update role_menu
		if err := tx.Delete(&entity.RoleMenu{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		if len(menuIds) > 0 {
			var roleMenus []entity.RoleMenu
			for _, mid := range menuIds {
				roleMenus = append(roleMenus, entity.RoleMenu{RoleId: id, MenuId: mid})
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *permissionRepository) DeleteRoles(db *gorm.DB, ids []int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&entity.Role{}, ids).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.RoleResource{}, "role_id IN ?", ids).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.RoleMenu{}, "role_id IN ?", ids).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entity.UserAuthRole{}, "role_id IN ?", ids).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *permissionRepository) GetResourceIdsByRoleId(db *gorm.DB, roleId int) ([]int, error) {
	var ids []int
	err := db.Model(&entity.RoleResource{}).Where("role_id = ?", roleId).Pluck("resource_id", &ids).Error
	return ids, err
}

func (r *permissionRepository) GetMenuIdsByRoleId(db *gorm.DB, roleId int) ([]int, error) {
	var ids []int
	err := db.Model(&entity.RoleMenu{}).Where("role_id = ?", roleId).Pluck("menu_id", &ids).Error
	return ids, err
}

// Menu implementations
func (r *permissionRepository) GetMenuList(db *gorm.DB, keyword string) ([]entity.Menu, error) {
	var list []entity.Menu
	db = db.Model(&entity.Menu{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := db.Order("order_num").Find(&list).Error
	return list, err
}

func (r *permissionRepository) GetMenuListByUserId(db *gorm.DB, userId int) ([]entity.Menu, error) {
	var list []entity.Menu
	err := db.Table("menu").
		Joins("JOIN role_menu ON menu.id = role_menu.menu_id").
		Joins("JOIN user_auth_role ON role_menu.role_id = user_auth_role.role_id").
		Where("user_auth_role.user_auth_id = ?", userId).
		Distinct("menu.*").
		Order("order_num").
		Find(&list).Error
	return list, err
}

func (r *permissionRepository) GetAllMenuList(db *gorm.DB) ([]entity.Menu, error) {
	var list []entity.Menu
	err := db.Model(&entity.Menu{}).Order("order_num").Find(&list).Error
	return list, err
}

func (r *permissionRepository) SaveOrUpdateMenu(db *gorm.DB, menu *entity.Menu) error {
	if menu.ID == 0 {
		return db.Create(menu).Error
	}
	return db.Model(menu).Updates(menu).Error
}

func (r *permissionRepository) DeleteMenu(db *gorm.DB, id int) error {
	return db.Delete(&entity.Menu{}, id).Error
}

func (r *permissionRepository) GetMenuById(db *gorm.DB, id int) (*entity.Menu, error) {
	var menu entity.Menu
	err := db.First(&menu, id).Error
	return &menu, err
}

func (r *permissionRepository) CheckMenuInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	err := db.Model(&entity.RoleMenu{}).Where("menu_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) CheckMenuHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	err := db.Model(&entity.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// Resource implementations
func (r *permissionRepository) GetResourceList(db *gorm.DB, keyword string) ([]entity.Resource, error) {
	var list []entity.Resource
	db = db.Model(&entity.Resource{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := db.Find(&list).Error
	return list, err
}

func (r *permissionRepository) SaveOrUpdateResource(db *gorm.DB, id, parentId int, name, url, method string) error {
	if id == 0 {
		return db.Create(&entity.Resource{
			ParentId: parentId,
			Name:     name,
			Url:      url,
			Method:   method,
		}).Error
	}
	return db.Model(&entity.Resource{Model: entity.Model{ID: id}}).Updates(entity.Resource{
		ParentId: parentId,
		Name:     name,
		Url:      url,
		Method:   method,
	}).Error
}

func (r *permissionRepository) DeleteResource(db *gorm.DB, id int) error {
	return db.Delete(&entity.Resource{}, id).Error
}

func (r *permissionRepository) UpdateResourceAnonymous(db *gorm.DB, id int, anonymous bool) error {
	return db.Model(&entity.Resource{Model: entity.Model{ID: id}}).Update("anonymous", anonymous).Error
}
