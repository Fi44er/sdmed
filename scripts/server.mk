run_app:
	go run ./cmd/main.go -migrate=${migrate_flag} -redis=${redis_mode}

run_tree:
	go run ./scripts/src/tree.go

gen_mock:
	@echo "Generating mock for $(src)"
	@mkdir -p $(dir $(src))/mock
	@go run github.com/golang/mock/mockgen@latest \
		-source=$(src) \
		-destination=$(dir $(src))/mock/$(notdir $(basename $(src)))_mock.go \
		-package=mock
