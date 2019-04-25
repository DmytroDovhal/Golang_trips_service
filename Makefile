validate:
	swagger validate ./swagger/swagger.yaml


serve:
	swagger serve -F=swagger ./swagger/swagger.yaml


all:
	swagger validate ./swagger/swagger.yaml
	swagger serve -F=swagger ./swagger/swagger.yaml

tests:
	go test -covermode=count ./...

mockgen: ## Run mockgen cli fro generate mocks
	mockgen \
		-destination=database/mock.go \
		-package database \
		team-project/database TicketRepository, UserCRUD, Model