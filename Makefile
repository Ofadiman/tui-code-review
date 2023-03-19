dev:
	@go run .

# Regenerate generated golang code after adding, removing or modifying graphql queries in genqlient.graphql file.
gql_generate:
	@go run github.com/Khan/genqlient

settings_setup:
	@cp ./ignored/settings.json ~/.tui-code-review.json
	
settings_clean:
	@rm ~/.tui-code-review.json
	