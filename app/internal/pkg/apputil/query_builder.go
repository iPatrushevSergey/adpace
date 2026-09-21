package apputil

import (
	"errors"
	"strconv"
	"strings"
)

var ErrNoChanges = errors.New("no fields to patch")

// PatchQueryBuilder builds a partial UPDATE query, setting only the columns
// passed via Set with a non-nil value. One builder is single-use: create it,
// call Set for each optional field, then call Build once.
//
//	b := NewPatchQueryBuilder("advertiser", "advertiser_id", id, "advertiser_id, name, updated_at", 2)
//	Set(b, "name", patch.Name)
//	Set(b, "country", patch.Country)
//	query, args, err := b.Build()
type PatchQueryBuilder struct {
	table     string
	idCol     string
	idVal     any
	returning string
	query     strings.Builder
	args      []any
	tmpNum    [5]byte
	hasField  bool
}

func NewPatchQueryBuilder(table, idCol string, idVal any, returning string, maxFields int) *PatchQueryBuilder {
	b := &PatchQueryBuilder{table: table, idCol: idCol, idVal: idVal, returning: returning}
	b.query.Grow(34 + len(table) + len(idCol) + len(returning) + maxFields*24)
	b.args = make([]any, 0, maxFields+1)
	b.query.WriteString("UPDATE ")
	b.query.WriteString(table)
	b.query.WriteString(" SET ")
	return b
}

func (b *PatchQueryBuilder) writeParamNum(n int) {
	b.query.Write(strconv.AppendInt(b.tmpNum[:0], int64(n), 10))
}

func (b *PatchQueryBuilder) Build() (string, []any, error) {
	if !b.hasField {
		return "", nil, ErrNoChanges
	}

	b.args = append(b.args, b.idVal)
	b.query.WriteString(" WHERE ")
	b.query.WriteString(b.idCol)
	b.query.WriteString(" = $")
	b.writeParamNum(len(b.args))

	if b.returning != "" {
		b.query.WriteString(" RETURNING ")
		b.query.WriteString(b.returning)
	}

	return b.query.String(), b.args, nil
}

func Set[T any](b *PatchQueryBuilder, col string, val *T) *PatchQueryBuilder {
	if val == nil {
		return b
	}
	if b.hasField {
		b.query.WriteByte(',')
	}
	b.hasField = true
	b.args = append(b.args, *val)
	b.query.WriteString(col)
	b.query.WriteString(" = $")
	b.writeParamNum(len(b.args))
	return b
}
