.PHONY: mock
mock:
	@mockgen -source=./service/article.go -destination=./service/mock/article_mock.go --package=articleServiceMock
	@mockgen -source=./repository/article_author.go -destination=./repository/mock/article_author_mock.go --package=articleRepoMock
	@mockgen -source=./repository/article_reader.go -destination=./repository/mock/article_reader_mock.go --package=articleRepoMock
	@go mod tidy