# TALink MCP Server (未链MCP服务器)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**未链MCP服务器** - 连接AI与优质教育内容的桥梁，基于Model Context Protocol（MCP）标准的教育素材服务器，为AI大语言模型提供安全、高效的教育内容访问能力。

## 项目概述

**TALink MCP Server** (未链MCP服务器) 是一个基于MCP协议的教育内容服务平台，致力于：

- 🎯 **标准化接口**：提供统一的MCP协议接口，简化AI应用集成
- 🔒 **安全可控**：多层权限控制和内容保护机制
- 🚀 **高性能**：支持高并发访问和智能缓存
- 📚 **丰富内容**：整合好未来优质教育素材资源
- 🤖 **AI赋能**：为教育AI应用提供结构化内容支持

## 核心特性

### 🛠️ 工具集合 (Tools)

| 类别 | 工具 | 功能描述 |
|-----|------|---------|
| 检索类 | `search_teaching_materials` | 关键词搜索教学素材 |
| | `search_by_grade_subject` | 按年级学科筛选 |
| | `get_recommended_materials` | 个性化推荐 |
| 内容类 | `get_material_detail` | 获取素材详情 |
| | `get_related_materials` | 关联素材推荐 |
| | `extract_key_points` | 知识点提取 |
| 生成类 | `generate_lesson_plan` | 教案生成 |
| | `generate_exercises` | 练习题生成 |
| | `generate_teaching_script` | 教学脚本生成 |
| 分析类 | `analyze_material_difficulty` | 难度分析 |
| | `check_curriculum_alignment` | 课标对齐检查 |

### 📚 资源集合 (Resources)

- **课程大纲**: `curriculum://grade-{grade}/subject-{subject}`
- **知识图谱**: `knowledge-graph://subject-{subject}/level-{level}`
- **教学模板**: `template://type-{type}/model-{model}`

## 项目概述

**TALink MCP Server** (未链MCP服务器) 是一个基于MCP协议的教育内容服务平台，致力于：

- 🎯 **标准化接口**：提供统一的MCP协议接口，简化AI应用集成
- 🔒 **安全可控**：多层权限控制和内容保护机制
- 🚀 **高性能**：支持高并发访问和智能缓存
- 📚 **丰富内容**：整合好未来优质教育素材资源
- 🤖 **AI赋能**：为教育AI应用提供结构化内容支持

## 核心特性

### 🛠️ 工具集合 (Tools)

| 类别 | 工具 | 功能描述 |
|-----|------|---------|
| 检索类 | `search_teaching_materials` | 关键词搜索教学素材 |
| | `search_by_grade_subject` | 按年级学科筛选 |
| | `get_recommended_materials` | 个性化推荐 |
| 内容类 | `get_material_detail` | 获取素材详情 |
| | `get_related_materials` | 关联素材推荐 |
| | `extract_key_points` | 知识点提取 |
| 生成类 | `generate_lesson_plan` | 教案生成 |
| | `generate_exercises` | 练习题生成 |
| | `generate_teaching_script` | 教学脚本生成 |
| 分析类 | `analyze_material_difficulty` | 难度分析 |
| | `check_curriculum_alignment` | 课标对齐检查 |

### 📚 资源集合 (Resources)

- **课程大纲**: `curriculum://grade-{grade}/subject-{subject}`
- **知识图谱**: `knowledge-graph://subject-{subject}/level-{level}`
- **教学模板**: `template://type-{type}/model-{model}`

## 快速开始

### 环境要求

- Go 1.22+
- PostgreSQL 12+
- Redis 6.0+

### 安装依赖

```bash
go mod download
```

### 配置环境

1. 复制配置文件模板：
```bash
cp config/config.example.yaml config/config.yaml
```

2. 编辑配置文件，设置数据库和Redis连接信息

### 运行服务器

```bash
# 开发模式
go run cmd/server/main.go

# 或使用Makefile
make dev
```

### 二进制部署

