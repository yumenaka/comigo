// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/yumenaka/comigo/internal/ent/predicate"
	"github.com/yumenaka/comigo/internal/ent/singlepageinfo"
)

// SinglePageInfoUpdate is the builder for updating SinglePageInfo entities.
type SinglePageInfoUpdate struct {
	config
	hooks    []Hook
	mutation *SinglePageInfoMutation
}

// Where appends a list predicates to the SinglePageInfoUpdate builder.
func (spiu *SinglePageInfoUpdate) Where(ps ...predicate.SinglePageInfo) *SinglePageInfoUpdate {
	spiu.mutation.Where(ps...)
	return spiu
}

// SetBookID sets the "BookID" field.
func (spiu *SinglePageInfoUpdate) SetBookID(s string) *SinglePageInfoUpdate {
	spiu.mutation.SetBookID(s)
	return spiu
}

// SetNillableBookID sets the "BookID" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableBookID(s *string) *SinglePageInfoUpdate {
	if s != nil {
		spiu.SetBookID(*s)
	}
	return spiu
}

// SetPageNum sets the "PageNum" field.
func (spiu *SinglePageInfoUpdate) SetPageNum(i int) *SinglePageInfoUpdate {
	spiu.mutation.ResetPageNum()
	spiu.mutation.SetPageNum(i)
	return spiu
}

// SetNillablePageNum sets the "PageNum" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillablePageNum(i *int) *SinglePageInfoUpdate {
	if i != nil {
		spiu.SetPageNum(*i)
	}
	return spiu
}

// AddPageNum adds i to the "PageNum" field.
func (spiu *SinglePageInfoUpdate) AddPageNum(i int) *SinglePageInfoUpdate {
	spiu.mutation.AddPageNum(i)
	return spiu
}

// SetPath sets the "Path" field.
func (spiu *SinglePageInfoUpdate) SetPath(s string) *SinglePageInfoUpdate {
	spiu.mutation.SetPath(s)
	return spiu
}

// SetNillablePath sets the "Path" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillablePath(s *string) *SinglePageInfoUpdate {
	if s != nil {
		spiu.SetPath(*s)
	}
	return spiu
}

// SetName sets the "Name" field.
func (spiu *SinglePageInfoUpdate) SetName(s string) *SinglePageInfoUpdate {
	spiu.mutation.SetName(s)
	return spiu
}

// SetNillableName sets the "Name" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableName(s *string) *SinglePageInfoUpdate {
	if s != nil {
		spiu.SetName(*s)
	}
	return spiu
}

// SetURL sets the "Url" field.
func (spiu *SinglePageInfoUpdate) SetURL(s string) *SinglePageInfoUpdate {
	spiu.mutation.SetURL(s)
	return spiu
}

// SetNillableURL sets the "Url" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableURL(s *string) *SinglePageInfoUpdate {
	if s != nil {
		spiu.SetURL(*s)
	}
	return spiu
}

// SetBlurHash sets the "BlurHash" field.
func (spiu *SinglePageInfoUpdate) SetBlurHash(s string) *SinglePageInfoUpdate {
	spiu.mutation.SetBlurHash(s)
	return spiu
}

// SetNillableBlurHash sets the "BlurHash" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableBlurHash(s *string) *SinglePageInfoUpdate {
	if s != nil {
		spiu.SetBlurHash(*s)
	}
	return spiu
}

// SetHeight sets the "Height" field.
func (spiu *SinglePageInfoUpdate) SetHeight(i int) *SinglePageInfoUpdate {
	spiu.mutation.ResetHeight()
	spiu.mutation.SetHeight(i)
	return spiu
}

// SetNillableHeight sets the "Height" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableHeight(i *int) *SinglePageInfoUpdate {
	if i != nil {
		spiu.SetHeight(*i)
	}
	return spiu
}

