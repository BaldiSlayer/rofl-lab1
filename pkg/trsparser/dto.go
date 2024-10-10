package trsparser

import (
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/models"
	"github.com/BaldiSlayer/rofl-lab1/internal/parser/trsinterprets"
	rulesparser "github.com/BaldiSlayer/rofl-lab1/internal/parser/trsrules"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../docs/trs-parser-api.yaml
