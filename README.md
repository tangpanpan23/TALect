# TALect - 好未来AI教育生态

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**连接AI与优质教育内容的桥梁** - 基于Model Context Protocol（MCP）标准构建的智能教育内容服务平台。

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

### 🎯 应用场景

TALect MCP服务针对教育行业的四大核心场景提供解决方案：

| 应用场景 | 目标用户 | 核心价值 | 文档链接 |
|---------|---------|---------|---------|
| **AI助教/智能备课** | 教师、教研员 | 解放教师生产力，提升教学质量 | [详细文档](./future-mcp-server/docs/use-cases/ai-assistant.md) |
| **个性化学习与智能辅导** | 学生、家长 | 规模化因材施教 | [详细文档](./future-mcp-server/docs/use-cases/personalized-learning.md) |
| **智能作业与学情分析** | 教师、学校管理者 | 教学反馈闭环 | [详细文档](./future-mcp-server/docs/use-cases/intelligent-assessment.md) |
| **教育智能体与行政自动化** | 教务管理人员 | 机构运营效率 | [详细文档](./future-mcp-server/docs/use-cases/educational-agent.md) |

**[查看完整应用场景分析](./future-mcp-server/docs/use-cases/README.md)**

### 📚 资源服务 (Resources)

- **课程体系**: 学而思培优完整课程框架 - 教学规划和进度控制
- **知识图谱**: 学科知识关联网络 - 个性化学习路径推荐
- **教学模板**: 标准化教学流程 - 教学质量保障

## 🏗️ 项目架构

```
TALect/
└── future-mcp-server/          # TALink MCP Server (未链MCP服务器)
    ├── cmd/server/             # 主服务器入口
    ├── config/                 # 配置文件
    ├── internal/               # 内部包
    │   ├── auth/              # 认证授权
    │   ├── cache/             # 缓存管理
    │   ├── database/          # 数据库操作
    │   ├── handler/           # HTTP处理器
    │   ├── service/           # 业务逻辑层
    │   └── types/             # 类型定义
    ├── pkg/                   # 公共包
    │   ├── logger/            # 日志包
    │   └── mcp/              # MCP协议实现
    └── docs/                  # 项目文档
```

## 🚀 快速开始

### 环境要求

- Go 1.22+
- PostgreSQL 12+
- Redis 6.0+

### 启动服务

```bash
# 进入项目目录
cd future-mcp-server

# 安装依赖
go mod download

# 配置环境
cp etc/talink.yaml etc/talink-local.yaml
# 根据需要修改配置文件中的数据库和Redis连接信息

# 运行服务器
go run cmd/server/main.go -f etc/talink.yaml
```

### 验证服务

```bash
curl http://localhost:8080/health
```

### MVP演示

启动服务器后，运行演示脚本验证所有核心功能：

```bash
# 进入项目目录
cd future-mcp-server

# 运行MVP功能演示
./demo.sh
```

演示脚本将自动测试：
- ✅ MCP协议初始化
- ✅ 工具列表获取
- ✅ 教学材料搜索
- ✅ 素材详情获取
- ✅ 教案自动生成

## 🏗️ 技术架构

项目采用 **go-zero 微服务框架**，提供企业级的服务治理能力和高性能架构。

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
- **文档**: [详细文档](future-mcp-server/docs/)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

---

⭐ 如果这个项目对你有帮助，请给我们一个星标！

*本文档是TALect项目的核心介绍。如需详细的技术文档和API说明，请查看 [future-mcp-server/docs/](future-mcp-server/docs/) 目录。*