// AddHeight adds i to the "Height" field.
func (spiu *SinglePageInfoUpdate) AddHeight(i int) *SinglePageInfoUpdate {
	spiu.mutation.AddHeight(i)
	return spiu
}

// SetWidth sets the "Width" field.
func (spiu *SinglePageInfoUpdate) SetWidth(i int) *SinglePageInfoUpdate {
	spiu.mutation.ResetWidth()
	spiu.mutation.SetWidth(i)
	return spiu
}

// SetNillableWidth sets the "Width" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableWidth(i *int) *SinglePageInfoUpdate {
	if i != nil {
		spiu.SetWidth(*i)
	}
	return spiu
}

// AddWidth adds i to the "Width" field.
func (spiu *SinglePageInfoUpdate) AddWidth(i int) *SinglePageInfoUpdate {
	spiu.mutation.AddWidth(i)
	return spiu
}

// SetModTime sets the "ModTime" field.
func (spiu *SinglePageInfoUpdate) SetModTime(t time.Time) *SinglePageInfoUpdate {
	spiu.mutation.SetModTime(t)
	return spiu
}

// SetNillableModTime sets the "ModTime" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableModTime(t *time.Time) *SinglePageInfoUpdate {
	if t != nil {
		spiu.SetModTime(*t)
	}
	return spiu
}

// SetSize sets the "Size" field.
func (spiu *SinglePageInfoUpdate) SetSize(i int64) *SinglePageInfoUpdate {
	spiu.mutation.ResetSize()
	spiu.mutation.SetSize(i)
	return spiu
}

// SetNillableSize sets the "Size" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableSize(i *int64) *SinglePageInfoUpdate {
	if i != nil {
		spiu.SetSize(*i)
	}
	return spiu
}

// AddSize adds i to the "Size" field.
func (spiu *SinglePageInfoUpdate) AddSize(i int64) *SinglePageInfoUpdate {
	spiu.mutation.AddSize(i)
	return spiu
}

// SetImgType sets the "ImgType" field.
func (spiu *SinglePageInfoUpdate) SetImgType(s string) *SinglePageInfoUpdate {
	spiu.mutation.SetImgType(s)
	return spiu
}

// SetNillableImgType sets the "ImgType" field if the given value is not nil.
func (spiu *SinglePageInfoUpdate) SetNillableImgType(s *string) *SinglePageInfoUpdate {
	if s != nil {
		spiu.SetImgType(*s)
	}
	return spiu
}

