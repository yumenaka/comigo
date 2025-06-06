// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/directory"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/image"
)

// DirectoryCreate is the builder for creating a Directory entity.
type DirectoryCreate struct {
	config
	mutation *DirectoryMutation
	hooks    []Hook
}

// SetPath sets the "path" field.
func (dc *DirectoryCreate) SetPath(s string) *DirectoryCreate {
	dc.mutation.SetPath(s)
	return dc
}

// SetName sets the "name" field.
func (dc *DirectoryCreate) SetName(s string) *DirectoryCreate {
	dc.mutation.SetName(s)
	return dc
}

// SetParentID sets the "parent" edge to the Directory entity by ID.
func (dc *DirectoryCreate) SetParentID(id int) *DirectoryCreate {
	dc.mutation.SetParentID(id)
	return dc
}

// SetNillableParentID sets the "parent" edge to the Directory entity by ID if the given value is not nil.
func (dc *DirectoryCreate) SetNillableParentID(id *int) *DirectoryCreate {
	if id != nil {
		dc = dc.SetParentID(*id)
	}
	return dc
}

// SetParent sets the "parent" edge to the Directory entity.
func (dc *DirectoryCreate) SetParent(d *Directory) *DirectoryCreate {
	return dc.SetParentID(d.ID)
}

// AddChildIDs adds the "children" edge to the Directory entity by IDs.
func (dc *DirectoryCreate) AddChildIDs(ids ...int) *DirectoryCreate {
	dc.mutation.AddChildIDs(ids...)
	return dc
}

// AddChildren adds the "children" edges to the Directory entity.
func (dc *DirectoryCreate) AddChildren(d ...*Directory) *DirectoryCreate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return dc.AddChildIDs(ids...)
}

// AddImageIDs adds the "images" edge to the Image entity by IDs.
func (dc *DirectoryCreate) AddImageIDs(ids ...int) *DirectoryCreate {
	dc.mutation.AddImageIDs(ids...)
	return dc
}

// AddImages adds the "images" edges to the Image entity.
func (dc *DirectoryCreate) AddImages(i ...*Image) *DirectoryCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return dc.AddImageIDs(ids...)
}

// Mutation returns the DirectoryMutation object of the builder.
func (dc *DirectoryCreate) Mutation() *DirectoryMutation {
	return dc.mutation
}

// Save creates the Directory in the database.
func (dc *DirectoryCreate) Save(ctx context.Context) (*Directory, error) {
	return withHooks(ctx, dc.sqlSave, dc.mutation, dc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (dc *DirectoryCreate) SaveX(ctx context.Context) *Directory {
	v, err := dc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dc *DirectoryCreate) Exec(ctx context.Context) error {
	_, err := dc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dc *DirectoryCreate) ExecX(ctx context.Context) {
	if err := dc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (dc *DirectoryCreate) check() error {
	if _, ok := dc.mutation.Path(); !ok {
		return &ValidationError{Name: "path", err: errors.New(`ent: missing required field "Directory.path"`)}
	}
	if v, ok := dc.mutation.Path(); ok {
		if err := directory.PathValidator(v); err != nil {
			return &ValidationError{Name: "path", err: fmt.Errorf(`ent: validator failed for field "Directory.path": %w`, err)}
		}
	}
	if _, ok := dc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Directory.name"`)}
	}
	if v, ok := dc.mutation.Name(); ok {
		if err := directory.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Directory.name": %w`, err)}
		}
	}
	return nil
}

func (dc *DirectoryCreate) sqlSave(ctx context.Context) (*Directory, error) {
	if err := dc.check(); err != nil {
		return nil, err
	}
	_node, _spec := dc.createSpec()
	if err := sqlgraph.CreateNode(ctx, dc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	dc.mutation.id = &_node.ID
	dc.mutation.done = true
	return _node, nil
}

func (dc *DirectoryCreate) createSpec() (*Directory, *sqlgraph.CreateSpec) {
	var (
		_node = &Directory{config: dc.config}
		_spec = sqlgraph.NewCreateSpec(directory.Table, sqlgraph.NewFieldSpec(directory.FieldID, field.TypeInt))
	)
	if value, ok := dc.mutation.Path(); ok {
		_spec.SetField(directory.FieldPath, field.TypeString, value)
		_node.Path = value
	}
	if value, ok := dc.mutation.Name(); ok {
		_spec.SetField(directory.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if nodes := dc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   directory.ParentTable,
			Columns: []string{directory.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(directory.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.directory_children = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := dc.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   directory.ChildrenTable,
			Columns: []string{directory.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(directory.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := dc.mutation.ImagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   directory.ImagesTable,
			Columns: []string{directory.ImagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(image.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// DirectoryCreateBulk is the builder for creating many Directory entities in bulk.
type DirectoryCreateBulk struct {
	config
	err      error
	builders []*DirectoryCreate
}

// Save creates the Directory entities in the database.
func (dcb *DirectoryCreateBulk) Save(ctx context.Context) ([]*Directory, error) {
	if dcb.err != nil {
		return nil, dcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(dcb.builders))
	nodes := make([]*Directory, len(dcb.builders))
	mutators := make([]Mutator, len(dcb.builders))
	for i := range dcb.builders {
		func(i int, root context.Context) {
			builder := dcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*DirectoryMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, dcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, dcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, dcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (dcb *DirectoryCreateBulk) SaveX(ctx context.Context) []*Directory {
	v, err := dcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dcb *DirectoryCreateBulk) Exec(ctx context.Context) error {
	_, err := dcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dcb *DirectoryCreateBulk) ExecX(ctx context.Context) {
	if err := dcb.Exec(ctx); err != nil {
		panic(err)
	}
}
