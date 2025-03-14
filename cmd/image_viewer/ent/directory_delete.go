// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/directory"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/predicate"
)

// DirectoryDelete is the builder for deleting a Directory entity.
type DirectoryDelete struct {
	config
	hooks    []Hook
	mutation *DirectoryMutation
}

// Where appends a list predicates to the DirectoryDelete builder.
func (dd *DirectoryDelete) Where(ps ...predicate.Directory) *DirectoryDelete {
	dd.mutation.Where(ps...)
	return dd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (dd *DirectoryDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, dd.sqlExec, dd.mutation, dd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (dd *DirectoryDelete) ExecX(ctx context.Context) int {
	n, err := dd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (dd *DirectoryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(directory.Table, sqlgraph.NewFieldSpec(directory.FieldID, field.TypeInt))
	if ps := dd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, dd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	dd.mutation.done = true
	return affected, err
}

// DirectoryDeleteOne is the builder for deleting a single Directory entity.
type DirectoryDeleteOne struct {
	dd *DirectoryDelete
}

// Where appends a list predicates to the DirectoryDelete builder.
func (ddo *DirectoryDeleteOne) Where(ps ...predicate.Directory) *DirectoryDeleteOne {
	ddo.dd.mutation.Where(ps...)
	return ddo
}

// Exec executes the deletion query.
func (ddo *DirectoryDeleteOne) Exec(ctx context.Context) error {
	n, err := ddo.dd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{directory.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ddo *DirectoryDeleteOne) ExecX(ctx context.Context) {
	if err := ddo.Exec(ctx); err != nil {
		panic(err)
	}
}
