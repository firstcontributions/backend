package graphql

import (
	"io/ioutil"
	"testing"

	graphqlschema "github.com/firstcontributions/backend/internal/graphql/schema"

	"github.com/graph-gophers/graphql-go"
)

func Test_SchemaWorks(t *testing.T) {
	schema, err := ioutil.ReadFile("../../assets/schema.graphql")
	if err != nil {
		t.Errorf("unexpected error on reading schema %v", err)
	}

	resolver := &graphqlschema.Resolver{}
	if _, err := graphql.ParseSchema(string(schema), resolver); err != nil {
		t.Errorf("unexpected error on parsing schema schema %v", err)
	}
}