```bash
# 构建生产版本
make build

# 运行服务
./build/future-mcp-server

# 或使用systemd服务管理
sudo make service-install  # 安装服务
sudo make service-start    # 启动服务
sudo make service-status   # 查看状态
```

### 自动化部署

```bash
# 生产环境完整部署
sudo ./deploy/deploy.sh prod

# 开发环境部署
sudo ./deploy/deploy.sh dev

# 回滚到上一版本
sudo ./deploy/deploy.sh rollback
```

## 项目结构

```
future-mcp-server/
├── cmd/                    # 应用入口
│   └── server/            # 主服务器
├── config/                # 配置文件
├── docs/                  # 项目文档
│   ├── api/              # API文档
│   ├── architecture/     # 架构设计
│   └── requirements/     # 需求文档
├── internal/              # 内部包（不对外暴露）
│   ├── auth/             # 认证授权
│   ├── cache/            # 缓存管理
│   ├── database/         # 数据库操作
│   ├── handler/          # HTTP处理器
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   ├── tools/            # MCP工具实现
│   └── types/            # 类型定义
├── pkg/                  # 公共包（可对外暴露）
│   ├── logger/           # 日志包
│   ├── mcp/             # MCP协议实现
│   ├── utils/            # 工具函数
│   └── validator/        # 验证器
├── scripts/              # 构建和部署脚本
├── third/                # 第三方服务集成
│   ├── pinecone/        # 向量搜索
│   └── storage/          # 对象存储
├── deploy/               # 部署配置
│   ├── docker/          # Docker配置
│   └── k8s/             # Kubernetes配置
├── docker-compose.yml    # Docker Compose配置
├── Makefile             # 构建脚本
├── go.mod               # Go模块文件
└── README.md            # 项目说明
```

## 配置说明

### 主要配置项

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug/release

database:
  host: "localhost"
  port: 5432
  user: "future_mcp"
  password: "password"
  dbname: "future_mcp"
  sslmode: "disable"

redis:
  host: "localhost:6379"
  password: ""
  db: 0

auth:
  jwt_secret: "your-jwt-secret"
  jwt_expire: 86400

vector_search:
  provider: "pinecone"  # pinecone/weaviate
  api_key: "your-api-key"
  environment: "us-east-1"
  index_name: "future-materials"

storage:
  provider: "local"  # local/s3/cos
  bucket: "future-materials"
  region: "us-east-1"
```

## API使用示例

### MCP协议调用

```json
// 初始化连接
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "Claude Desktop",
      "version": "0.1.0"
    }
  }
}

// 调用搜索工具
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "search_teaching_materials",
    "arguments": {
      "query": "一元二次方程",
      "grade": "grade_2",
      "subject": "math",
      "limit": 10
    }
  }
}
```

### REST API调用

```bash
# 搜索素材
curl -X POST "http://localhost:8080/api/v1/materials/search" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "数学",
    "filters": {
      "grade": ["grade_1", "grade_2"],
      "subject": "math"
    },
    "pagination": {
      "page": 1,
      "size": 20
    }
  }'

# 获取素材详情
curl -X GET "http://localhost:8080/api/v1/materials/123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 开发指南

### 代码规范

项目遵循以下编码规范：

- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 使用 `go vet` 进行静态分析
- 单元测试覆盖率 > 80%

### 添加新工具

1. 在 `internal/tools/` 目录下实现工具逻辑
2. 在 `internal/service/mcp/` 中注册工具
3. 更新相关文档和测试

### 数据库迁移

```bash
# 创建新迁移
make migrate-new name=add_user_table

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down
```

## 部署说明

### 开发环境

```bash
make dev
```

### 测试环境

```bash
make test
```

### 生产环境

```bash
make build
make deploy
```

## 监控和日志

### 指标监控

- 请求响应时间
- 错误率统计
- 资源使用情况
- 工具调用统计

### 日志级别

- DEBUG: 详细调试信息
- INFO: 一般信息
- WARN: 警告信息
- ERROR: 错误信息
- FATAL: 致命错误

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系我们

