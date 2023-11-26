// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/yumenaka/comi/ent/book"
	"github.com/yumenaka/comi/ent/schema"
	"github.com/yumenaka/comi/ent/singlepageinfo"
	"github.com/yumenaka/comi/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	bookFields := schema.Book{}.Fields()
	_ = bookFields
	// bookDescTitle is the schema descriptor for Title field.
	bookDescTitle := bookFields[0].Descriptor()
	// book.TitleValidator is a validator for the "Title" field. It is called by the builders before save.
	book.TitleValidator = bookDescTitle.Validators[0].(func(string) error)
	// bookDescOwner is the schema descriptor for Owner field.
	bookDescOwner := bookFields[2].Descriptor()
	// book.DefaultOwner holds the default value on creation for the Owner field.
	book.DefaultOwner = bookDescOwner.Default.(string)
	// bookDescChildBookNum is the schema descriptor for ChildBookNum field.
	bookDescChildBookNum := bookFields[6].Descriptor()
	// book.ChildBookNumValidator is a validator for the "ChildBookNum" field. It is called by the builders before save.
	book.ChildBookNumValidator = bookDescChildBookNum.Validators[0].(func(int) error)
	// bookDescDepth is the schema descriptor for Depth field.
	bookDescDepth := bookFields[7].Descriptor()
	// book.DepthValidator is a validator for the "Depth" field. It is called by the builders before save.
	book.DepthValidator = bookDescDepth.Validators[0].(func(int) error)
	// bookDescPageCount is the schema descriptor for PageCount field.
	bookDescPageCount := bookFields[9].Descriptor()
	// book.PageCountValidator is a validator for the "PageCount" field. It is called by the builders before save.
	book.PageCountValidator = bookDescPageCount.Validators[0].(func(int) error)
	// bookDescModified is the schema descriptor for Modified field.
	bookDescModified := bookFields[16].Descriptor()
	// book.DefaultModified holds the default value on creation for the Modified field.
	book.DefaultModified = bookDescModified.Default.(func() time.Time)
	singlepageinfoFields := schema.SinglePageInfo{}.Fields()
	_ = singlepageinfoFields
	// singlepageinfoDescModeTime is the schema descriptor for ModeTime field.
	singlepageinfoDescModeTime := singlepageinfoFields[7].Descriptor()
	// singlepageinfo.DefaultModeTime holds the default value on creation for the ModeTime field.
	singlepageinfo.DefaultModeTime = singlepageinfoDescModeTime.Default.(func() time.Time)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.NameValidator is a validator for the "name" field. It is called by the builders before save.
	user.NameValidator = userDescName.Validators[0].(func(string) error)
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userFields[1].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUsername is the schema descriptor for username field.
	userDescUsername := userFields[2].Descriptor()
	// user.UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	user.UsernameValidator = userDescUsername.Validators[0].(func(string) error)
	// userDescLastLogin is the schema descriptor for last_login field.
	userDescLastLogin := userFields[4].Descriptor()
	// user.DefaultLastLogin holds the default value on creation for the last_login field.
	user.DefaultLastLogin = userDescLastLogin.Default.(func() time.Time)
	// userDescAge is the schema descriptor for age field.
	userDescAge := userFields[5].Descriptor()
	// user.AgeValidator is a validator for the "age" field. It is called by the builders before save.
	user.AgeValidator = userDescAge.Validators[0].(func(int) error)
}
