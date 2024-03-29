%{
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

package parser

import (
	"math"
	"strconv"

	"github.com/arana-db/parser/ast"
	"github.com/arana-db/parser/model"
)

%}

%union {
	offset  int
	ident   string
	number  uint64
	hint    *ast.TableOptimizerHint
	hints []*ast.TableOptimizerHint
	table 	ast.HintTable
	modelIdents []model.CIStr
}

%token	<number>

	/*yy:token "%d" */
	hintIntLit "a 64-bit unsigned integer"

%token	<ident>

	/*yy:token "%c" */
	hintIdentifier
	hintInvalid    "a special token never used by parser, used by lexer to indicate error"

	/*yy:token "@%c" */
	hintSingleAtIdentifier "identifier with single leading at"

	/*yy:token "'%c'" */
	hintStringLit

	/* MySQL 8.0 hint names */
	hintJoinFixedOrder             "JOIN_FIXED_ORDER"
	hintJoinOrder                  "JOIN_ORDER"
	hintJoinPrefix                 "JOIN_PREFIX"
	hintJoinSuffix                 "JOIN_SUFFIX"
	hintBKA                        "BKA"
	hintNoBKA                      "NO_BKA"
	hintBNL                        "BNL"
	hintNoBNL                      "NO_BNL"
	hintHashJoin                   "HASH_JOIN"
	hintNoHashJoin                 "NO_HASH_JOIN"
	hintMerge                      "MERGE"
	hintNoMerge                    "NO_MERGE"
	hintIndexMerge                 "INDEX_MERGE"
	hintNoIndexMerge               "NO_INDEX_MERGE"
	hintMRR                        "MRR"
	hintNoMRR                      "NO_MRR"
	hintNoICP                      "NO_ICP"
	hintNoRangeOptimization        "NO_RANGE_OPTIMIZATION"
	hintSkipScan                   "SKIP_SCAN"
	hintNoSkipScan                 "NO_SKIP_SCAN"
	hintSemijoin                   "SEMIJOIN"
	hintNoSemijoin                 "NO_SEMIJOIN"
	hintMaxExecutionTime           "MAX_EXECUTION_TIME"
	hintSetVar                     "SET_VAR"
	hintResourceGroup              "RESOURCE_GROUP"
	hintQBName                     "QB_NAME"
	hintDerivedConditionPushdown   "DERIVED_CONDITION_PUSHDOWN"
	hintNoDerivedConditionPushdown "NO_DERIVED_CONDITION_PUSHDOWN"
	hintGroupIndex                 "GROUP_INDEX"
	hintNoGroupIndex               "NO_GROUP_INDEX"
	hintIndex                      "INDEX"
	hintNoIndex                    "NO_INDEX"
	hintJoinIndex                  "JOIN_INDEX"
	hintNoJoinIndex                "NO_JOIN_INDEX"
	hintOrderIndex                 "ORDER_INDEX"
	hintNoOrderIndex               "NO_ORDER_INDEX"
	hintSubQuery                   "SUBQUERY"

	/* TiDB hint names */
	hintAggToCop              "AGG_TO_COP"
	hintIgnorePlanCache       "IGNORE_PLAN_CACHE"
	hintHashAgg               "HASH_AGG"
	hintIgnoreIndex           "IGNORE_INDEX"
	hintInlHashJoin           "INL_HASH_JOIN"
	hintInlJoin               "INL_JOIN"
	hintInlMergeJoin          "INL_MERGE_JOIN"
	hintMemoryQuota           "MEMORY_QUOTA"
	hintNoSwapJoinInputs      "NO_SWAP_JOIN_INPUTS"
	hintQueryType             "QUERY_TYPE"
	hintReadConsistentReplica "READ_CONSISTENT_REPLICA"
	hintReadFromStorage       "READ_FROM_STORAGE"
	hintSMJoin                "MERGE_JOIN"
	hintBCJoin                "BROADCAST_JOIN"
	hintBCJoinPreferLocal     "BROADCAST_JOIN_LOCAL"
	hintStreamAgg             "STREAM_AGG"
	hintSwapJoinInputs        "SWAP_JOIN_INPUTS"
	hintUseIndexMerge         "USE_INDEX_MERGE"
	hintUseIndex              "USE_INDEX"
	hintUsePlanCache          "USE_PLAN_CACHE"
	hintUseToja               "USE_TOJA"
	hintTimeRange             "TIME_RANGE"
	hintUseCascades           "USE_CASCADES"
	hintNthPlan               "NTH_PLAN"
	hintLimitToCop            "LIMIT_TO_COP"
	hintForceIndex            "FORCE_INDEX"

	/* Other keywords */
	hintOLAP            "OLAP"
	hintOLTP            "OLTP"
	hintPartition       "PARTITION"
	hintTiKV            "TIKV"
	hintTiFlash         "TIFLASH"
	hintFalse           "FALSE"
	hintTrue            "TRUE"
	hintMB              "MB"
	hintGB              "GB"
	hintDupsWeedOut     "DUPSWEEDOUT"
	hintFirstMatch      "FIRSTMATCH"
	hintLooseScan       "LOOSESCAN"
	hintMaterialization "MATERIALIZATION"
	hintIntoExist       "INTOEXISTS"

