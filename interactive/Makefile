.PHONY: mock
mock:
	@mockgen -source=./service/like.go 					-destination=./service/mock/like_mock.go 					--package=svcMocks
	@mockgen -source=./repository/cache/like.go 		-destination=./repository/cache/mock/like_mock.go 			--package=cacheMocks
	@go mod tidy