// Mutation returns the SinglePageInfoMutation object of the builder.
func (spiu *SinglePageInfoUpdate) Mutation() *SinglePageInfoMutation {
	return spiu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (spiu *SinglePageInfoUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, spiu.sqlSave, spiu.mutation, spiu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (spiu *SinglePageInfoUpdate) SaveX(ctx context.Context) int {
	affected, err := spiu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (spiu *SinglePageInfoUpdate) Exec(ctx context.Context) error {
	_, err := spiu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (spiu *SinglePageInfoUpdate) ExecX(ctx context.Context) {
	if err := spiu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (spiu *SinglePageInfoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(singlepageinfo.Table, singlepageinfo.Columns, sqlgraph.NewFieldSpec(singlepageinfo.FieldID, field.TypeInt))
	if ps := spiu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := spiu.mutation.BookID(); ok {
		_spec.SetField(singlepageinfo.FieldBookID, field.TypeString, value)
	}
	if value, ok := spiu.mutation.PageNum(); ok {
		_spec.SetField(singlepageinfo.FieldPageNum, field.TypeInt, value)
	}
	if value, ok := spiu.mutation.AddedPageNum(); ok {
		_spec.AddField(singlepageinfo.FieldPageNum, field.TypeInt, value)
	}
	if value, ok := spiu.mutation.Path(); ok {
		_spec.SetField(singlepageinfo.FieldPath, field.TypeString, value)
	}
	if value, ok := spiu.mutation.Name(); ok {
		_spec.SetField(singlepageinfo.FieldName, field.TypeString, value)
	}
	if value, ok := spiu.mutation.URL(); ok {
		_spec.SetField(singlepageinfo.FieldURL, field.TypeString, value)
	}
	if value, ok := spiu.mutation.BlurHash(); ok {
		_spec.SetField(singlepageinfo.FieldBlurHash, field.TypeString, value)
	}
	if value, ok := spiu.mutation.Height(); ok {
		_spec.SetField(singlepageinfo.FieldHeight, field.TypeInt, value)
	}
	if value, ok := spiu.mutation.AddedHeight(); ok {
		_spec.AddField(singlepageinfo.FieldHeight, field.TypeInt, value)
	}
	if value, ok := spiu.mutation.Width(); ok {
		_spec.SetField(singlepageinfo.FieldWidth, field.TypeInt, value)
	}
	if value, ok := spiu.mutation.AddedWidth(); ok {
		_spec.AddField(singlepageinfo.FieldWidth, field.TypeInt, value)
	}
	if value, ok := spiu.mutation.ModTime(); ok {
		_spec.SetField(singlepageinfo.FieldModTime, field.TypeTime, value)
	}
	if value, ok := spiu.mutation.Size(); ok {
		_spec.SetField(singlepageinfo.FieldSize, field.TypeInt64, value)
	}
	if value, ok := spiu.mutation.AddedSize(); ok {
		_spec.AddField(singlepageinfo.FieldSize, field.TypeInt64, value)
	}
	if value, ok := spiu.mutation.ImgType(); ok {
		_spec.SetField(singlepageinfo.FieldImgType, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, spiu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{singlepageinfo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	spiu.mutation.done = true
	return n, nil
}

// SinglePageInfoUpdateOne is the builder for updating a single SinglePageInfo entity.
type SinglePageInfoUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SinglePageInfoMutation
}

// SetBookID sets the "BookID" field.
func (spiuo *SinglePageInfoUpdateOne) SetBookID(s string) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetBookID(s)
	return spiuo
}

// SetNillableBookID sets the "BookID" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableBookID(s *string) *SinglePageInfoUpdateOne {
	if s != nil {
		spiuo.SetBookID(*s)
	}
	return spiuo
}

// SetPageNum sets the "PageNum" field.
func (spiuo *SinglePageInfoUpdateOne) SetPageNum(i int) *SinglePageInfoUpdateOne {
	spiuo.mutation.ResetPageNum()
	spiuo.mutation.SetPageNum(i)
	return spiuo
}

// SetNillablePageNum sets the "PageNum" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillablePageNum(i *int) *SinglePageInfoUpdateOne {
	if i != nil {
		spiuo.SetPageNum(*i)
	}
	return spiuo
}

// AddPageNum adds i to the "PageNum" field.
func (spiuo *SinglePageInfoUpdateOne) AddPageNum(i int) *SinglePageInfoUpdateOne {
	spiuo.mutation.AddPageNum(i)
	return spiuo
}

// SetPath sets the "Path" field.
func (spiuo *SinglePageInfoUpdateOne) SetPath(s string) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetPath(s)
	return spiuo
}

// SetNillablePath sets the "Path" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillablePath(s *string) *SinglePageInfoUpdateOne {
	if s != nil {
		spiuo.SetPath(*s)
	}
	return spiuo
}

// SetName sets the "Name" field.
func (spiuo *SinglePageInfoUpdateOne) SetName(s string) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetName(s)
	return spiuo
}

// SetNillableName sets the "Name" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableName(s *string) *SinglePageInfoUpdateOne {
	if s != nil {
		spiuo.SetName(*s)
	}
	return spiuo
}

// SetURL sets the "Url" field.
func (spiuo *SinglePageInfoUpdateOne) SetURL(s string) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetURL(s)
	return spiuo
}

// SetNillableURL sets the "Url" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableURL(s *string) *SinglePageInfoUpdateOne {
	if s != nil {
		spiuo.SetURL(*s)
	}
	return spiuo
}

