package convert

import (
	"fmt"

	"github.com/Quad4-Software/cloud/pkg/kql"
	"github.com/Quad4-Software/cloud/services/search/pkg/opensearch/internal/osu"
	"github.com/Quad4-Software/cloud/services/search/pkg/query"
)

var (
	ErrUnsupportedNodeType = fmt.Errorf("unsupported node type")
)

func KQLToOpenSearchBoolQuery(kqlQuery string) (*osu.BoolQuery, error) {
	kqlAst, err := kql.Builder{}.Build(kqlQuery)
	if err != nil {
		return nil, err
	}

	// shared lowering: field resolution, media-type expansion, value lowercasing.
	kqlAst = query.Normalize(kqlAst, query.ResolveField)

	builder, err := TranspileKQLToOpenSearch(kqlAst.Nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile query: %w", err)
	}

	if q, ok := builder.(*osu.BoolQuery); !ok {
		return osu.NewBoolQuery().Must(builder), nil
	} else {
		return q, nil
	}
}
