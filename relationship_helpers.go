package main

import (
	"fmt"
	"strings"

	"github.com/nullbio/sqlboiler/bdb"
	"github.com/nullbio/sqlboiler/strmangle"
)

// RelationshipToManyTexts contains text that will be used by templates.
type RelationshipToManyTexts struct {
	LocalTable struct {
		NameGo       string
		NameSingular string
	}

	ForeignTable struct {
		NameGo            string
		NameSingular      string
		NamePluralGo      string
		NameHumanReadable string
		Slice             string
	}

	Function struct {
		Name     string
		Receiver string

		LocalAssignment   string
		ForeignAssignment string
	}
}

// createTextsFromRelationship creates a struct that does a lot of the text
// transformation in advance for a given relationship.
func createTextsFromRelationship(tables []bdb.Table, table bdb.Table, rel bdb.ToManyRelationship) RelationshipToManyTexts {
	r := RelationshipToManyTexts{}
	r.LocalTable.NameSingular = strmangle.Singular(table.Name)
	r.LocalTable.NameGo = strmangle.TitleCase(r.LocalTable.NameSingular)

	r.ForeignTable.NameSingular = strmangle.Singular(rel.ForeignTable)
	r.ForeignTable.NamePluralGo = strmangle.TitleCase(strmangle.Plural(rel.ForeignTable))
	r.ForeignTable.NameGo = strmangle.TitleCase(r.ForeignTable.NameSingular)
	r.ForeignTable.Slice = fmt.Sprintf("%sSlice", strmangle.CamelCase(r.ForeignTable.NameSingular))
	r.ForeignTable.NameHumanReadable = strings.Replace(rel.ForeignTable, "_", " ", -1)

	r.Function.Receiver = strings.ToLower(table.Name[:1])

	// Check to see if the foreign key name is the same as the local table name.
	// Simple case: yes - we can name the function the same as the plural table name
	// Not simple case: We have to name the function based off the foreign key and
	if colName := strings.TrimSuffix(rel.ForeignColumn, "_id"); r.LocalTable.NameSingular == colName {
		r.Function.Name = r.ForeignTable.NamePluralGo
	} else {
		r.Function.Name = r.ForeignTable.NamePluralGo + strmangle.TitleCase(colName)
	}

	if rel.Nullable {
		col := table.GetColumn(rel.Column)
		r.Function.LocalAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(rel.Column), strings.TrimPrefix(col.Type, "null."))
	} else {
		r.Function.LocalAssignment = strmangle.TitleCase(rel.Column)
	}

	if !rel.ToJoinTable {
		if rel.ForeignColumnNullable {
			foreignTable := bdb.GetTable(tables, rel.ForeignTable)
			col := foreignTable.GetColumn(rel.ForeignColumn)
			r.Function.ForeignAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(rel.ForeignColumn), strings.TrimPrefix(col.Type, "null."))
		} else {
			r.Function.ForeignAssignment = strmangle.TitleCase(rel.ForeignColumn)
		}

		return r
	}

	/*if rel.ForeignColumnNullable {
		foreignTable := bdb.GetTable(tables, rel.ForeignTable)
		col := foreignTable.GetColumn(rel.Column)
		r.Function.ForeignAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(rel.Column), strings.TrimPrefix("Null", col.Type))
	} else {
		r.Function.ForeignAssignment = strmangle.TitleCase(rel.ForeignColumn)
	}*/
	/*
		      {{if not .JoinTable -}}
		        {{- if .ForeignColumnNullable -}}
		          {{- $ftable := getTable $dot.Tables .ForeignTable -}}
		          {{- $fcol := getColumn $dot.Tables .ForeignColumn -}}
		  b.{{.ForeignColumn | camelCase}}.{{replace $fcol.Type "Null" ""}}, c.{{.ForeignColumn | camelCase}}.{{replace $fcol.Type "Null" ""}} = a.{{.Column}}{{if .Nullable}}, a.{{.Column}}
		  {{if .ForeignColumnNullable}}b.{{.ForeignColumn}}, c.{{.ForeignColumn}}{{else}}b.{{.ForeignColumn}}, c.{{.ForeignColumn}}{{end}} = a.ID, a.ID
		        {{- end -}}
		      {{- else -}}
		  b.user_id, c.user_id = a.ID, a.ID
		      {{- end -}}
		      {{- end}}
				// {{$fnName}}X retrieves all the {{$localTableSing}}'s {{$foreignTableHumanReadable}} with an executor.
				{{- if not $isForeignKeySimplyTableName}} via {{.ForeignColumn}} column.{{- end}}
				func ({{$receiver}} *{{$localTable}}) {{$fnName}}X(exec boil.Executor, selectCols ...string) ({{$foreignSlice}}, error) {
				  var ret {{$foreignSlice}}
	*/
	return r
}
