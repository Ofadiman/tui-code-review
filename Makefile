dev:
	go run .

# Regenerate generated golang code after adding, removing or modifying graphql queries in genqlient.graphql file.
generate:
	go run github.com/Khan/genqlient
