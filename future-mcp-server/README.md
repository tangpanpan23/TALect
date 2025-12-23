# TALink MCP Server (未链MCP服务器)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**连接AI与优质教育内容的桥梁** - 基于Model Context Protocol（MCP）标准的教育素材服务器，为AI大语言模型提供安全、高效的教育内容访问能力。

## 🌟 项目特色

- 🎯 **标准化接口**: 基于MCP协议，简化AI应用集成
- 🔒 **企业级安全**: 多层权限控制和内容保护机制
- 🚀 **高性能架构**: Go语言驱动，支持高并发访问
- 📚 **独家内容**: 整合好未来20年精品教育内容资源
- 🤖 **AI教育技术**: 结合个性化学习算法和大数据能力
- 🌍 **生态开放**: 打造教育AI内容服务平台

## 📋 核心功能

### 🛠️ AI工具生态 (Tools)

| 类别 | 核心工具 | 技术亮点 |
|-----|---------|---------|
| **智能检索** | 语义搜索、个性化推荐 | 向量检索 + AI算法 |
| **内容处理** | 知识点提取、智能摘要 | NLP技术 + 教育专业性 |
| **教学生成** | 教案生成、练习题生成 | AI创作 + 教学标准化 |
| **学习分析** | 难度评估、学习路径 | 大数据分析 + 自适应学习 |

### 📚 资源服务 (Resources)

- **课程体系**: 学而思培优完整课程框架 - 教学规划和进度控制
- **知识图谱**: 学科知识关联网络 - 个性化学习路径推荐
- **教学模板**: 标准化教学流程 - 教学质量保障

## 🚀 快速开始

### 环境要求

- Go 1.22+
- PostgreSQL 12+
- Redis 6.0+

### 安装依赖

```bash
go mod download
```

### 配置环境

```bash
# 复制配置文件模板
cp config/config.example.yaml config/config.yaml

# 编辑配置文件，设置数据库和Redis连接信息
```

### 运行服务器

```bash
# 开发模式
go run cmd/server/main.go

# 或使用Makefile
make dev
```

### 验证安装

```bash
curl http://localhost:8080/health
```

## 🏗️ 项目架构

```
future-mcp-server/
├── cmd/server/          # 主服务器入口
├── config/              # 配置文件
├── internal/            # 内部包
│   ├── auth/           # 认证授权
│   ├── cache/          # 缓存管理
│   ├── database/       # 数据库操作
│   ├── handler/        # HTTP处理器
│   ├── middleware/     # 中间件
│   ├── service/        # 业务逻辑层
│   └── types/          # 类型定义
├── pkg/                # 公共包
│   ├── logger/         # 日志包
│   └── mcp/           # MCP协议实现
├── docs/               # 项目文档
└── deploy/             # 部署配置
```

## 💡 核心价值主张

### 1. **内容独家性**
- **20年教育积累**: 整合学而思培优、考研帮等业务线的精品教育内容
- **教研专业性**: 基于好未来强大的教研团队和教学经验
- **持续更新**: 紧跟教育改革和教学发展趋势

### 2. **技术领先性**
- **AI原生设计**: 深度整合好未来的AI教育算法和大数据能力
- **性能极致**: Go语言驱动，支持大规模并发访问
- **架构现代化**: 微服务架构，支持水平扩展

### 3. **生态开放性**
- **标准协议**: 基于MCP协议，降低外部集成门槛
- **开发者友好**: 提供完整的SDK和开发工具
- **商业可持续**: 支持多种商业模式和合作方式

## 🔧 技术栈

- **后端框架**: Go 1.22+ + Gin
- **数据库**: PostgreSQL + pgvector (向量搜索)
- **缓存**: Redis
- **协议**: MCP (Model Context Protocol)
- **认证**: JWT + API Key
- **日志**: Zap (结构化日志)

## 📖 API使用示例

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
    }
  }'
```

## 🚀 部署方式

### 开发环境
```bash
make dev
```

### 生产环境
```bash
# 构建
make build

# 二进制部署
./build/future-mcp-server

# 或systemd服务管理
sudo make service-install
sudo make service-start
```

### 自动化部署
```bash
# 生产环境完整部署
sudo ./deploy/deploy.sh prod

# 开发环境部署
sudo ./deploy/deploy.sh dev
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 开发规范

- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 使用 `go vet` 进行静态分析
- 单元测试覆盖率 > 80%

## 📈 实施路线图

### 第一阶段：MVP（2-3个月）
- ✅ MCP服务器框架搭建 (Go + Gin)
- ✅ 3个核心工具实现：搜索、详情获取、教案生成
- ✅ 基础权限验证和审计体系
- ✅ 对接好未来精品课程内容库

### 第二阶段：功能完善（3-4个月）
- 🔄 增加5-8个专业教育工具
- 🔄 实现向量检索能力
- 🔄 添加资源集合支持
- 🔄 完善监控和日志系统

### 第三阶段：生态建设（4-6个月）
- 📋 SDK开发和开发者门户
- 📋 使用量分析和计费系统
- 📋 高级功能：个性化推荐、A/B测试

## 📞 联系我们

- **项目维护者**: Future Education Team
- **邮箱**: tangpan1@tal.com
- **文档**: [项目文档](docs/)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！

*本文档是TALink MCP Server项目的核心介绍文档。如需详细的项目规划和技术文档，请查看 [docs/](docs/) 目录。*