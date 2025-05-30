// Code generated by ent, DO NOT EDIT.

package singlepageinfo

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the singlepageinfo type in the database.
	Label = "single_page_info"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldBookID holds the string denoting the bookid field in the database.
	FieldBookID = "book_id"
	// FieldPageNum holds the string denoting the pagenum field in the database.
	FieldPageNum = "page_num"
	// FieldPath holds the string denoting the path field in the database.
	FieldPath = "path"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldBlurHash holds the string denoting the blurhash field in the database.
	FieldBlurHash = "blur_hash"
	// FieldHeight holds the string denoting the height field in the database.
	FieldHeight = "height"
	// FieldWidth holds the string denoting the width field in the database.
	FieldWidth = "width"
	// FieldModTime holds the string denoting the modtime field in the database.
	FieldModTime = "mod_time"
	// FieldSize holds the string denoting the size field in the database.
	FieldSize = "size"
	// FieldImgType holds the string denoting the imgtype field in the database.
	FieldImgType = "img_type"
	// Table holds the table name of the singlepageinfo in the database.
	Table = "single_page_infos"
)

// Columns holds all SQL columns for singlepageinfo fields.
var Columns = []string{
	FieldID,
	FieldBookID,
	FieldPageNum,
	FieldPath,
	FieldName,
	FieldURL,
	FieldBlurHash,
	FieldHeight,
	FieldWidth,
	FieldModTime,
	FieldSize,
	FieldImgType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "single_page_infos"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"book_page_infos",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultModTime holds the default value on creation for the "ModTime" field.
	DefaultModTime func() time.Time
)

// OrderOption defines the ordering options for the SinglePageInfo queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByBookID orders the results by the BookID field.
func ByBookID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBookID, opts...).ToFunc()
}

// ByPageNum orders the results by the PageNum field.
func ByPageNum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPageNum, opts...).ToFunc()
}

// ByPath orders the results by the Path field.
func ByPath(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPath, opts...).ToFunc()
}

// ByName orders the results by the Name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByURL orders the results by the Url field.
func ByURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURL, opts...).ToFunc()
}

// ByBlurHash orders the results by the BlurHash field.
func ByBlurHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBlurHash, opts...).ToFunc()
}

// ByHeight orders the results by the Height field.
func ByHeight(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHeight, opts...).ToFunc()
}

// ByWidth orders the results by the Width field.
func ByWidth(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWidth, opts...).ToFunc()
}

// ByModTime orders the results by the ModTime field.
func ByModTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModTime, opts...).ToFunc()
}

// BySize orders the results by the Size field.
func BySize(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSize, opts...).ToFunc()
}

// ByImgType orders the results by the ImgType field.
func ByImgType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImgType, opts...).ToFunc()
}