- 项目维护者: Future Education Team
- 邮箱: dev@future-mcp.com
- 文档: [https://docs.future-mcp.com](https://docs.future-mcp.com)

## 致谢

感谢所有为本项目做出贡献的开发者和教育工作者！

## 实施路线图

### 第一阶段：MVP（1-2个月）
**目标**: 基础检索功能上线

- ✅ MCP服务器框架搭建
- ✅ 3个核心工具实现：
  - `search_teaching_materials` (关键词搜索)
  - `get_material_detail` (详情获取)
  - `generate_lesson_plan` (教案生成)
- ✅ 基础权限验证
- ✅ 对接好未来素材库测试环境

### 第二阶段：功能扩展（2-3个月）
**目标**: 完善工具生态

- 🔄 增加5-8个专业教育工具
- 🔄 实现向量检索能力
- 🔄 添加资源集合支持
- 🔄 完善监控和日志系统
- 🔄 性能优化和缓存策略

### 第三阶段：生态建设（3-4个月）
**目标**: 开发者生态与商业化

- 📋 SDK开发（Python/JavaScript）
- 📋 开发者门户和文档
- 📋 使用量分析和计费系统
- 📋 高级功能：个性化推荐、A/B测试

## 安全与权限设计

### 多层安全防护

1. **认证层**: API密钥 + JWT令牌
2. **授权层**: RBAC（角色权限控制）
3. **访问层**: 素材使用配额限制
4. **审计层**: 完整操作日志追踪
5. **内容层**: 水印+DRM保护

### 权限级别

| 角色 | 权限说明 | 配额限制 |
|-----|---------|---------|
| **游客** | 基础素材搜索 | 100次/天 |
| **开发者** | 全部工具调用 | 1000次/天 |
| **合作伙伴** | 批量素材访问 | 自定义配额 |
| **内部团队** | 高级分析工具 | 无限制 |

## 商业模式

### 收费策略

1. **免费层**: 基础检索，限制调用次数
2. **开发者计划**: 按调用量计费
3. **企业方案**: 定制化+技术支持
4. **内容授权**: 素材使用授权费

### 市场定位

**目标用户**: 教育科技公司、AI开发者、在线教育平台

**竞争优势**: 好未来独家高质量内容 + 标准化接口

**合作伙伴**: 与主流AI平台集成（OpenAI、Claude、文心一言等）

## 风险与应对

### 技术风险
- **风险**: MCP协议变更
- **应对**: 抽象协议层，保持向后兼容

### 内容风险
- **风险**: 版权泄露或滥用
- **应对**: 数字水印+访问追踪+法律条款

### 业务风险
- **风险**: 市场需求不足
- **应对**: 先内部试用，逐步开放，收集反馈

## 成功指标

### 技术指标
- **可用性**: 99.9% SLA
- **响应时间**: P95 < 200ms
- **并发支持**: 1000+ QPS

### 业务指标
- **开发者数量**: 首年目标100+
- **API调用量**: 月均100万+
- **素材使用率**: 热门素材覆盖率80%+
- **合作伙伴**: 与3+主流AI平台集成

## 团队与资源需求

### 核心团队需求
- **后端开发**: 2-3人 (Go/TypeScript)
- **AI算法工程师**: 1-2人 (NLP/推荐算法)
- **产品经理**: 1人 (教育+AI背景)
- **测试/运维**: 1-2人

### 资源需求
- **基础设施**: 云服务器、CDN、数据库
- **开发工具**: Git、CI/CD、监控系统
- **内容准备**: 素材数字化、标注、向量化

## 下一步建议

1. **启动技术验证**: 先用小规模素材验证技术可行性
2. **内部试用**: 先让好未来内部产品团队试用
3. **寻找早期合作伙伴**: 与1-2个教育AI初创公司合作试点
4. **参加AI开发者大会**: 展示MCP服务器的能力

## 联系我们

- **项目维护者**: Future Education Team
- **邮箱**: tangpan1@tal.com
- **文档**: [项目文档](docs/)

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！