// SetBlurHash sets the "BlurHash" field.
func (spiuo *SinglePageInfoUpdateOne) SetBlurHash(s string) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetBlurHash(s)
	return spiuo
}

// SetNillableBlurHash sets the "BlurHash" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableBlurHash(s *string) *SinglePageInfoUpdateOne {
	if s != nil {
		spiuo.SetBlurHash(*s)
	}
	return spiuo
}

// SetHeight sets the "Height" field.
func (spiuo *SinglePageInfoUpdateOne) SetHeight(i int) *SinglePageInfoUpdateOne {
	spiuo.mutation.ResetHeight()
	spiuo.mutation.SetHeight(i)
	return spiuo
}

// SetNillableHeight sets the "Height" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableHeight(i *int) *SinglePageInfoUpdateOne {
	if i != nil {
		spiuo.SetHeight(*i)
	}
	return spiuo
}

// AddHeight adds i to the "Height" field.
func (spiuo *SinglePageInfoUpdateOne) AddHeight(i int) *SinglePageInfoUpdateOne {
	spiuo.mutation.AddHeight(i)
	return spiuo
}

// SetWidth sets the "Width" field.
func (spiuo *SinglePageInfoUpdateOne) SetWidth(i int) *SinglePageInfoUpdateOne {
	spiuo.mutation.ResetWidth()
	spiuo.mutation.SetWidth(i)
	return spiuo
}

// SetNillableWidth sets the "Width" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableWidth(i *int) *SinglePageInfoUpdateOne {
	if i != nil {
		spiuo.SetWidth(*i)
	}
	return spiuo
}

// AddWidth adds i to the "Width" field.
func (spiuo *SinglePageInfoUpdateOne) AddWidth(i int) *SinglePageInfoUpdateOne {
	spiuo.mutation.AddWidth(i)
	return spiuo
}

// SetModTime sets the "ModTime" field.
func (spiuo *SinglePageInfoUpdateOne) SetModTime(t time.Time) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetModTime(t)
	return spiuo
}

// SetNillableModTime sets the "ModTime" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableModTime(t *time.Time) *SinglePageInfoUpdateOne {
	if t != nil {
		spiuo.SetModTime(*t)
	}
	return spiuo
}

// SetSize sets the "Size" field.
func (spiuo *SinglePageInfoUpdateOne) SetSize(i int64) *SinglePageInfoUpdateOne {
	spiuo.mutation.ResetSize()
	spiuo.mutation.SetSize(i)
	return spiuo
}

// SetNillableSize sets the "Size" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableSize(i *int64) *SinglePageInfoUpdateOne {
	if i != nil {
		spiuo.SetSize(*i)
	}
	return spiuo
}

// AddSize adds i to the "Size" field.
func (spiuo *SinglePageInfoUpdateOne) AddSize(i int64) *SinglePageInfoUpdateOne {
	spiuo.mutation.AddSize(i)
	return spiuo
}

// SetImgType sets the "ImgType" field.
func (spiuo *SinglePageInfoUpdateOne) SetImgType(s string) *SinglePageInfoUpdateOne {
	spiuo.mutation.SetImgType(s)
	return spiuo
}

// SetNillableImgType sets the "ImgType" field if the given value is not nil.
func (spiuo *SinglePageInfoUpdateOne) SetNillableImgType(s *string) *SinglePageInfoUpdateOne {
	if s != nil {
		spiuo.SetImgType(*s)
	}
	return spiuo
}

// Mutation returns the SinglePageInfoMutation object of the builder.
func (spiuo *SinglePageInfoUpdateOne) Mutation() *SinglePageInfoMutation {
	return spiuo.mutation
}

