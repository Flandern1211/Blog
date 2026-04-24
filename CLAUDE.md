

# === 代码红线（IMPORTANT 标注） ===
-IMPORTANT: 所有数据库操作必须在repository层，Controller和Service层 禁止直接操作 DB
-IMPORTANT: 对外 API 的响应格式必须遵循 @docs/api-response-format.md
-IMPORTANT: 数据库迁移必须可逆（有 up 必须有 down），大表变更必须考虑锁
-YOU MUST: 新增 API 端点必须同时更新@docs中的文档
-YOU MUST: 涉及用户数据的操作必须有审计日志


# === 代码风格 ===
- 错误处理:所有错误都必须显式处理，不允许不处理

# === 性能约定 ===
-列表接口必须分页，默认 20 条，最大 100 条
-关联数据用 JOIN 或 eager loading，严禁循环查询（N+1）
-频繁读取的配置数据使用缓存，过期时间不超过 5 分钟

# === 上下文压缩保护 ===
When compacting, ALWAYS preserve:
-已修改文件列表和每个文件改了什么
-测试运行结果
-架构红线约束
-当前 spec 的引用路径