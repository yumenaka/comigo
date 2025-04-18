// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/directory"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/image"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/predicate"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeDirectory = "Directory"
	TypeImage     = "Image"
)

// DirectoryMutation represents an operation that mutates the Directory nodes in the graph.
type DirectoryMutation struct {
	config
	op              Op
	typ             string
	id              *int
	_path           *string
	name            *string
	clearedFields   map[string]struct{}
	parent          *int
	clearedparent   bool
	children        map[int]struct{}
	removedchildren map[int]struct{}
	clearedchildren bool
	images          map[int]struct{}
	removedimages   map[int]struct{}
	clearedimages   bool
	done            bool
	oldValue        func(context.Context) (*Directory, error)
	predicates      []predicate.Directory
}

var _ ent.Mutation = (*DirectoryMutation)(nil)

// directoryOption allows management of the mutation configuration using functional options.
type directoryOption func(*DirectoryMutation)

// newDirectoryMutation creates new mutation for the Directory entity.
func newDirectoryMutation(c config, op Op, opts ...directoryOption) *DirectoryMutation {
	m := &DirectoryMutation{
		config:        c,
		op:            op,
		typ:           TypeDirectory,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withDirectoryID sets the ID field of the mutation.
func withDirectoryID(id int) directoryOption {
	return func(m *DirectoryMutation) {
		var (
			err   error
			once  sync.Once
			value *Directory
		)
		m.oldValue = func(ctx context.Context) (*Directory, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Directory.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withDirectory sets the old Directory of the mutation.
func withDirectory(node *Directory) directoryOption {
	return func(m *DirectoryMutation) {
		m.oldValue = func(context.Context) (*Directory, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m DirectoryMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m DirectoryMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *DirectoryMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *DirectoryMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Directory.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetPath sets the "path" field.
func (m *DirectoryMutation) SetPath(s string) {
	m._path = &s
}

// Path returns the value of the "path" field in the mutation.
func (m *DirectoryMutation) Path() (r string, exists bool) {
	v := m._path
	if v == nil {
		return
	}
	return *v, true
}

// OldPath returns the old "path" field's value of the Directory entity.
// If the Directory object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *DirectoryMutation) OldPath(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPath is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPath requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPath: %w", err)
	}
	return oldValue.Path, nil
}

// ResetPath resets all changes to the "path" field.
func (m *DirectoryMutation) ResetPath() {
	m._path = nil
}

// SetName sets the "name" field.
func (m *DirectoryMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *DirectoryMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Directory entity.
// If the Directory object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *DirectoryMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *DirectoryMutation) ResetName() {
	m.name = nil
}

// SetParentID sets the "parent" edge to the Directory entity by id.
func (m *DirectoryMutation) SetParentID(id int) {
	m.parent = &id
}

// ClearParent clears the "parent" edge to the Directory entity.
func (m *DirectoryMutation) ClearParent() {
	m.clearedparent = true
}

// ParentCleared reports if the "parent" edge to the Directory entity was cleared.
func (m *DirectoryMutation) ParentCleared() bool {
	return m.clearedparent
}

// ParentID returns the "parent" edge ID in the mutation.
func (m *DirectoryMutation) ParentID() (id int, exists bool) {
	if m.parent != nil {
		return *m.parent, true
	}
	return
}

// ParentIDs returns the "parent" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ParentID instead. It exists only for internal usage by the builders.
func (m *DirectoryMutation) ParentIDs() (ids []int) {
	if id := m.parent; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetParent resets all changes to the "parent" edge.
func (m *DirectoryMutation) ResetParent() {
	m.parent = nil
	m.clearedparent = false
}

// AddChildIDs adds the "children" edge to the Directory entity by ids.
func (m *DirectoryMutation) AddChildIDs(ids ...int) {
	if m.children == nil {
		m.children = make(map[int]struct{})
	}
	for i := range ids {
		m.children[ids[i]] = struct{}{}
	}
}

// ClearChildren clears the "children" edge to the Directory entity.
func (m *DirectoryMutation) ClearChildren() {
	m.clearedchildren = true
}

// ChildrenCleared reports if the "children" edge to the Directory entity was cleared.
func (m *DirectoryMutation) ChildrenCleared() bool {
	return m.clearedchildren
}

// RemoveChildIDs removes the "children" edge to the Directory entity by IDs.
func (m *DirectoryMutation) RemoveChildIDs(ids ...int) {
	if m.removedchildren == nil {
		m.removedchildren = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.children, ids[i])
		m.removedchildren[ids[i]] = struct{}{}
	}
}

// RemovedChildren returns the removed IDs of the "children" edge to the Directory entity.
func (m *DirectoryMutation) RemovedChildrenIDs() (ids []int) {
	for id := range m.removedchildren {
		ids = append(ids, id)
	}
	return
}

// ChildrenIDs returns the "children" edge IDs in the mutation.
func (m *DirectoryMutation) ChildrenIDs() (ids []int) {
	for id := range m.children {
		ids = append(ids, id)
	}
	return
}

// ResetChildren resets all changes to the "children" edge.
func (m *DirectoryMutation) ResetChildren() {
	m.children = nil
	m.clearedchildren = false
	m.removedchildren = nil
}

// AddImageIDs adds the "images" edge to the Image entity by ids.
func (m *DirectoryMutation) AddImageIDs(ids ...int) {
	if m.images == nil {
		m.images = make(map[int]struct{})
	}
	for i := range ids {
		m.images[ids[i]] = struct{}{}
	}
}

// ClearImages clears the "images" edge to the Image entity.
func (m *DirectoryMutation) ClearImages() {
	m.clearedimages = true
}

// ImagesCleared reports if the "images" edge to the Image entity was cleared.
func (m *DirectoryMutation) ImagesCleared() bool {
	return m.clearedimages
}

// RemoveImageIDs removes the "images" edge to the Image entity by IDs.
func (m *DirectoryMutation) RemoveImageIDs(ids ...int) {
	if m.removedimages == nil {
		m.removedimages = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.images, ids[i])
		m.removedimages[ids[i]] = struct{}{}
	}
}

// RemovedImages returns the removed IDs of the "images" edge to the Image entity.
func (m *DirectoryMutation) RemovedImagesIDs() (ids []int) {
	for id := range m.removedimages {
		ids = append(ids, id)
	}
	return
}

// ImagesIDs returns the "images" edge IDs in the mutation.
func (m *DirectoryMutation) ImagesIDs() (ids []int) {
	for id := range m.images {
		ids = append(ids, id)
	}
	return
}

// ResetImages resets all changes to the "images" edge.
func (m *DirectoryMutation) ResetImages() {
	m.images = nil
	m.clearedimages = false
	m.removedimages = nil
}

// Where appends a list predicates to the DirectoryMutation builder.
func (m *DirectoryMutation) Where(ps ...predicate.Directory) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the DirectoryMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *DirectoryMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Directory, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *DirectoryMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *DirectoryMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Directory).
func (m *DirectoryMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *DirectoryMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m._path != nil {
		fields = append(fields, directory.FieldPath)
	}
	if m.name != nil {
		fields = append(fields, directory.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *DirectoryMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case directory.FieldPath:
		return m.Path()
	case directory.FieldName:
		return m.Name()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *DirectoryMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case directory.FieldPath:
		return m.OldPath(ctx)
	case directory.FieldName:
		return m.OldName(ctx)
	}
	return nil, fmt.Errorf("unknown Directory field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *DirectoryMutation) SetField(name string, value ent.Value) error {
	switch name {
	case directory.FieldPath:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPath(v)
		return nil
	case directory.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown Directory field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *DirectoryMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *DirectoryMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *DirectoryMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Directory numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *DirectoryMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *DirectoryMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *DirectoryMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Directory nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *DirectoryMutation) ResetField(name string) error {
	switch name {
	case directory.FieldPath:
		m.ResetPath()
		return nil
	case directory.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown Directory field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *DirectoryMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.parent != nil {
		edges = append(edges, directory.EdgeParent)
	}
	if m.children != nil {
		edges = append(edges, directory.EdgeChildren)
	}
	if m.images != nil {
		edges = append(edges, directory.EdgeImages)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *DirectoryMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case directory.EdgeParent:
		if id := m.parent; id != nil {
			return []ent.Value{*id}
		}
	case directory.EdgeChildren:
		ids := make([]ent.Value, 0, len(m.children))
		for id := range m.children {
			ids = append(ids, id)
		}
		return ids
	case directory.EdgeImages:
		ids := make([]ent.Value, 0, len(m.images))
		for id := range m.images {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *DirectoryMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedchildren != nil {
		edges = append(edges, directory.EdgeChildren)
	}
	if m.removedimages != nil {
		edges = append(edges, directory.EdgeImages)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *DirectoryMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case directory.EdgeChildren:
		ids := make([]ent.Value, 0, len(m.removedchildren))
		for id := range m.removedchildren {
			ids = append(ids, id)
		}
		return ids
	case directory.EdgeImages:
		ids := make([]ent.Value, 0, len(m.removedimages))
		for id := range m.removedimages {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *DirectoryMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	if m.clearedparent {
		edges = append(edges, directory.EdgeParent)
	}
	if m.clearedchildren {
		edges = append(edges, directory.EdgeChildren)
	}
	if m.clearedimages {
		edges = append(edges, directory.EdgeImages)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *DirectoryMutation) EdgeCleared(name string) bool {
	switch name {
	case directory.EdgeParent:
		return m.clearedparent
	case directory.EdgeChildren:
		return m.clearedchildren
	case directory.EdgeImages:
		return m.clearedimages
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *DirectoryMutation) ClearEdge(name string) error {
	switch name {
	case directory.EdgeParent:
		m.ClearParent()
		return nil
	}
	return fmt.Errorf("unknown Directory unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *DirectoryMutation) ResetEdge(name string) error {
	switch name {
	case directory.EdgeParent:
		m.ResetParent()
		return nil
	case directory.EdgeChildren:
		m.ResetChildren()
		return nil
	case directory.EdgeImages:
		m.ResetImages()
		return nil
	}
	return fmt.Errorf("unknown Directory edge %s", name)
}

// ImageMutation represents an operation that mutates the Image nodes in the graph.
type ImageMutation struct {
	config
	op               Op
	typ              string
	id               *int
	_path            *string
	name             *string
	size             *int64
	addsize          *int64
	mod_time         *time.Time
	create_time      *time.Time
	clearedFields    map[string]struct{}
	directory        *int
	cleareddirectory bool
	done             bool
	oldValue         func(context.Context) (*Image, error)
	predicates       []predicate.Image
}

var _ ent.Mutation = (*ImageMutation)(nil)

// imageOption allows management of the mutation configuration using functional options.
type imageOption func(*ImageMutation)

// newImageMutation creates new mutation for the Image entity.
func newImageMutation(c config, op Op, opts ...imageOption) *ImageMutation {
	m := &ImageMutation{
		config:        c,
		op:            op,
		typ:           TypeImage,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withImageID sets the ID field of the mutation.
func withImageID(id int) imageOption {
	return func(m *ImageMutation) {
		var (
			err   error
			once  sync.Once
			value *Image
		)
		m.oldValue = func(ctx context.Context) (*Image, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Image.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withImage sets the old Image of the mutation.
func withImage(node *Image) imageOption {
	return func(m *ImageMutation) {
		m.oldValue = func(context.Context) (*Image, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ImageMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ImageMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ImageMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ImageMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Image.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetPath sets the "path" field.
func (m *ImageMutation) SetPath(s string) {
	m._path = &s
}

// Path returns the value of the "path" field in the mutation.
func (m *ImageMutation) Path() (r string, exists bool) {
	v := m._path
	if v == nil {
		return
	}
	return *v, true
}

// OldPath returns the old "path" field's value of the Image entity.
// If the Image object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ImageMutation) OldPath(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPath is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPath requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPath: %w", err)
	}
	return oldValue.Path, nil
}

// ResetPath resets all changes to the "path" field.
func (m *ImageMutation) ResetPath() {
	m._path = nil
}

// SetName sets the "name" field.
func (m *ImageMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ImageMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Image entity.
// If the Image object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ImageMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *ImageMutation) ResetName() {
	m.name = nil
}

// SetSize sets the "size" field.
func (m *ImageMutation) SetSize(i int64) {
	m.size = &i
	m.addsize = nil
}

// Size returns the value of the "size" field in the mutation.
func (m *ImageMutation) Size() (r int64, exists bool) {
	v := m.size
	if v == nil {
		return
	}
	return *v, true
}

// OldSize returns the old "size" field's value of the Image entity.
// If the Image object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ImageMutation) OldSize(ctx context.Context) (v int64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSize is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSize requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSize: %w", err)
	}
	return oldValue.Size, nil
}

// AddSize adds i to the "size" field.
func (m *ImageMutation) AddSize(i int64) {
	if m.addsize != nil {
		*m.addsize += i
	} else {
		m.addsize = &i
	}
}

// AddedSize returns the value that was added to the "size" field in this mutation.
func (m *ImageMutation) AddedSize() (r int64, exists bool) {
	v := m.addsize
	if v == nil {
		return
	}
	return *v, true
}

// ResetSize resets all changes to the "size" field.
func (m *ImageMutation) ResetSize() {
	m.size = nil
	m.addsize = nil
}

// SetModTime sets the "mod_time" field.
func (m *ImageMutation) SetModTime(t time.Time) {
	m.mod_time = &t
}

// ModTime returns the value of the "mod_time" field in the mutation.
func (m *ImageMutation) ModTime() (r time.Time, exists bool) {
	v := m.mod_time
	if v == nil {
		return
	}
	return *v, true
}

// OldModTime returns the old "mod_time" field's value of the Image entity.
// If the Image object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ImageMutation) OldModTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldModTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldModTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldModTime: %w", err)
	}
	return oldValue.ModTime, nil
}

// ResetModTime resets all changes to the "mod_time" field.
func (m *ImageMutation) ResetModTime() {
	m.mod_time = nil
}

// SetCreateTime sets the "create_time" field.
func (m *ImageMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the value of the "create_time" field in the mutation.
func (m *ImageMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "create_time" field's value of the Image entity.
// If the Image object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ImageMutation) OldCreateTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "create_time" field.
func (m *ImageMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetDirectoryID sets the "directory" edge to the Directory entity by id.
func (m *ImageMutation) SetDirectoryID(id int) {
	m.directory = &id
}

// ClearDirectory clears the "directory" edge to the Directory entity.
func (m *ImageMutation) ClearDirectory() {
	m.cleareddirectory = true
}

// DirectoryCleared reports if the "directory" edge to the Directory entity was cleared.
func (m *ImageMutation) DirectoryCleared() bool {
	return m.cleareddirectory
}

// DirectoryID returns the "directory" edge ID in the mutation.
func (m *ImageMutation) DirectoryID() (id int, exists bool) {
	if m.directory != nil {
		return *m.directory, true
	}
	return
}

// DirectoryIDs returns the "directory" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// DirectoryID instead. It exists only for internal usage by the builders.
func (m *ImageMutation) DirectoryIDs() (ids []int) {
	if id := m.directory; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetDirectory resets all changes to the "directory" edge.
func (m *ImageMutation) ResetDirectory() {
	m.directory = nil
	m.cleareddirectory = false
}

// Where appends a list predicates to the ImageMutation builder.
func (m *ImageMutation) Where(ps ...predicate.Image) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ImageMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ImageMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Image, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ImageMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ImageMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Image).
func (m *ImageMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ImageMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m._path != nil {
		fields = append(fields, image.FieldPath)
	}
	if m.name != nil {
		fields = append(fields, image.FieldName)
	}
	if m.size != nil {
		fields = append(fields, image.FieldSize)
	}
	if m.mod_time != nil {
		fields = append(fields, image.FieldModTime)
	}
	if m.create_time != nil {
		fields = append(fields, image.FieldCreateTime)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ImageMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case image.FieldPath:
		return m.Path()
	case image.FieldName:
		return m.Name()
	case image.FieldSize:
		return m.Size()
	case image.FieldModTime:
		return m.ModTime()
	case image.FieldCreateTime:
		return m.CreateTime()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ImageMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case image.FieldPath:
		return m.OldPath(ctx)
	case image.FieldName:
		return m.OldName(ctx)
	case image.FieldSize:
		return m.OldSize(ctx)
	case image.FieldModTime:
		return m.OldModTime(ctx)
	case image.FieldCreateTime:
		return m.OldCreateTime(ctx)
	}
	return nil, fmt.Errorf("unknown Image field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ImageMutation) SetField(name string, value ent.Value) error {
	switch name {
	case image.FieldPath:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPath(v)
		return nil
	case image.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case image.FieldSize:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSize(v)
		return nil
	case image.FieldModTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModTime(v)
		return nil
	case image.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	}
	return fmt.Errorf("unknown Image field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ImageMutation) AddedFields() []string {
	var fields []string
	if m.addsize != nil {
		fields = append(fields, image.FieldSize)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ImageMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case image.FieldSize:
		return m.AddedSize()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ImageMutation) AddField(name string, value ent.Value) error {
	switch name {
	case image.FieldSize:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddSize(v)
		return nil
	}
	return fmt.Errorf("unknown Image numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ImageMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ImageMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ImageMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Image nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ImageMutation) ResetField(name string) error {
	switch name {
	case image.FieldPath:
		m.ResetPath()
		return nil
	case image.FieldName:
		m.ResetName()
		return nil
	case image.FieldSize:
		m.ResetSize()
		return nil
	case image.FieldModTime:
		m.ResetModTime()
		return nil
	case image.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	}
	return fmt.Errorf("unknown Image field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ImageMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.directory != nil {
		edges = append(edges, image.EdgeDirectory)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ImageMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case image.EdgeDirectory:
		if id := m.directory; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ImageMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ImageMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ImageMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.cleareddirectory {
		edges = append(edges, image.EdgeDirectory)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ImageMutation) EdgeCleared(name string) bool {
	switch name {
	case image.EdgeDirectory:
		return m.cleareddirectory
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ImageMutation) ClearEdge(name string) error {
	switch name {
	case image.EdgeDirectory:
		m.ClearDirectory()
		return nil
	}
	return fmt.Errorf("unknown Image unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ImageMutation) ResetEdge(name string) error {
	switch name {
	case image.EdgeDirectory:
		m.ResetDirectory()
		return nil
	}
	return fmt.Errorf("unknown Image edge %s", name)
}
