// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// BooksColumns holds the columns for the "books" table.
	BooksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "title", Type: field.TypeString, Size: 1024},
		{Name: "book_id", Type: field.TypeString, Unique: true},
		{Name: "owner", Type: field.TypeString, Default: "admin"},
		{Name: "file_path", Type: field.TypeString},
		{Name: "book_store_path", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "child_book_num", Type: field.TypeInt},
		{Name: "depth", Type: field.TypeInt},
		{Name: "parent_folder", Type: field.TypeString},
		{Name: "page_count", Type: field.TypeInt},
		{Name: "size", Type: field.TypeInt64},
		{Name: "authors", Type: field.TypeString},
		{Name: "isbn", Type: field.TypeString},
		{Name: "press", Type: field.TypeString},
		{Name: "published_at", Type: field.TypeString},
		{Name: "extract_path", Type: field.TypeString},
		{Name: "modified", Type: field.TypeTime},
		{Name: "extract_num", Type: field.TypeInt},
		{Name: "init_complete", Type: field.TypeBool},
		{Name: "read_percent", Type: field.TypeFloat64},
		{Name: "non_utf8zip", Type: field.TypeBool},
		{Name: "zip_text_encoding", Type: field.TypeString},
	}
	// BooksTable holds the schema information for the "books" table.
	BooksTable = &schema.Table{
		Name:       "books",
		Columns:    BooksColumns,
		PrimaryKey: []*schema.Column{BooksColumns[0]},
	}
	// SinglePageInfosColumns holds the columns for the "single_page_infos" table.
	SinglePageInfosColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "book_id", Type: field.TypeString},
		{Name: "page_num", Type: field.TypeInt},
		{Name: "path", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "url", Type: field.TypeString},
		{Name: "blur_hash", Type: field.TypeString},
		{Name: "height", Type: field.TypeInt},
		{Name: "width", Type: field.TypeInt},
		{Name: "mod_time", Type: field.TypeTime},
		{Name: "size", Type: field.TypeInt64},
		{Name: "img_type", Type: field.TypeString},
		{Name: "book_page_infos", Type: field.TypeInt, Nullable: true},
	}
	// SinglePageInfosTable holds the schema information for the "single_page_infos" table.
	SinglePageInfosTable = &schema.Table{
		Name:       "single_page_infos",
		Columns:    SinglePageInfosColumns,
		PrimaryKey: []*schema.Column{SinglePageInfosColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "single_page_infos_books_PageInfos",
				Columns:    []*schema.Column{SinglePageInfosColumns[12]},
				RefColumns: []*schema.Column{BooksColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true, Size: 50},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "username", Type: field.TypeString, Unique: true, Size: 50},
		{Name: "password", Type: field.TypeString},
		{Name: "last_login", Type: field.TypeTime},
		{Name: "age", Type: field.TypeInt},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		BooksTable,
		SinglePageInfosTable,
		UsersTable,
	}
)

func init() {
	SinglePageInfosTable.ForeignKeys[0].RefTable = BooksTable
}