%type	<ident>
	Identifier                           "identifier (including keywords)"
	QueryBlockOpt                        "Query block identifier optional"
	JoinOrderOptimizerHintName
	SupportedTableLevelOptimizerHintName
	SupportedIndexLevelOptimizerHintName
	SubqueryOptimizerHintName
	BooleanHintName                      "name of hints which take a boolean input"
	NullaryHintName                      "name of hints which take no input"
	SubqueryStrategy
	Value                                "the value in the SET_VAR() hint"
	HintQueryType                        "query type in optimizer hint (OLAP or OLTP)"
	HintStorageType                      "storage type in optimizer hint (TiKV or TiFlash)"

%type	<number>
	UnitOfBytes "unit of bytes (MB or GB)"
	CommaOpt    "optional ','"

%type	<hints>
	OptimizerHintList           "optimizer hint list"
	StorageOptimizerHintOpt     "storage level optimizer hint"
	HintStorageTypeAndTableList "storage type and tables list in optimizer hint"

%type	<hint>
	TableOptimizerHintOpt   "optimizer hint"
	HintTableList           "table list in optimizer hint"
	HintTableListOpt        "optional table list in optimizer hint"
	HintIndexList           "table name with index list in optimizer hint"
	IndexNameList           "index list in optimizer hint"
	IndexNameListOpt        "optional index list in optimizer hint"
	HintTrueOrFalse         "true or false in optimizer hint"
	HintStorageTypeAndTable "storage type and tables in optimizer hint"

%type	<table>
	HintTable "Table in optimizer hint"

%type	<modelIdents>
	PartitionList         "partition name list in optimizer hint"
	PartitionListOpt      "optional partition name list in optimizer hint"
	SubqueryStrategies    "subquery strategies"
	SubqueryStrategiesOpt "optional subquery strategies"


%start	Start

%%

Start:
	OptimizerHintList
	{
		parser.result = $1
	}

OptimizerHintList:
	TableOptimizerHintOpt
	{
		if $1 != nil {
			$$ = []*ast.TableOptimizerHint{$1}
		}
	}
|	OptimizerHintList CommaOpt TableOptimizerHintOpt
	{
		if $3 != nil {
			$$ = append($1, $3)
		} else {
			$$ = $1
		}
	}
|	StorageOptimizerHintOpt
	{
		$$ = $1
	}
|	OptimizerHintList CommaOpt StorageOptimizerHintOpt
	{
		$$ = append($1, $3...)
	}

TableOptimizerHintOpt:
	JoinOrderOptimizerHintName '(' HintTableList ')'
	{
		h := $3
		h.HintName = model.NewCIStr($1)
		$$ = h
	}
|	SupportedTableLevelOptimizerHintName '(' HintTableListOpt ')'
	{
		h := $3
		h.HintName = model.NewCIStr($1)
		$$ = h
	}
|	SupportedIndexLevelOptimizerHintName '(' HintIndexList ')'
	{
		h := $3
		h.HintName = model.NewCIStr($1)
		$$ = h
	}
