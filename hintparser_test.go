// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package parser_test

import (
	"github.com/arana-db/parser/model"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arana-db/parser"
	"github.com/arana-db/parser/ast"
	"github.com/arana-db/parser/mysql"
)

func TestParseHint(t *testing.T) {
	testCases := []struct {
		input  string
		mode   mysql.SQLMode
		output []*ast.TableOptimizerHint
		errs   []string
	}{
		{
			input: "",
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "MEMORY_QUOTA(8 MB) MEMORY_QUOTA(6 GB)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("MEMORY_QUOTA"),
					HintData: int64(8 * 1024 * 1024),
				},
				{
					HintName: model.NewCIStr("MEMORY_QUOTA"),
					HintData: int64(6 * 1024 * 1024 * 1024),
				},
			},
		},
		{
			input: "QB_NAME(qb1) QB_NAME(`qb2`), QB_NAME(TRUE) QB_NAME(\"ANSI quoted\") QB_NAME(_utf8), QB_NAME(0b10) QB_NAME(0x1a)",
			mode:  mysql.ModeANSIQuotes,
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("qb2"),
				},
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("TRUE"),
				},
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("ANSI quoted"),
				},
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("_utf8"),
				},
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("0b10"),
				},
				{
					HintName: model.NewCIStr("QB_NAME"),
					QBName:   model.NewCIStr("0x1a"),
				},
			},
		},
		{
			input: "QB_NAME(1)",
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "QB_NAME('string literal')",
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "QB_NAME(many identifiers)",
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "QB_NAME(@qb1)",
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "QB_NAME(b'10')",
			errs: []string{
				`Cannot use bit-value literal`,
				`Optimizer hint syntax error at line 1 `,
			},
		},
		{
			input: "QB_NAME(x'1a')",
			errs: []string{
				`Cannot use hexadecimal literal`,
				`Optimizer hint syntax error at line 1 `,
			},
		},
		// ** Table **//
		// BKA and NO_BKA's Applicable Scopes `Query block` and `Table`
		{
			input: "BKA() BKA(@qb1) BKA(@qb1 tbl1) BKA(@qb1 tbl1,tbl2) BKA(tbl1) BKA(tbl1@qb1) BKA(tbl1@qb1, tbl2@qb2) " +
				"NO_BKA() NO_BKA(@qb1) NO_BKA(@qb1 tbl1) NO_BKA(@qb1 tbl1,tbl2) NO_BKA(tbl1) NO_BKA(tbl1@qb1) NO_BKA(tbl1@qb1, tbl2@qb2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("BKA"),
				},
				{
					HintName: model.NewCIStr("BKA"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("BKA"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("BKA"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("BKA"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("BKA"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("BKA"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_BKA"),
				},
				{
					HintName: model.NewCIStr("NO_BKA"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("NO_BKA"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_BKA"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("NO_BKA"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_BKA"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_BKA"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		// BNL and NO_BNL's Applicable Scopes `Query block` and `Table`
		{
			input: "BNL() BNL(@qb1) BNL(@qb1 tbl1) BNL(@qb1 tbl1,tbl2) BNL(tbl1) BNL(tbl1@qb1) BNL(tbl1@qb1, tbl2@qb2) " +
				"NO_BNL() NO_BNL(@qb1) NO_BNL(@qb1 tbl1) NO_BNL(@qb1 tbl1,tbl2) NO_BNL(tbl1) NO_BNL(tbl1@qb1) NO_BNL(tbl1@qb1, tbl2@qb2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("BNL"),
				},
				{
					HintName: model.NewCIStr("BNL"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("BNL"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("BNL"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("BNL"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("BNL"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("BNL"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_BNL"),
				},
				{
					HintName: model.NewCIStr("NO_BNL"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("NO_BNL"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_BNL"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("NO_BNL"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_BNL"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_BNL"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		// DERIVED_CONDITION_PUSHDOWN and NO_DERIVED_CONDITION_PUSHDOWN's Applicable Scopes `Query block` and `Table`
		{
			input: "DERIVED_CONDITION_PUSHDOWN() DERIVED_CONDITION_PUSHDOWN(@qb1) DERIVED_CONDITION_PUSHDOWN(@qb1 tbl1) DERIVED_CONDITION_PUSHDOWN(@qb1 tbl1,tbl2) DERIVED_CONDITION_PUSHDOWN(tbl1) DERIVED_CONDITION_PUSHDOWN(tbl1@qb1) DERIVED_CONDITION_PUSHDOWN(tbl1@qb1, tbl2@qb2) " +
				"NO_DERIVED_CONDITION_PUSHDOWN() NO_DERIVED_CONDITION_PUSHDOWN(@qb1) NO_DERIVED_CONDITION_PUSHDOWN(@qb1 tbl1) NO_DERIVED_CONDITION_PUSHDOWN(@qb1 tbl1,tbl2) NO_DERIVED_CONDITION_PUSHDOWN(tbl1) NO_DERIVED_CONDITION_PUSHDOWN(tbl1@qb1) NO_DERIVED_CONDITION_PUSHDOWN(tbl1@qb1, tbl2@qb2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
				},
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("DERIVED_CONDITION_PUSHDOWN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
				},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_DERIVED_CONDITION_PUSHDOWN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		// HASH_JOIN and NO_HASH_JOIN's Applicable Scopes `Query block` and `Table`
		{
			input: "HASH_JOIN() HASH_JOIN(@qb1) HASH_JOIN(@qb1 tbl1) HASH_JOIN(@qb1 tbl1,tbl2) HASH_JOIN(tbl1) HASH_JOIN(tbl1@qb1) HASH_JOIN(tbl1@qb1, tbl2@qb2) " +
				"NO_HASH_JOIN() NO_HASH_JOIN(@qb1) NO_HASH_JOIN(@qb1 tbl1) NO_HASH_JOIN(@qb1 tbl1,tbl2) NO_HASH_JOIN(tbl1) NO_HASH_JOIN(tbl1@qb1) NO_HASH_JOIN(tbl1@qb1, tbl2@qb2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("HASH_JOIN"),
				},
				{
					HintName: model.NewCIStr("HASH_JOIN"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("HASH_JOIN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("HASH_JOIN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("HASH_JOIN"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("HASH_JOIN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("HASH_JOIN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
				},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_HASH_JOIN"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		// INDEX_MERGE NO_INDEX_MERGE
		{
			input: "INDEX_MERGE(t1) INDEX_MERGE(t1 idx1,idx2) INDEX_MERGE(@qb1 t1 idx1,idx2) INDEX_MERGE(t1@qb1 idx1,idx2) " +
				"NO_INDEX_MERGE(t1) NO_INDEX_MERGE(t1 idx1,idx2) NO_INDEX_MERGE(@qb1 t1 idx1,idx2) NO_INDEX_MERGE(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("NO_INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		// JOIN_FIXED_ORDER's Applicable Scopes `Query block`
		{
			input: "JOIN_FIXED_ORDER() JOIN_FIXED_ORDER(@qb1) JOIN_FIXED_ORDER(@qb1 tbl1) JOIN_FIXED_ORDER(@qb1 tbl1,tbl2) " +
				"JOIN_FIXED_ORDER(tbl1) JOIN_FIXED_ORDER(tbl1@qb1) JOIN_FIXED_ORDER(tbl1@qb1, tbl2@qb2)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
				},
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_FIXED_ORDER"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},

		// ** index level **//
		// GROUP_INDEX
		{
			input: "GROUP_INDEX(t1) GROUP_INDEX(t1 idx1,idx2) GROUP_INDEX(@qb1 t1 idx1,idx2) GROUP_INDEX(t1@qb1 idx1,idx2) " +
				"NO_GROUP_INDEX(t1) NO_GROUP_INDEX(t1 idx1,idx2) NO_GROUP_INDEX(@qb1 t1 idx1,idx2) NO_GROUP_INDEX(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("NO_GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_GROUP_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},

		// MERGE and NO_MERGE's Applicable Scopes `Table`
		{
			input: "MERGE() MERGE(@qb1) MERGE(@qb1 tbl1) MERGE(@qb1 tbl1,tbl2) MERGE(tbl1) MERGE(tbl1@qb1) MERGE(tbl1@qb1, tbl2@qb2) " +
				"NO_MERGE() NO_MERGE(@qb1) NO_MERGE(@qb1 tbl1) NO_MERGE(@qb1 tbl1,tbl2) NO_MERGE(tbl1) NO_MERGE(tbl1@qb1) NO_MERGE(tbl1@qb1, tbl2@qb2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("MERGE"),
				},
				{
					HintName: model.NewCIStr("MERGE"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("MERGE"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_MERGE"),
				},
				{
					HintName: model.NewCIStr("NO_MERGE"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("NO_MERGE"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_MERGE"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("NO_MERGE"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("NO_MERGE"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("NO_MERGE"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		{
			input: "TIDB_HJ(@qb1) INL_JOIN(x, `y y`.z) MERGE_JOIN(w@`First QB`)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("TIDB_HJ"),
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("INL_JOIN"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("x")},
						{DBName: model.NewCIStr("y y"), TableName: model.NewCIStr("z")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE_JOIN"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("w"), QBName: model.NewCIStr("First QB")},
					},
				},
			},
		},
		{
			input: "USE_INDEX_MERGE(@qb1 tbl1 x, y, z) IGNORE_INDEX(tbl2@qb2) USE_INDEX(tbl3 PRIMARY) FORCE_INDEX(tbl4@qb3 c1) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("USE_INDEX_MERGE"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("tbl1")}},
					QBName:   model.NewCIStr("qb1"),
					Indexes:  []model.CIStr{model.NewCIStr("x"), model.NewCIStr("y"), model.NewCIStr("z")},
				},
				{
					HintName: model.NewCIStr("IGNORE_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("tbl2"), QBName: model.NewCIStr("qb2")}},
				},
				{
					HintName: model.NewCIStr("USE_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("tbl3")}},
					Indexes:  []model.CIStr{model.NewCIStr("PRIMARY")},
				},
				{
					HintName: model.NewCIStr("FORCE_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("tbl4"), QBName: model.NewCIStr("qb3")}},
					Indexes:  []model.CIStr{model.NewCIStr("c1")},
				},
			},
		},
		{
			input: "USE_INDEX(@qb1 tbl1 partition(p0) x) USE_INDEX_MERGE(@qb2 tbl2@qb2 partition(p0, p1) x, y, z)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("USE_INDEX"),
					Tables: []ast.HintTable{{
						TableName:     model.NewCIStr("tbl1"),
						PartitionList: []model.CIStr{model.NewCIStr("p0")},
					}},
					QBName:  model.NewCIStr("qb1"),
					Indexes: []model.CIStr{model.NewCIStr("x")},
				},
				{
					HintName: model.NewCIStr("USE_INDEX_MERGE"),
					Tables: []ast.HintTable{{
						TableName:     model.NewCIStr("tbl2"),
						QBName:        model.NewCIStr("qb2"),
						PartitionList: []model.CIStr{model.NewCIStr("p0"), model.NewCIStr("p1")},
					}},
					QBName:  model.NewCIStr("qb2"),
					Indexes: []model.CIStr{model.NewCIStr("x"), model.NewCIStr("y"), model.NewCIStr("z")},
				},
			},
		},
		{
			input: `SET_VAR(sbs = 16M) SET_VAR(fkc=OFF) SET_VAR(os="mcb=off") set_var(abc=1) set_var(os2='mcb2=off')`,
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("SET_VAR"),
					HintData: ast.HintSetVar{
						VarName: "sbs",
						Value:   "16M",
					},
				},
				{
					HintName: model.NewCIStr("SET_VAR"),
					HintData: ast.HintSetVar{
						VarName: "fkc",
						Value:   "OFF",
					},
				},
				{
					HintName: model.NewCIStr("SET_VAR"),
					HintData: ast.HintSetVar{
						VarName: "os",
						Value:   "mcb=off",
					},
				},
				{
					HintName: model.NewCIStr("set_var"),
					HintData: ast.HintSetVar{
						VarName: "abc",
						Value:   "1",
					},
				},
				{
					HintName: model.NewCIStr("set_var"),
					HintData: ast.HintSetVar{
						VarName: "os2",
						Value:   "mcb2=off",
					},
				},
			},
		},
		{
			input: "USE_TOJA(TRUE) IGNORE_PLAN_CACHE() USE_CASCADES(TRUE) QUERY_TYPE(@qb1 OLAP) QUERY_TYPE(OLTP)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("USE_TOJA"),
					HintData: true,
				},
				{
					HintName: model.NewCIStr("IGNORE_PLAN_CACHE"),
				},
				{
					HintName: model.NewCIStr("USE_CASCADES"),
					HintData: true,
				},
				{
					HintName: model.NewCIStr("QUERY_TYPE"),
					QBName:   model.NewCIStr("qb1"),
					HintData: model.NewCIStr("OLAP"),
				},
				{
					HintName: model.NewCIStr("QUERY_TYPE"),
					HintData: model.NewCIStr("OLTP"),
				},
			},
		},
		{
			input: "READ_FROM_STORAGE(@foo TIKV[a, b], TIFLASH[c, d]) HASH_AGG() READ_FROM_STORAGE(TIKV[e])",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("READ_FROM_STORAGE"),
					HintData: model.NewCIStr("TIKV"),
					QBName:   model.NewCIStr("foo"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("a")},
						{TableName: model.NewCIStr("b")},
					},
				},
				{
					HintName: model.NewCIStr("READ_FROM_STORAGE"),
					HintData: model.NewCIStr("TIFLASH"),
					QBName:   model.NewCIStr("foo"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("c")},
						{TableName: model.NewCIStr("d")},
					},
				},
				{
					HintName: model.NewCIStr("HASH_AGG"),
				},
				{
					HintName: model.NewCIStr("READ_FROM_STORAGE"),
					HintData: model.NewCIStr("TIKV"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("e")},
					},
				},
			},
		},
		{
			input: "unknown_hint()",
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "set_var(timestamp = 1.5)",
			errs: []string{
				`Cannot use decimal number`,
				`Optimizer hint syntax error at line 1 `,
			},
		},
		{
			input: "set_var(timestamp = _utf8mb4'1234')", // Optimizer hint doesn't recognize _charset'strings'.
			errs:  []string{`Optimizer hint syntax error at line 1 `},
		},
		{
			input: "set_var(timestamp = 9999999999999999999999999999999999999)",
			errs: []string{
				`integer value is out of range`,
				`Optimizer hint syntax error at line 1 `,
			},
		},
		{
			input: "time_range('2020-02-20 12:12:12',456)",
			errs: []string{
				`Optimizer hint syntax error at line 1 `,
			},
		},
		{
			input: "time_range(456,'2020-02-20 12:12:12')",
			errs: []string{
				`Optimizer hint syntax error at line 1 `,
			},
		},
		{
			input: "TIME_RANGE('2020-02-20 12:12:12','2020-02-20 13:12:12')",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("TIME_RANGE"),
					HintData: ast.HintTimeRange{
						From: "2020-02-20 12:12:12",
						To:   "2020-02-20 13:12:12",
					},
				},
			},
		},
		{
			input: "JOIN_ORDER(@qb1 tbl1) JOIN_ORDER(@qb1 tbl1,tbl2) " +
				"JOIN_ORDER(tbl1) JOIN_ORDER(tbl1@qb1) JOIN_ORDER(tbl1@qb1, tbl2@qb2)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("JOIN_ORDER"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_ORDER"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_ORDER"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_ORDER"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_ORDER"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		{
			input: "JOIN_PREFIX(@qb1 tbl1) JOIN_PREFIX(@qb1 tbl1,tbl2) " +
				"JOIN_PREFIX(tbl1) JOIN_PREFIX(tbl1@qb1) JOIN_PREFIX(tbl1@qb1, tbl2@qb2)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("JOIN_PREFIX"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_PREFIX"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_PREFIX"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_PREFIX"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_PREFIX"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		{
			input: "JOIN_SUFFIX(@qb1 tbl1) JOIN_SUFFIX(@qb1 tbl1,tbl2) " +
				"JOIN_SUFFIX(tbl1) JOIN_SUFFIX(tbl1@qb1) JOIN_SUFFIX(tbl1@qb1, tbl2@qb2)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("JOIN_SUFFIX"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_SUFFIX"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_SUFFIX"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_SUFFIX"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("JOIN_SUFFIX"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		{
			input: "MAX_EXECUTION_TIME(1000)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("MAX_EXECUTION_TIME"),
					HintData: uint64(1000),
				},
			},
		},
		// MERGE,NOMERGE
		{
			input: "MERGE() MERGE(@qb1) MERGE(@qb1 tbl1) MERGE(@qb1 tbl1,tbl2) " +
				"MERGE(tbl1) MERGE(tbl1@qb1) MERGE(tbl1@qb1, tbl2@qb2)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("MERGE"),
				},
				{
					HintName: model.NewCIStr("MERGE"),
					QBName:   model.NewCIStr("qb1")},
				{
					HintName: model.NewCIStr("MERGE"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					QBName:   model.NewCIStr("qb1"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
						{TableName: model.NewCIStr("tbl2")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					Tables: []ast.HintTable{
						{TableName: model.NewCIStr("tbl1")},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
					},
				},
				{
					HintName: model.NewCIStr("MERGE"),
					Tables: []ast.HintTable{
						{
							TableName: model.NewCIStr("tbl1"),
							QBName:    model.NewCIStr("qb1"),
						},
						{
							TableName: model.NewCIStr("tbl2"),
							QBName:    model.NewCIStr("qb2"),
						},
					},
				},
			},
		},
		// MRR, NOMRR
		{
			input: "MRR(t1) MRR(t1 idx1,idx2) MRR(@qb1 t1 idx1,idx2) MRR(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("MRR"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("MRR"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("MRR"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("MRR"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		{
			input: "NO_ICP(t1) NO_ICP(t1 idx1,idx2) NO_ICP(@qb1 t1 idx1,idx2) NO_ICP(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		{
			input: "NO_ICP(t1) NO_ICP(t1 idx1,idx2) NO_ICP(@qb1 t1 idx1,idx2) NO_ICP(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_ICP"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		{
			input: "NO_RANGE_OPTIMIZATION(t1) NO_RANGE_OPTIMIZATION(t1 idx1,idx2) NO_RANGE_OPTIMIZATION(@qb1 t1 idx1,idx2)" +
				"NO_RANGE_OPTIMIZATION(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("NO_RANGE_OPTIMIZATION"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("NO_RANGE_OPTIMIZATION"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_RANGE_OPTIMIZATION"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_RANGE_OPTIMIZATION"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		// ORDER_INDEX, NO_ORDER_INDEX
		{
			input: "ORDER_INDEX(t1) ORDER_INDEX(t1 idx1,idx2) NO_ORDER_INDEX(@qb1 t1 idx1,idx2) NO_ORDER_INDEX(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("ORDER_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("ORDER_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_ORDER_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_ORDER_INDEX"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		//NO_SEMIJOIN, SEMIJOIN
		{
			input: "NO_SEMIJOIN(@subq1 FIRSTMATCH, LOOSESCAN) SEMIJOIN(@subq1 MATERIALIZATION, DUPSWEEDOUT) SEMIJOIN()",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("NO_SEMIJOIN"),
					QBName:   model.NewCIStr("subq1"),
					HintData: []model.CIStr{model.NewCIStr("FIRSTMATCH"), model.NewCIStr("LOOSESCAN")},
				},
				{
					HintName: model.NewCIStr("SEMIJOIN"),
					QBName:   model.NewCIStr("subq1"),
					HintData: []model.CIStr{model.NewCIStr("MATERIALIZATION"), model.NewCIStr("DUPSWEEDOUT")},
				},
				{
					HintName: model.NewCIStr("SEMIJOIN"),
					HintData: []model.CIStr{},
				},
			},
		},
		// SKIP_SCAN, NO_SKIP_SCAN
		{
			input: "SKIP_SCAN(t1) SKIP_SCAN(t1 idx1,idx2) NO_SKIP_SCAN(@qb1 t1 idx1,idx2) NO_SKIP_SCAN(t1@qb1 idx1,idx2) ",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("SKIP_SCAN"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
				},
				{
					HintName: model.NewCIStr("SKIP_SCAN"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
				{
					HintName: model.NewCIStr("NO_SKIP_SCAN"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
					QBName:   model.NewCIStr("qb1"),
				},
				{
					HintName: model.NewCIStr("NO_SKIP_SCAN"),
					Tables:   []ast.HintTable{{TableName: model.NewCIStr("t1"), QBName: model.NewCIStr("qb1")}},
					Indexes:  []model.CIStr{model.NewCIStr("idx1"), model.NewCIStr("idx2")},
				},
			},
		},
		{
			input: "RESOURCE_GROUP(group_name)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("RESOURCE_GROUP"),
					HintData: model.NewCIStr("group_name"),
				},
			},
		},
		// For SUBQUERY hints, these strategy values are permitted: INTOEXISTS, MATERIALIZATION.
		{
			input: "SUBQUERY(@subq1 MATERIALIZATION) SUBQUERY(@subq1 MATERIALIZATION, INTOEXISTS)",
			output: []*ast.TableOptimizerHint{
				{
					HintName: model.NewCIStr("SUBQUERY"),
					QBName:   model.NewCIStr("subq1"),
					HintData: []model.CIStr{model.NewCIStr("MATERIALIZATION")},
				},
				{
					HintName: model.NewCIStr("SUBQUERY"),
					QBName:   model.NewCIStr("subq1"),
					HintData: []model.CIStr{model.NewCIStr("MATERIALIZATION"), model.NewCIStr("INTOEXISTS")},
				},
			},
		},
	}

	for _, tc := range testCases {
		output, errs := parser.ParseHint("/*+"+tc.input+"*/", tc.mode, parser.Pos{Line: 1})
		require.Lenf(t, errs, len(tc.errs), "input = %s,\n... errs = %q", tc.input, errs)
		for i, err := range errs {
			require.Errorf(t, err, "input = %s, i = %d", tc.input, i)
			require.Containsf(t, err.Error(), tc.errs[i], "input = %s, i = %d", tc.input, i)
		}
		require.Equalf(t, tc.output, output, "input = %s,\n... output = %q", tc.input, output)
	}
}