// Where appends a list predicates to the SinglePageInfoUpdate builder.
func (spiuo *SinglePageInfoUpdateOne) Where(ps ...predicate.SinglePageInfo) *SinglePageInfoUpdateOne {
	spiuo.mutation.Where(ps...)
	return spiuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (spiuo *SinglePageInfoUpdateOne) Select(field string, fields ...string) *SinglePageInfoUpdateOne {
	spiuo.fields = append([]string{field}, fields...)
	return spiuo
}

// Save executes the query and returns the updated SinglePageInfo entity.
func (spiuo *SinglePageInfoUpdateOne) Save(ctx context.Context) (*SinglePageInfo, error) {
	return withHooks(ctx, spiuo.sqlSave, spiuo.mutation, spiuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (spiuo *SinglePageInfoUpdateOne) SaveX(ctx context.Context) *SinglePageInfo {
	node, err := spiuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (spiuo *SinglePageInfoUpdateOne) Exec(ctx context.Context) error {
	_, err := spiuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (spiuo *SinglePageInfoUpdateOne) ExecX(ctx context.Context) {
	if err := spiuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (spiuo *SinglePageInfoUpdateOne) sqlSave(ctx context.Context) (_node *SinglePageInfo, err error) {
	_spec := sqlgraph.NewUpdateSpec(singlepageinfo.Table, singlepageinfo.Columns, sqlgraph.NewFieldSpec(singlepageinfo.FieldID, field.TypeInt))
	id, ok := spiuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SinglePageInfo.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := spiuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, singlepageinfo.FieldID)
		for _, f := range fields {
			if !singlepageinfo.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != singlepageinfo.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := spiuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := spiuo.mutation.BookID(); ok {
		_spec.SetField(singlepageinfo.FieldBookID, field.TypeString, value)
	}
	if value, ok := spiuo.mutation.PageNum(); ok {
		_spec.SetField(singlepageinfo.FieldPageNum, field.TypeInt, value)
	}
	if value, ok := spiuo.mutation.AddedPageNum(); ok {
		_spec.AddField(singlepageinfo.FieldPageNum, field.TypeInt, value)
	}
	if value, ok := spiuo.mutation.Path(); ok {
		_spec.SetField(singlepageinfo.FieldPath, field.TypeString, value)
	}
	if value, ok := spiuo.mutation.Name(); ok {
		_spec.SetField(singlepageinfo.FieldName, field.TypeString, value)
	}
	if value, ok := spiuo.mutation.URL(); ok {
		_spec.SetField(singlepageinfo.FieldURL, field.TypeString, value)
	}
	if value, ok := spiuo.mutation.BlurHash(); ok {
		_spec.SetField(singlepageinfo.FieldBlurHash, field.TypeString, value)
	}
	if value, ok := spiuo.mutation.Height(); ok {
		_spec.SetField(singlepageinfo.FieldHeight, field.TypeInt, value)
	}
	if value, ok := spiuo.mutation.AddedHeight(); ok {
		_spec.AddField(singlepageinfo.FieldHeight, field.TypeInt, value)
	}
	if value, ok := spiuo.mutation.Width(); ok {
		_spec.SetField(singlepageinfo.FieldWidth, field.TypeInt, value)
	}
	if value, ok := spiuo.mutation.AddedWidth(); ok {
		_spec.AddField(singlepageinfo.FieldWidth, field.TypeInt, value)
	}
	if value, ok := spiuo.mutation.ModTime(); ok {
		_spec.SetField(singlepageinfo.FieldModTime, field.TypeTime, value)
	}
	if value, ok := spiuo.mutation.Size(); ok {
		_spec.SetField(singlepageinfo.FieldSize, field.TypeInt64, value)
	}
	if value, ok := spiuo.mutation.AddedSize(); ok {
		_spec.AddField(singlepageinfo.FieldSize, field.TypeInt64, value)
	}
	if value, ok := spiuo.mutation.ImgType(); ok {
		_spec.SetField(singlepageinfo.FieldImgType, field.TypeString, value)
	}
	_node = &SinglePageInfo{config: spiuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, spiuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{singlepageinfo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	spiuo.mutation.done = true
	return _node, nil
}
