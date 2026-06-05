#!/bin/sh
set -e
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

export MISE_HIDE_UPDATE_WARNING=1

# 日志输出格式
COLOR_PREFIX="\033[1;36m[Entrypoint]\033[0m"
log() {
  printf "${COLOR_PREFIX} %s\n" "$1"
}

# ============================
# 数据目录结构（青龙风格）
# ============================
# /app/data/
# ├── config/    配置文件（config.ini）
# ├── db/        数据库（xuanwu.db）
# ├── scripts/   用户脚本
# ├── bak/       备份文件
# ├── deps/      语言运行时（mise）
# └── log/       运行日志
DATA_DIR="/app/data"
export XW_CONFIG_PATH="$DATA_DIR/config/config.ini"

log "Initializing data directory: $DATA_DIR ..."

mkdir -p \
  "$DATA_DIR/config" \
  "$DATA_DIR/db" \
  "$DATA_DIR/scripts" \
  "$DATA_DIR/bak" \
  "$DATA_DIR/deps" \
  "$DATA_DIR/log"

# 同步示例配置（如果宿主机目录为空）
if [ -d "/app/config" ] && [ -z "$(ls -A "$DATA_DIR/config" 2>/dev/null)" ]; then
  cp -a /app/config/. "$DATA_DIR/config/" 2>/dev/null || true
  log "Example config synced"
fi

# 兼容旧版：如果旧的独立目录存在且不是符号链接，迁移到 data/ 下
for dir_name in configs envs; do
  if [ -d "/app/$dir_name" ] && [ ! -L "/app/$dir_name" ]; then
    log "Migrating /app/$dir_name to $DATA_DIR/ ..."
    case "$dir_name" in
      configs)
        if [ -z "$(ls -A "$DATA_DIR/config" 2>/dev/null)" ] && [ -n "$(ls -A "/app/$dir_name" 2>/dev/null)" ]; then
          cp -a "/app/$dir_name/." "$DATA_DIR/config/"
        fi
        ;;
      envs)
        if [ -z "$(ls -A "$DATA_DIR/deps" 2>/dev/null)" ] && [ -n "$(ls -A "/app/$dir_name" 2>/dev/null)" ]; then
          cp -a "/app/$dir_name/." "$DATA_DIR/deps/"
        fi
        ;;
    esac
    rm -rf "/app/$dir_name"
    log "Migrated /app/$dir_name"
  fi
done

# 兼容旧版：如果 /app/data/xuanwu.db 存在，移到 data/db/
if [ -f "$DATA_DIR/xuanwu.db" ]; then
  mv "$DATA_DIR/xuanwu.db" "$DATA_DIR/db/xuanwu.db" 2>/dev/null || true
  log "Migrated xuanwu.db to $DATA_DIR/db/"
fi

# 兼容旧版：如果 /app/data/backups/ 存在，重命名为 data/bak/
if [ -d "$DATA_DIR/backups" ] && [ ! -L "$DATA_DIR/backups" ]; then
  if [ -z "$(ls -A "$DATA_DIR/bak" 2>/dev/null)" ] && [ -n "$(ls -A "$DATA_DIR/backups" 2>/dev/null)" ]; then
    cp -a "$DATA_DIR/backups/." "$DATA_DIR/bak/"
  fi
  rm -rf "$DATA_DIR/backups"
  log "Migrated backups -> bak"
fi

# 同步示例脚本
if [ -d "/app/example" ]; then
  mkdir -p "$DATA_DIR/scripts/example"
  rsync -a --ignore-existing /app/example/ "$DATA_DIR/scripts/example/" || true
  log "Example scripts synced"
else
  log "No example directory found, skipping example sync"
fi

# ============================
# Mise 环境初始化
# ============================
MISE_DIR="$DATA_DIR/deps/mise"
mkdir -p "$MISE_DIR"
if [ -d "/opt/mise-base" ]; then
  log "Syncing mise environment from base..."
  rsync -a --ignore-existing /opt/mise-base/ "$MISE_DIR/" || true
  log "Mise environment synced"
else
  log "No base mise environment found, skipping sync"
fi

# ============================
# 环境变量注入
# ============================
export MISE_DATA_DIR="$MISE_DIR"
export MISE_CONFIG_DIR="$MISE_DIR"
export PATH="$MISE_DIR/shims:$MISE_DIR/bin:$PATH"

log "Mise PATH configured, verifying runtimes..."

# 默认启用 Python 镜像源
export PIP_INDEX_URL=${PIP_INDEX_URL:-https://pypi.org/simple}

# Node 内存限制
export NODE_OPTIONS="--max-old-space-size=256"
export PYTHONPATH=$DATA_DIR/scripts:$PYTHONPATH

# ============================
# 打印确认
# ============================
log "Checking mise..."
log "  - mise: $(mise --version 2>/dev/null | head -n 1 || echo "not found")"

log "Checking python..."
log "  - python: $(python --version 2>&1 | head -n 1 || echo "not found")"

log "Checking node..."
log "  - node: $(node --version 2>&1 | head -n 1 || echo "not found")"

log "Checking npm..."
log "  - npm: $(npm --version 2>&1 | head -n 1 || echo "not found")"

# ============================
# 启动应用
# ============================
printf "\n\033[1;32m>>> Environment setup complete. Starting Xuanwu Server...\033[0m\n\n"
printf "\033[1;36m[Data]\033[0m All data is stored in: $DATA_DIR\n"
printf "\033[1;36m[Data]\033[0m   config:   $DATA_DIR/config/\n"
printf "\033[1;36m[Data]\033[0m   database: $DATA_DIR/db/\n"
printf "\033[1;36m[Data]\033[0m   scripts:  $DATA_DIR/scripts/\n"
printf "\033[1;36m[Data]\033[0m   backup:   $DATA_DIR/bak/\n"
printf "\033[1;36m[Data]\033[0m   runtimes: $DATA_DIR/deps/\n"
printf "\033[1;36m[Data]\033[0m   logs:     $DATA_DIR/log/\n\n"

cd /app
exec xuanwu server