|	SubqueryOptimizerHintName '(' QueryBlockOpt SubqueryStrategiesOpt ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			QBName:   model.NewCIStr($3),
			HintData: $4,
		}
	}
|	"MAX_EXECUTION_TIME" '(' QueryBlockOpt hintIntLit ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			QBName:   model.NewCIStr($3),
			HintData: $4,
		}
	}
|	"NTH_PLAN" '(' QueryBlockOpt hintIntLit ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			QBName:   model.NewCIStr($3),
			HintData: int64($4),
		}
	}
|	"SET_VAR" '(' Identifier '=' Value ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			HintData: ast.HintSetVar{
				VarName: $3,
				Value:   $5,
			},
		}
	}
|	"RESOURCE_GROUP" '(' Identifier ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			HintData: model.NewCIStr($3),
		}
	}
|	"QB_NAME" '(' Identifier ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			QBName:   model.NewCIStr($3),
		}
	}
|	"MEMORY_QUOTA" '(' QueryBlockOpt hintIntLit UnitOfBytes ')'
	{
		maxValue := uint64(math.MaxInt64) / $5
		if $4 <= maxValue {
			$$ = &ast.TableOptimizerHint{
				HintName: model.NewCIStr($1),
				HintData: int64($4 * $5),
				QBName:   model.NewCIStr($3),
			}
		} else {
			yylex.AppendError(ErrWarnMemoryQuotaOverflow.GenWithStackByArgs(math.MaxInt64))
			parser.lastErrorAsWarn()
			$$ = nil
		}
	}
|	"TIME_RANGE" '(' hintStringLit CommaOpt hintStringLit ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			HintData: ast.HintTimeRange{
				From: $3,
				To:   $5,
			},
		}
	}
|	BooleanHintName '(' QueryBlockOpt HintTrueOrFalse ')'
	{
		h := $4
		h.HintName = model.NewCIStr($1)
		h.QBName = model.NewCIStr($3)
		$$ = h
	}
|	NullaryHintName '(' QueryBlockOpt ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			QBName:   model.NewCIStr($3),
		}
	}
|	"QUERY_TYPE" '(' QueryBlockOpt HintQueryType ')'
	{
		$$ = &ast.TableOptimizerHint{
			HintName: model.NewCIStr($1),
			QBName:   model.NewCIStr($3),
			HintData: model.NewCIStr($4),
		}
	}

StorageOptimizerHintOpt:
	"READ_FROM_STORAGE" '(' QueryBlockOpt HintStorageTypeAndTableList ')'
	{
		hs := $4
		name := model.NewCIStr($1)
		qb := model.NewCIStr($3)
		for _, h := range hs {
			h.HintName = name
			h.QBName = qb
		}
		$$ = hs
	}

HintStorageTypeAndTableList:
	HintStorageTypeAndTable
	{
		$$ = []*ast.TableOptimizerHint{$1}
	}
|	HintStorageTypeAndTableList ',' HintStorageTypeAndTable
	{
		$$ = append($1, $3)
	}

HintStorageTypeAndTable:
	HintStorageType '[' HintTableList ']'
	{
		h := $3
		h.HintData = model.NewCIStr($1)
		$$ = h
	}

QueryBlockOpt:
	/* empty */
	{
		$$ = ""
	}
|	hintSingleAtIdentifier

CommaOpt:
	/*empty*/
	{}
|	','
	{}

PartitionListOpt:
	/* empty */
	{
		$$ = nil
	}
|	"PARTITION" '(' PartitionList ')'
	{
		$$ = $3
	}

PartitionList:
	Identifier
	{
		$$ = []model.CIStr{model.NewCIStr($1)}
	}
|	PartitionList CommaOpt Identifier
	{
		$$ = append($1, model.NewCIStr($3))
	}

/**
 * HintTableListOpt:
 *
 *	[@query_block_name] [tbl_name [, tbl_name] ...]
 *	[tbl_name@query_block_name [, tbl_name@query_block_name] ...]
 *
 */
HintTableListOpt:
	HintTableList
|	QueryBlockOpt
	{
		$$ = &ast.TableOptimizerHint{
			QBName: model.NewCIStr($1),
		}
	}

