package gen

import (
	`fmt`
	`go/token`
	`go/types`
	`strings`
)

var (
	globalIdent = names(
		"AggregateFunc",
		"As",
		"Asc",
		"Client",
		"config",
		"Count",
		"Debug",
		"Desc",
		"Driver",
		"Hook",
		"Log",
		"MutateFunc",
		"Mutation",
		"Mutator",
		"Op",
		"Option",
		"OrderFunc",
		"Max",
		"Mean",
		"Min",
		"Sum",
		"Policy",
		"Query",
		"Value",
	)
)

func ValidSchemaName(name string) error {
	pkg := strings.ToLower(name)
	if token.Lookup(pkg).IsKeyword() {
		return fmt.Errorf("schema lowercase name conflicts with Go keyword %q", pkg)
	}
	if types.Universe.Lookup(pkg) != nil {
		return fmt.Errorf("schema lowercase name conflicts with Go predeclared identifier %q", pkg)
	}
	if _, ok := globalIdent[pkg]; ok {
		return fmt.Errorf("schema lowercase name conflicts ent predeclared identifier %q", pkg)
	}
	if _, ok := globalIdent[name]; ok {
		return fmt.Errorf("schema name conflicts with ent predeclared identifier %q", name)
	}
	return nil
}

func names(ids ...string) map[string]struct{} {
	m := make(map[string]struct{})
	for i := range ids {
		m[ids[i]] = struct{}{}
	}
	return m
}
