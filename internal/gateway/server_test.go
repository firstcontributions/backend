package gateway

import (
	"io/ioutil"
	"testing"

	graphqlschema "github.com/firstcontributions/backend/internal/graphql/schema"
	"github.com/graph-gophers/graphql-go"
	otelgraphql "github.com/graph-gophers/graphql-go/trace/otel"
)

func TestServer_GetGraphqlSchema(t *testing.T) {
	schema, err := ioutil.ReadFile("../../assets/schema.graphql")
	if err != nil {
		t.Errorf("error on reading schema %v", err)
		return
	}
	resolver := &graphqlschema.Resolver{}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.Tracer(otelgraphql.DefaultTracer())}
	_, err = graphql.ParseSchema(string(schema), resolver, opts...)
	if err != nil {
		t.Errorf("error on parsing schema %v", err)
		return
	}
}
