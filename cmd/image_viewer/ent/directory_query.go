// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/directory"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/image"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent/predicate"
)

// DirectoryQuery is the builder for querying Directory entities.
type DirectoryQuery struct {
	config
	ctx          *QueryContext
	order        []directory.OrderOption
	inters       []Interceptor
	predicates   []predicate.Directory
	withParent   *DirectoryQuery
	withChildren *DirectoryQuery
	withImages   *ImageQuery
	withFKs      bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the DirectoryQuery builder.
func (dq *DirectoryQuery) Where(ps ...predicate.Directory) *DirectoryQuery {
	dq.predicates = append(dq.predicates, ps...)
	return dq
}

// Limit the number of records to be returned by this query.
func (dq *DirectoryQuery) Limit(limit int) *DirectoryQuery {
	dq.ctx.Limit = &limit
	return dq
}

// Offset to start from.
func (dq *DirectoryQuery) Offset(offset int) *DirectoryQuery {
	dq.ctx.Offset = &offset
	return dq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (dq *DirectoryQuery) Unique(unique bool) *DirectoryQuery {
	dq.ctx.Unique = &unique
	return dq
}

// Order specifies how the records should be ordered.
func (dq *DirectoryQuery) Order(o ...directory.OrderOption) *DirectoryQuery {
	dq.order = append(dq.order, o...)
	return dq
}

// QueryParent chains the current query on the "parent" edge.
func (dq *DirectoryQuery) QueryParent() *DirectoryQuery {
	query := (&DirectoryClient{config: dq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := dq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := dq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(directory.Table, directory.FieldID, selector),
			sqlgraph.To(directory.Table, directory.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, directory.ParentTable, directory.ParentColumn),
		)
		fromU = sqlgraph.SetNeighbors(dq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryChildren chains the current query on the "children" edge.
func (dq *DirectoryQuery) QueryChildren() *DirectoryQuery {
	query := (&DirectoryClient{config: dq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := dq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := dq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(directory.Table, directory.FieldID, selector),
			sqlgraph.To(directory.Table, directory.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, directory.ChildrenTable, directory.ChildrenColumn),
		)
		fromU = sqlgraph.SetNeighbors(dq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryImages chains the current query on the "images" edge.
func (dq *DirectoryQuery) QueryImages() *ImageQuery {
	query := (&ImageClient{config: dq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := dq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := dq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(directory.Table, directory.FieldID, selector),
			sqlgraph.To(image.Table, image.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, directory.ImagesTable, directory.ImagesColumn),
		)
		fromU = sqlgraph.SetNeighbors(dq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Directory entity from the query.
// Returns a *NotFoundError when no Directory was found.
func (dq *DirectoryQuery) First(ctx context.Context) (*Directory, error) {
	nodes, err := dq.Limit(1).All(setContextOp(ctx, dq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{directory.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (dq *DirectoryQuery) FirstX(ctx context.Context) *Directory {
	node, err := dq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Directory ID from the query.
// Returns a *NotFoundError when no Directory ID was found.
func (dq *DirectoryQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dq.Limit(1).IDs(setContextOp(ctx, dq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{directory.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (dq *DirectoryQuery) FirstIDX(ctx context.Context) int {
	id, err := dq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Directory entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Directory entity is found.
// Returns a *NotFoundError when no Directory entities are found.
func (dq *DirectoryQuery) Only(ctx context.Context) (*Directory, error) {
	nodes, err := dq.Limit(2).All(setContextOp(ctx, dq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{directory.Label}
	default:
		return nil, &NotSingularError{directory.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (dq *DirectoryQuery) OnlyX(ctx context.Context) *Directory {
	node, err := dq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Directory ID in the query.
// Returns a *NotSingularError when more than one Directory ID is found.
// Returns a *NotFoundError when no entities are found.
func (dq *DirectoryQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dq.Limit(2).IDs(setContextOp(ctx, dq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{directory.Label}
	default:
		err = &NotSingularError{directory.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (dq *DirectoryQuery) OnlyIDX(ctx context.Context) int {
	id, err := dq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Directories.
func (dq *DirectoryQuery) All(ctx context.Context) ([]*Directory, error) {
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryAll)
	if err := dq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Directory, *DirectoryQuery]()
	return withInterceptors[[]*Directory](ctx, dq, qr, dq.inters)
}

// AllX is like All, but panics if an error occurs.
func (dq *DirectoryQuery) AllX(ctx context.Context) []*Directory {
	nodes, err := dq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Directory IDs.
func (dq *DirectoryQuery) IDs(ctx context.Context) (ids []int, err error) {
	if dq.ctx.Unique == nil && dq.path != nil {
		dq.Unique(true)
	}
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryIDs)
	if err = dq.Select(directory.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (dq *DirectoryQuery) IDsX(ctx context.Context) []int {
	ids, err := dq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (dq *DirectoryQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryCount)
	if err := dq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, dq, querierCount[*DirectoryQuery](), dq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (dq *DirectoryQuery) CountX(ctx context.Context) int {
	count, err := dq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (dq *DirectoryQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryExist)
	switch _, err := dq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (dq *DirectoryQuery) ExistX(ctx context.Context) bool {
	exist, err := dq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the DirectoryQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (dq *DirectoryQuery) Clone() *DirectoryQuery {
	if dq == nil {
		return nil
	}
	return &DirectoryQuery{
		config:       dq.config,
		ctx:          dq.ctx.Clone(),
		order:        append([]directory.OrderOption{}, dq.order...),
		inters:       append([]Interceptor{}, dq.inters...),
		predicates:   append([]predicate.Directory{}, dq.predicates...),
		withParent:   dq.withParent.Clone(),
		withChildren: dq.withChildren.Clone(),
		withImages:   dq.withImages.Clone(),
		// clone intermediate query.
		sql:  dq.sql.Clone(),
		path: dq.path,
	}
}

// WithParent tells the query-builder to eager-load the nodes that are connected to
// the "parent" edge. The optional arguments are used to configure the query builder of the edge.
func (dq *DirectoryQuery) WithParent(opts ...func(*DirectoryQuery)) *DirectoryQuery {
	query := (&DirectoryClient{config: dq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	dq.withParent = query
	return dq
}

// WithChildren tells the query-builder to eager-load the nodes that are connected to
// the "children" edge. The optional arguments are used to configure the query builder of the edge.
func (dq *DirectoryQuery) WithChildren(opts ...func(*DirectoryQuery)) *DirectoryQuery {
	query := (&DirectoryClient{config: dq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	dq.withChildren = query
	return dq
}

// WithImages tells the query-builder to eager-load the nodes that are connected to
// the "images" edge. The optional arguments are used to configure the query builder of the edge.
func (dq *DirectoryQuery) WithImages(opts ...func(*ImageQuery)) *DirectoryQuery {
	query := (&ImageClient{config: dq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	dq.withImages = query
	return dq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Path string `json:"path,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Directory.Query().
//		GroupBy(directory.FieldPath).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (dq *DirectoryQuery) GroupBy(field string, fields ...string) *DirectoryGroupBy {
	dq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &DirectoryGroupBy{build: dq}
	grbuild.flds = &dq.ctx.Fields
	grbuild.label = directory.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Path string `json:"path,omitempty"`
//	}
//
//	client.Directory.Query().
//		Select(directory.FieldPath).
//		Scan(ctx, &v)
func (dq *DirectoryQuery) Select(fields ...string) *DirectorySelect {
	dq.ctx.Fields = append(dq.ctx.Fields, fields...)
	sbuild := &DirectorySelect{DirectoryQuery: dq}
	sbuild.label = directory.Label
	sbuild.flds, sbuild.scan = &dq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a DirectorySelect configured with the given aggregations.
func (dq *DirectoryQuery) Aggregate(fns ...AggregateFunc) *DirectorySelect {
	return dq.Select().Aggregate(fns...)
}

func (dq *DirectoryQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range dq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, dq); err != nil {
				return err
			}
		}
	}
	for _, f := range dq.ctx.Fields {
		if !directory.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if dq.path != nil {
		prev, err := dq.path(ctx)
		if err != nil {
			return err
		}
		dq.sql = prev
	}
	return nil
}

func (dq *DirectoryQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Directory, error) {
	var (
		nodes       = []*Directory{}
		withFKs     = dq.withFKs
		_spec       = dq.querySpec()
		loadedTypes = [3]bool{
			dq.withParent != nil,
			dq.withChildren != nil,
			dq.withImages != nil,
		}
	)
	if dq.withParent != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, directory.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Directory).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Directory{config: dq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, dq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := dq.withParent; query != nil {
		if err := dq.loadParent(ctx, query, nodes, nil,
			func(n *Directory, e *Directory) { n.Edges.Parent = e }); err != nil {
			return nil, err
		}
	}
	if query := dq.withChildren; query != nil {
		if err := dq.loadChildren(ctx, query, nodes,
			func(n *Directory) { n.Edges.Children = []*Directory{} },
			func(n *Directory, e *Directory) { n.Edges.Children = append(n.Edges.Children, e) }); err != nil {
			return nil, err
		}
	}
	if query := dq.withImages; query != nil {
		if err := dq.loadImages(ctx, query, nodes,
			func(n *Directory) { n.Edges.Images = []*Image{} },
			func(n *Directory, e *Image) { n.Edges.Images = append(n.Edges.Images, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (dq *DirectoryQuery) loadParent(ctx context.Context, query *DirectoryQuery, nodes []*Directory, init func(*Directory), assign func(*Directory, *Directory)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Directory)
	for i := range nodes {
		if nodes[i].directory_children == nil {
			continue
		}
		fk := *nodes[i].directory_children
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(directory.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "directory_children" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (dq *DirectoryQuery) loadChildren(ctx context.Context, query *DirectoryQuery, nodes []*Directory, init func(*Directory), assign func(*Directory, *Directory)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Directory)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Directory(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(directory.ChildrenColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.directory_children
		if fk == nil {
			return fmt.Errorf(`foreign-key "directory_children" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "directory_children" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (dq *DirectoryQuery) loadImages(ctx context.Context, query *ImageQuery, nodes []*Directory, init func(*Directory), assign func(*Directory, *Image)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Directory)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Image(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(directory.ImagesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.directory_images
		if fk == nil {
			return fmt.Errorf(`foreign-key "directory_images" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "directory_images" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (dq *DirectoryQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := dq.querySpec()
	_spec.Node.Columns = dq.ctx.Fields
	if len(dq.ctx.Fields) > 0 {
		_spec.Unique = dq.ctx.Unique != nil && *dq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, dq.driver, _spec)
}

func (dq *DirectoryQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(directory.Table, directory.Columns, sqlgraph.NewFieldSpec(directory.FieldID, field.TypeInt))
	_spec.From = dq.sql
	if unique := dq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if dq.path != nil {
		_spec.Unique = true
	}
	if fields := dq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, directory.FieldID)
		for i := range fields {
			if fields[i] != directory.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := dq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := dq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := dq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := dq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (dq *DirectoryQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(dq.driver.Dialect())
	t1 := builder.Table(directory.Table)
	columns := dq.ctx.Fields
	if len(columns) == 0 {
		columns = directory.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if dq.sql != nil {
		selector = dq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if dq.ctx.Unique != nil && *dq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range dq.predicates {
		p(selector)
	}
	for _, p := range dq.order {
		p(selector)
	}
	if offset := dq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := dq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// DirectoryGroupBy is the group-by builder for Directory entities.
type DirectoryGroupBy struct {
	selector
	build *DirectoryQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (dgb *DirectoryGroupBy) Aggregate(fns ...AggregateFunc) *DirectoryGroupBy {
	dgb.fns = append(dgb.fns, fns...)
	return dgb
}

// Scan applies the selector query and scans the result into the given value.
func (dgb *DirectoryGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, dgb.build.ctx, ent.OpQueryGroupBy)
	if err := dgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*DirectoryQuery, *DirectoryGroupBy](ctx, dgb.build, dgb, dgb.build.inters, v)
}

func (dgb *DirectoryGroupBy) sqlScan(ctx context.Context, root *DirectoryQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(dgb.fns))
	for _, fn := range dgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*dgb.flds)+len(dgb.fns))
		for _, f := range *dgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*dgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := dgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// DirectorySelect is the builder for selecting fields of Directory entities.
type DirectorySelect struct {
	*DirectoryQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ds *DirectorySelect) Aggregate(fns ...AggregateFunc) *DirectorySelect {
	ds.fns = append(ds.fns, fns...)
	return ds
}

// Scan applies the selector query and scans the result into the given value.
func (ds *DirectorySelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ds.ctx, ent.OpQuerySelect)
	if err := ds.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*DirectoryQuery, *DirectorySelect](ctx, ds.DirectoryQuery, ds, ds.inters, v)
}

func (ds *DirectorySelect) sqlScan(ctx context.Context, root *DirectoryQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ds.fns))
	for _, fn := range ds.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ds.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