HintTableList:
	QueryBlockOpt HintTable
	{
		$$ = &ast.TableOptimizerHint{
			Tables: []ast.HintTable{$2},
			QBName: model.NewCIStr($1),
		}
	}
|	HintTableList ',' HintTable
	{
		h := $1
		h.Tables = append(h.Tables, $3)
		$$ = h
	}

HintTable:
	Identifier QueryBlockOpt PartitionListOpt
	{
		$$ = ast.HintTable{
			TableName:     model.NewCIStr($1),
			QBName:        model.NewCIStr($2),
			PartitionList: $3,
		}
	}
|	Identifier '.' Identifier QueryBlockOpt PartitionListOpt
	{
		$$ = ast.HintTable{
			DBName:        model.NewCIStr($1),
			TableName:     model.NewCIStr($3),
			QBName:        model.NewCIStr($4),
			PartitionList: $5,
		}
	}

/**
 * HintIndexList:
 *
 *	[@query_block_name] tbl_name [index_name [, index_name] ...]
 *	tbl_name@query_block_name [index_name [, index_name] ...]
 */
HintIndexList:
	QueryBlockOpt HintTable CommaOpt IndexNameListOpt
	{
		h := $4
		h.Tables = []ast.HintTable{$2}
		h.QBName = model.NewCIStr($1)
		$$ = h
	}

IndexNameListOpt:
	/* empty */
	{
		$$ = &ast.TableOptimizerHint{}
	}
|	IndexNameList

IndexNameList:
	Identifier
	{
		$$ = &ast.TableOptimizerHint{
			Indexes: []model.CIStr{model.NewCIStr($1)},
		}
	}
|	IndexNameList ',' Identifier
	{
		h := $1
		h.Indexes = append(h.Indexes, model.NewCIStr($3))
		$$ = h
	}

/**
 * Miscellaneous rules
 */
SubqueryStrategiesOpt:
	/* empty */
	{
		$$ = []model.CIStr{}
	}
|	SubqueryStrategies

SubqueryStrategies:
	SubqueryStrategy
	{
		$$ = []model.CIStr{model.NewCIStr($1)}
	}
|	SubqueryStrategies CommaOpt SubqueryStrategy
	{
		$$ = append($1, model.NewCIStr($3))
	}

Value:
	hintStringLit
|	Identifier
|	hintIntLit
	{
		$$ = strconv.FormatUint($1, 10)
	}

UnitOfBytes:
	"MB"
	{
		$$ = 1024 * 1024
	}
|	"GB"
	{
		$$ = 1024 * 1024 * 1024
	}

HintTrueOrFalse:
	"TRUE"
	{
		$$ = &ast.TableOptimizerHint{HintData: true}
	}
|	"FALSE"
	{
		$$ = &ast.TableOptimizerHint{HintData: false}
	}

JoinOrderOptimizerHintName:
	"JOIN_ORDER"
|	"JOIN_PREFIX"
|	"JOIN_SUFFIX"

SupportedTableLevelOptimizerHintName:
	"MERGE_JOIN"
|	"BROADCAST_JOIN"
|	"BROADCAST_JOIN_LOCAL"
|	"INL_JOIN"
|	"INL_HASH_JOIN"
|	"SWAP_JOIN_INPUTS"
|	"NO_SWAP_JOIN_INPUTS"
|	"INL_MERGE_JOIN"
|	"HASH_JOIN"
|	"MERGE"
|	"NO_MERGE"
|	"BKA"
|	"NO_BKA"
|	"JOIN_FIXED_ORDER"
|	"BNL"
|	"NO_BNL"
|	"NO_HASH_JOIN"
|	"DERIVED_CONDITION_PUSHDOWN"
|	"NO_DERIVED_CONDITION_PUSHDOWN"

SupportedIndexLevelOptimizerHintName:
	"USE_INDEX"
|	"IGNORE_INDEX"
|	"USE_INDEX_MERGE"
|	"GROUP_INDEX"
|	"FORCE_INDEX"
|	"ORDER_INDEX"
|	"NO_ORDER_INDEX"
|	"SKIP_SCAN"
|	"NO_SKIP_SCAN"
|	"MRR"
|	"NO_MRR"
|	"NO_ICP"
|	"NO_RANGE_OPTIMIZATION"
|	"INDEX_MERGE"
|	"NO_INDEX_MERGE"
|	"NO_GROUP_INDEX"
|	"INDEX"
|	"NO_INDEX"
|	"JOIN_INDEX"
|	"NO_JOIN_INDEX"

