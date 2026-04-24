## 问题
- 在后台发布的新文章前端没有显示，查看数据库发现对应的article更新了，并且响应给前端的数据里没有新的文章。
- 我判断问题出在后端获取数据库数据时或者返回响应时。

请按以下步骤：

- 阅读gin-blog-server/internal/handle/handle_front.go中的GetArticleList()和gin-blog-server/internal/model/article.go中的GetBlogArticleList()
- 看问题出现在哪里并进行修改
- 写一个测试或者直接用现有的测试文件来测试
- 确保测试通过
- 顺便检查是否有类似的隐患