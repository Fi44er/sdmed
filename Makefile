include ./scripts/server.mk
include ./scripts/docker.mk

# Цвета для текста
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
CYAN=\033[0;36m
NC=\033[0m # No Color

# Определение переменных для параметров
migrate_flag ?= false
redis_mode ?= 0

.PHONY: help run run_with_migrate generate_dock start_postgres connect_postgres start_redis connect_redis start_all stop_all view_logs

# Ширина столбца для команд (можно настроить)
WIDTH=20

help:
	@echo "${BLUE}=== Makefile Help ===${NC}"
	@echo ""
	@echo "${YELLOW}Usage:${NC}"
	@echo "  make [command] [options]"
	@echo ""
	@echo "${YELLOW}App Commands:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "run_app" "Start the app with default settings (migrate_flag=${migrate_flag}, redis_mode=${redis_mode})"
	@printf "  ${CYAN}%-${WIDTH}s${NC}%s\n" "" "Override example: make run_app redis_mode=2 migrate_flag=true"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "gen_mock" "Generate mock for src"
	@printf "  ${CYAN}%-${WIDTH}s${NC}%s\n" "" "Override example: make gen_mock src=./cmd/handler/user.go"
	@echo ""
	@echo "${YELLOW}Swagger Commands:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "gen_dock" "Generate Swagger documentation"
	@echo ""
	@echo "${YELLOW}Docker Commands:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "start_postgres" "Start PostgreSQL container"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "connect_postgres" "Connect to PostgreSQL container via bash"
	@echo ""
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "start_redis" "Start Redis container"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "connect_redis" "Connect to Redis container via bash"
	@echo ""
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "start_all" "Start all services (Postgres, Redis, Backend)"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "stop_all" "Stop all running containers"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "docker_logs" "View docker logs for all services"
	@echo ""
	@echo "${YELLOW}Options:${NC}"
	@printf "  ${CYAN}%-${WIDTH}s${NC}%s\n" "MIGRATE_FLAG" "Enable or disable migrations (true|false, default: ${MIGRATE_FLAG})"
	@printf "  ${CYAN}%-${WIDTH}s${NC}%s\n" "REDIS_MODE" "Set Redis cache mode (0|1|2, default: ${REDIS_MODE})"
	@printf "  %-${WIDTH}s%s\n" "" " 0: No flush, 1: Selective flush, 2: Full flush"

