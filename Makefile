.PHONY: mock
mock:
	@mockgen -source=./repository/article_author.go 	-destination=./repository/mock/article_author_mock.go 		--package=articleRepoMock
	@mockgen -source=./repository/article_reader.go 	-destination=./repository/mock/article_reader_mock.go 		--package=articleRepoMock
	@mockgen -source=./service/like.go 					-destination=./service/mock/like_mock.go 					--package=svcMocks
	@mockgen -source=./service/article.go 				-destination=./service/mock/article_mock.go 				--package=svcMocks
	@mockgen -source=./repository/cache/article_topn.go -destination=./repository/cache/mock/article_topn_mock.go 	--package=cacheMocks
	@go mod tidy