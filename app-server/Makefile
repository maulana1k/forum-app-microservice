# -------------------------
# Makefile for DB Migrations
# -------------------------

# Go migration runner
MIGRATE_CMD = go run platform/main.go

# Migration directory
MIGRATION_DIR = platform/migrations

# Database URL (set via env or default)
DB_URL ?= postgres://user:password@localhost:5432/forum_db?sslmode=disable

# -------------------------
# Migration commands
# -------------------------

# Apply all migrations
migration.up:
	$(MIGRATE_CMD) -action=up

# Rollback last migration
migration.down:
	$(MIGRATE_CMD) -action=down

# Drop all tables and reset (use with caution!)
migration.reset:
	$(MIGRATE_CMD) -action=reset

# Show current migration version
migration.status:
	$(MIGRATE_CMD) -action=status

# Create a new migration (up + down files)
# Usage: make migration.create name=create_posts_table
migration.create:
	@if [ -z "$(name)" ]; then \
		echo "Error: please provide a migration name: make migration.create name=create_posts_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)