SubqueryOptimizerHintName:
	"SEMIJOIN"
|	"NO_SEMIJOIN"
/*For SUBQUERY hints, permitted strategy values are different to SEMIJOIN and NO_SEMIJOIN, but not limit strictly here*/
|	"SUBQUERY"

SubqueryStrategy:
	"DUPSWEEDOUT"
|	"FIRSTMATCH"
|	"LOOSESCAN"
|	"MATERIALIZATION"
|	"INTOEXISTS"

BooleanHintName:
	"USE_TOJA"
|	"USE_CASCADES"

NullaryHintName:
	"USE_PLAN_CACHE"
|	"HASH_AGG"
|	"STREAM_AGG"
|	"AGG_TO_COP"
|	"LIMIT_TO_COP"
|	"READ_CONSISTENT_REPLICA"
|	"IGNORE_PLAN_CACHE"

HintQueryType:
	"OLAP"
|	"OLTP"

HintStorageType:
	"TIKV"
|	"TIFLASH"

Identifier:
	hintIdentifier
/* MySQL 8.0 hint names */
|	"JOIN_FIXED_ORDER"
|	"JOIN_ORDER"
|	"JOIN_PREFIX"
|	"JOIN_SUFFIX"
|	"BKA"
|	"NO_BKA"
|	"BNL"
|	"NO_BNL"
|	"HASH_JOIN"
|	"NO_HASH_JOIN"
|	"MERGE"
|	"NO_MERGE"
|	"INDEX_MERGE"
|	"NO_INDEX_MERGE"
|	"MRR"
|	"NO_MRR"
|	"NO_ICP"
|	"NO_RANGE_OPTIMIZATION"
|	"SKIP_SCAN"
|	"NO_SKIP_SCAN"
|	"SEMIJOIN"
|	"NO_SEMIJOIN"
|	"MAX_EXECUTION_TIME"
|	"SET_VAR"
|	"RESOURCE_GROUP"
|	"QB_NAME"
|	"DERIVED_CONDITION_PUSHDOWN"
|	"NO_DERIVED_CONDITION_PUSHDOWN"
|	"GROUP_INDEX"
|	"NO_GROUP_INDEX"
|	"INDEX"
|	"NO_INDEX"
|	"JOIN_INDEX"
|	"NO_JOIN_INDEX"
|	"ORDER_INDEX"
|	"NO_ORDER_INDEX"
|	"SUBQUERY"
/* TiDB hint names */
|	"AGG_TO_COP"
|	"LIMIT_TO_COP"
|	"IGNORE_PLAN_CACHE"
|	"HASH_AGG"
|	"IGNORE_INDEX"
|	"INL_HASH_JOIN"
|	"INL_JOIN"
|	"INL_MERGE_JOIN"
|	"MEMORY_QUOTA"
|	"NO_SWAP_JOIN_INPUTS"
|	"QUERY_TYPE"
|	"READ_CONSISTENT_REPLICA"
|	"READ_FROM_STORAGE"
|	"MERGE_JOIN"
|	"BROADCAST_JOIN"
|	"BROADCAST_JOIN_LOCAL"
|	"STREAM_AGG"
|	"SWAP_JOIN_INPUTS"
|	"USE_INDEX_MERGE"
|	"USE_INDEX"
|	"USE_PLAN_CACHE"
|	"USE_TOJA"
|	"TIME_RANGE"
|	"USE_CASCADES"
|	"NTH_PLAN"
|	"FORCE_INDEX"
/* other keywords */
|	"OLAP"
|	"OLTP"
|	"TIKV"
|	"TIFLASH"
|	"FALSE"
|	"TRUE"
|	"MB"
|	"GB"
|	"DUPSWEEDOUT"
|	"FIRSTMATCH"
|	"LOOSESCAN"
|	"MATERIALIZATION"
|	"INTOEXISTS"
%%
