package g

// Redis Key

const (
	EMAIL_CODE   = "email_code:"   // 邮箱验证码
	ONLINE_USER  = "online_user:"  // 在线用户
	OFFLINE_USER = "offline_user:" // 强制下线用户
	VISITOR_AREA = "visitor_area"  // 地域统计
	VIEW_COUNT   = "view_count"    // 访问数量

	KEY_UNIQUE_VISITOR_SET = "unique_visitor" // 唯一用户记录 set

	ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set
	ARTICLE_LIKE_COUNT    = "article_like_count" // 文章点赞数
	ARTICLE_VIEW_COUNT    = "article_view_count" // 文章查看数

	COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set
	COMMENT_LIKE_COUNT    = "comment_like_count" // 评论点赞数

	PAGE   = "page"   // 页面封面
	CONFIG = "config" // 博客配置
)

// Gin Context Key | Session Key

const (
	CTX_DB        = "_db_field"
	CTX_RDB       = "_rdb_field"
	CTX_USER_AUTH = "_user_auth_field"
	CTX_IS_SUPER  = "_is_super_field"
)

// Config Key

const (
	CONFIG_ARTICLE_COVER     = "article_cover"
	CONFIG_IS_COMMENT_REVIEW = "is_comment_review"
	CONFIG_ABOUT             = "about"
)

// Resource 标识
const (
	RESOURCE_BACKEND_LOGIN = "/api/admin/login"
	METHOD_BACKEND_LOGIN   = "POST"
)
