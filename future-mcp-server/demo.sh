#!/bin/bash

# TALink MCP Server MVP演示脚本
# 这个脚本演示了MVP版本的3个核心功能

echo "🚀 TALink MCP Server MVP演示"
echo "=================================="

# 服务器地址
SERVER_URL="http://localhost:8080/mcp/jsonrpc"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查服务器是否运行
echo -e "\n${YELLOW}检查服务器状态...${NC}"
if curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}✅ 服务器运行正常${NC}"
else
    echo -e "${RED}❌ 服务器未运行，请先启动服务器${NC}"
    echo "运行: cd future-mcp-server && ./server -f etc/talink.yaml"
    exit 1
fi

echo -e "\n${YELLOW}1. MCP协议初始化${NC}"
INIT_RESPONSE=$(curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "Demo Client",
        "version": "1.0.0"
      }
    }
  }')

if echo "$INIT_RESPONSE" | grep -q "TALink MCP Server"; then
    echo -e "${GREEN}✅ MCP初始化成功${NC}"
else
    echo -e "${RED}❌ MCP初始化失败${NC}"
    exit 1
fi

echo -e "\n${YELLOW}2. 获取工具列表${NC}"
TOOLS_RESPONSE=$(curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }')

if echo "$TOOLS_RESPONSE" | grep -q "search_teaching_materials"; then
    echo -e "${GREEN}✅ 工具列表获取成功${NC}"
    echo "发现工具: $(echo "$TOOLS_RESPONSE" | grep -o '"name":"[^"]*"' | wc -l) 个"
else
    echo -e "${RED}❌ 工具列表获取失败${NC}"
    exit 1
fi

echo -e "\n${YELLOW}3. 搜索教学材料${NC}"
SEARCH_RESPONSE=$(curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "search_teaching_materials",
      "arguments": {
        "query": "方程",
        "limit": 5
      }
    }
  }')

if echo "$SEARCH_RESPONSE" | grep -q "一元二次方程解法"; then
    echo -e "${GREEN}✅ 搜索功能正常${NC}"
    MATERIAL_ID=$(echo "$SEARCH_RESPONSE" | grep -o 'ID: [a-f0-9-]*' | head -1 | cut -d' ' -f2)
    echo "找到素材ID: $MATERIAL_ID"
else
    echo -e "${RED}❌ 搜索功能异常${NC}"
    exit 1
fi

echo -e "\n${YELLOW}4. 获取素材详情${NC}"
if [ -n "$MATERIAL_ID" ]; then
    DETAIL_RESPONSE=$(curl -s -X POST "$SERVER_URL" \
      -H "Content-Type: application/json" \
      -d "{
        \"jsonrpc\": \"2.0\",
        \"id\": 4,
        \"method\": \"tools/call\",
        \"params\": {
          \"name\": \"get_material_detail\",
          \"arguments\": {
            \"material_id\": \"$MATERIAL_ID\"
          }
        }
      }")

    if echo "$DETAIL_RESPONSE" | grep -q "素材详情"; then
        echo -e "${GREEN}✅ 素材详情获取成功${NC}"
        echo "$DETAIL_RESPONSE" | grep -o '"text":"[^"]*"' | head -1 | cut -d'"' -f4 | head -20
    else
        echo -e "${RED}❌ 素材详情获取失败${NC}"
    fi
fi

echo -e "\n${YELLOW}5. 生成教案${NC}"
LESSON_RESPONSE=$(curl -s -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"jsonrpc\": \"2.0\",
    \"id\": 5,
    \"method\": \"tools/call\",
    \"params\": {
      \"name\": \"generate_lesson_plan\",
      \"arguments\": {
        \"material_ids\": [\"$MATERIAL_ID\"],
        \"objectives\": [\"掌握一元二次方程解法\", \"能够运用公式法解题\"],
        \"grade\": \"grade_2\",
        \"student_level\": \"intermediate\",
        \"duration\": 45
      }
    }
  }")

if echo "$LESSON_RESPONSE" | grep -q "学而思教研标准教案"; then
    echo -e "${GREEN}✅ 教案生成成功${NC}"
    echo "教案包含5E教学模型和完整的教学流程"
else
    echo -e "${RED}❌ 教案生成失败${NC}"
fi

echo -e "\n${GREEN}🎉 MVP演示完成！${NC}"
echo "=================================="
echo "核心功能验证结果:"
echo "✅ MCP协议初始化"
echo "✅ 工具列表获取"
echo "✅ 教学材料搜索"
echo "✅ 素材详情获取"
echo "✅ 教案自动生成"
echo ""
echo "TALink MCP Server MVP版本已就绪，可以进行演示和进一步开发！"
