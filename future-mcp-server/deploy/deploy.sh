#!/bin/bash

# TALink MCP Server 部署脚本
# 用法: ./deploy.sh [dev|prod|rollback]

set -e

# 配置变量
APP_NAME="future-mcp-server"
APP_USER="talink"
INSTALL_DIR="/opt/talink-mcp-server"
BACKUP_DIR="/opt/talink-mcp-server/backups"
CONFIG_DIR="${INSTALL_DIR}/config"
LOGS_DIR="${INSTALL_DIR}/logs"
STORAGE_DIR="${INSTALL_DIR}/storage"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查root权限
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "此脚本需要root权限运行"
        exit 1
    fi
}

# 创建用户
create_user() {
    if ! id "$APP_USER" &>/dev/null; then
        log_info "创建用户 $APP_USER"
        useradd -r -s /bin/false "$APP_USER"
    else
        log_info "用户 $APP_USER 已存在"
    fi
}

# 创建目录
create_directories() {
    log_info "创建应用目录"
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$LOGS_DIR"
    mkdir -p "$STORAGE_DIR"
    mkdir -p "$BACKUP_DIR"

    chown -R "$APP_USER:$APP_USER" "$INSTALL_DIR"
    chmod 755 "$INSTALL_DIR"
    chmod 755 "$CONFIG_DIR"
    chmod 755 "$LOGS_DIR"
    chmod 755 "$STORAGE_DIR"
}

# 备份当前版本
backup_current() {
    if [[ -f "$INSTALL_DIR/$APP_NAME" ]]; then
        local timestamp=$(date +%Y%m%d_%H%M%S)
        local backup_file="$BACKUP_DIR/${APP_NAME}_${timestamp}.bak"

        log_info "备份当前版本到 $backup_file"
        cp "$INSTALL_DIR/$APP_NAME" "$backup_file"
        cp "$CONFIG_DIR/config.yaml" "$BACKUP_DIR/config_${timestamp}.yaml" 2>/dev/null || true
    fi
}

# 复制文件
copy_files() {
    log_info "复制应用文件"

    # 复制二进制文件
    if [[ ! -f "build/$APP_NAME" ]]; then
        log_error "找不到构建文件 build/$APP_NAME，请先运行 make build"
        exit 1
    fi

    cp "build/$APP_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$APP_NAME"

    # 复制配置文件（如果不存在）
    if [[ ! -f "$CONFIG_DIR/config.yaml" ]]; then
        if [[ -f "config/config.yaml" ]]; then
            cp "config/config.yaml" "$CONFIG_DIR/"
        else
            log_warn "未找到配置文件，请手动配置 $CONFIG_DIR/config.yaml"
        fi
    fi

    chown -R "$APP_USER:$APP_USER" "$INSTALL_DIR"
}

# 安装systemd服务
install_service() {
    log_info "安装systemd服务"

    # 复制服务文件
    cp "deploy/$APP_NAME.service" "/etc/systemd/system/"

    # 重新加载systemd
    systemctl daemon-reload

    log_info "服务已安装，可使用以下命令管理："
    echo "  启动: systemctl start $APP_NAME"
    echo "  停止: systemctl stop $APP_NAME"
    echo "  重启: systemctl restart $APP_NAME"
    echo "  状态: systemctl status $APP_NAME"
    echo "  开机自启: systemctl enable $APP_NAME"
}

# 启动服务
start_service() {
    log_info "启动服务"
    systemctl start "$APP_NAME"

    # 等待服务启动
    sleep 3

    # 检查服务状态
    if systemctl is-active --quiet "$APP_NAME"; then
        log_info "服务启动成功"
    else
        log_error "服务启动失败，请检查日志"
        journalctl -u "$APP_NAME" -n 20 --no-pager
        exit 1
    fi
}

# 停止服务
stop_service() {
    log_info "停止服务"
    systemctl stop "$APP_NAME" || true
}

# 回滚到上一版本
rollback() {
    log_info "执行回滚"

    # 停止服务
    stop_service

    # 查找最新的备份文件
    local latest_backup=$(ls -t "$BACKUP_DIR/${APP_NAME}"_*.bak 2>/dev/null | head -1)

    if [[ -z "$latest_backup" ]]; then
        log_error "未找到备份文件"
        exit 1
    fi

    log_info "回滚到 $latest_backup"
    cp "$latest_backup" "$INSTALL_DIR/$APP_NAME"
    chmod +x "$INSTALL_DIR/$APP_NAME"
    chown "$APP_USER:$APP_USER" "$INSTALL_DIR/$APP_NAME"

    # 恢复配置文件
    local config_backup=$(ls -t "$BACKUP_DIR/config_"*.yaml 2>/dev/null | head -1)
    if [[ -n "$config_backup" ]]; then
        cp "$config_backup" "$CONFIG_DIR/config.yaml"
        chown "$APP_USER:$APP_USER" "$CONFIG_DIR/config.yaml"
    fi

    # 启动服务
    start_service

    log_info "回滚完成"
}

# 检查系统要求
check_system() {
    log_info "检查系统要求"

    # 检查操作系统
    if [[ ! -f /etc/os-release ]]; then
        log_error "不支持的操作系统"
        exit 1
    fi

    # 检查必要的命令
    local commands=("systemctl" "useradd" "chown" "chmod")
    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            log_error "缺少必要的命令: $cmd"
            exit 1
        fi
    done

    log_info "系统检查通过"
}

# 显示使用帮助
show_help() {
    echo "TALink MCP Server 部署脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  dev      开发环境部署（不安装systemd服务）"
    echo "  prod     生产环境部署（完整部署）"
    echo "  rollback 回滚到上一版本"
    echo "  help     显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 prod    # 生产环境部署"
    echo "  $0 dev     # 开发环境部署"
}

# 开发环境部署
deploy_dev() {
    log_info "开始开发环境部署"

    create_user
    create_directories
    backup_current
    copy_files

    log_info "开发环境部署完成"
    log_info "手动启动服务: cd $INSTALL_DIR && ./$APP_NAME"
}

# 生产环境部署
deploy_prod() {
    log_info "开始生产环境部署"

    check_system
    create_user
    create_directories
    backup_current
    copy_files
    install_service
    start_service

    log_info "生产环境部署完成"
    log_info "服务状态: systemctl status $APP_NAME"
}

# 主函数
main() {
    local command=${1:-help}

    case $command in
        dev)
            check_root
            deploy_dev
            ;;
        prod)
            check_root
            deploy_prod
            ;;
        rollback)
            check_root
            rollback
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
