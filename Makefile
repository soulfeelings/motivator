.PHONY: all install dev build test clean \
	backend-install backend-dev backend-build backend-test \
	web-install web-dev web-build web-test \
	admin-install admin-dev admin-build admin-test \
	mobile-install mobile-dev mobile-build mobile-test

# ─── All ──────────────────────────────────────────────
all: install build

install: backend-install web-install admin-install mobile-install

dev:
	@echo "Usage: make backend-dev | web-dev | admin-dev | mobile-dev"

build: backend-build web-build admin-build

test: backend-test web-test admin-test mobile-test

clean: backend-clean web-clean admin-clean mobile-clean

# ─── Backend ──────────────────────────────────────────
backend-install:
	cd backend && $(MAKE) install

backend-dev:
	cd backend && $(MAKE) dev

backend-build:
	cd backend && $(MAKE) build

backend-test:
	cd backend && $(MAKE) test

backend-clean:
	cd backend && $(MAKE) clean

# ─── Web ──────────────────────────────────────────────
web-install:
	cd web && $(MAKE) install

web-dev:
	cd web && $(MAKE) dev

web-build:
	cd web && $(MAKE) build

web-test:
	cd web && $(MAKE) test

web-clean:
	cd web && $(MAKE) clean

# ─── Admin ────────────────────────────────────────────
admin-install:
	cd admin && $(MAKE) install

admin-dev:
	cd admin && $(MAKE) dev

admin-build:
	cd admin && $(MAKE) build

admin-test:
	cd admin && $(MAKE) test

admin-clean:
	cd admin && $(MAKE) clean

# ─── Mobile ───────────────────────────────────────────
mobile-install:
	cd mobile && $(MAKE) install

mobile-dev:
	cd mobile && $(MAKE) dev

mobile-build:
	cd mobile && $(MAKE) build

mobile-test:
	cd mobile && $(MAKE) test

mobile-clean:
	cd mobile && $(MAKE) clean
