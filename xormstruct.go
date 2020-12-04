package main

import "time"

type tableSchema struct {
	tableStruct interface{}
	foreignKeys []string
	data        interface{}
}

// int8,uint8,int16,uint16,int32,uint32,int64,uint64,int >> int
// int64,uint64 >> bigint
// time.Time datetime
// bool bit
// float32, float64 real
// rune int
// []byte varbinary(50)

// define Table User
type User struct {
	Id                  int `xorm:"not null unique pk autoincr"`
	Id_user_modified_by int
	Display_name        string `xorm:"nvarchar(255)"`
	First_name          string `xorm:"nvarchar(255)"`
	Last_name           string `xorm:"nvarchar(255)"`
	Mail                string `xorm:"nvarchar(255)"`
	Phone_number        string `xorm:"nvarchar(255)"`
	Language            string `xorm:"nvarchar(255)"`
	Password            string `xorm:"nvarchar(255)"`
	Extended            string `xorm:"nvarchar(255)"`
	Is_active           bool   `xorm:"not null"`
	Is_locked           bool   `xorm:"not null"`
	Confirmation_token  string `xorm:"nvarchar(255)"`
	Reset_token         string `xorm:"nvarchar(255)"`
	Login_count         int    `xorm:"not null"`
	Login_failed        int    `xorm:"not null"`
	Last_login          time.Time
	User_type           string    `xorm:"nvarchar(255)"`
	Created_at          time.Time `xorm:"not null created"`
	Updated_at          time.Time `xorm:"not null updated"`
}

var UserFK = []string{"ALTER TABLE [dbo].[user] ADD CONSTRAINT [FK_user_moduser] FOREIGN KEY([id_user_modified_by]) REFERENCES [dbo].[user]([id]) ON UPDATE NO ACTION ON DELETE NO ACTION"}

// define Table Role
type Role struct {
	Id                  int `xorm:"not null unique pk autoincr"`
	Id_user_modified_by int
	Name                string    `xorm:"nvarchar(255) unique index"`
	Description         string    `xorm:"nvarchar(255)"`
	Is_active           bool      `xorm:"not null"`
	Created_at          time.Time `xorm:"not null created"`
	Updated_at          time.Time `xorm:"not null updated"`
}

var RoleFK = []string{"ALTER TABLE [dbo].[role] ADD CONSTRAINT [FK_role_moduser] FOREIGN KEY([id_user_modified_by]) REFERENCES [dbo].[user]([id]) ON UPDATE NO ACTION ON DELETE SET NULL"}

// define Table UserRole
type UserRole struct {
	Id                  int `xorm:"not null unique pk autoincr"`
	Id_user_modified_by int
	Id_user             int       `xorm:"not null index"`
	Id_role             int       `xorm:"not null index"`
	Created_at          time.Time `xorm:"not null created"`
	Updated_at          time.Time `xorm:"not null updated"`
}

var UserRoleFK = []string{
	"ALTER TABLE [dbo].[user_role] ADD CONSTRAINT [FK_userrole_moduser] FOREIGN KEY([id_user_modified_by]) REFERENCES [dbo].[user]([id]) ON UPDATE NO ACTION ON DELETE NO ACTION",
	"ALTER TABLE [dbo].[user_role] ADD CONSTRAINT [FK_userrole_user] FOREIGN KEY([id_user]) REFERENCES [dbo].[user]([id]) ON UPDATE NO ACTION ON DELETE CASCADE",
	"ALTER TABLE [dbo].[user_role] ADD CONSTRAINT [FK_userrole_role] FOREIGN KEY([id_role]) REFERENCES [dbo].[role]([id]) ON UPDATE NO ACTION ON DELETE CASCADE"}

// define Table Documents
type Documents struct {
	Id                  int `xorm:"not null unique pk autoincr"`
	Id_user_modified_by int
	Name                string    `xorm:"nvarchar(255) unique index"`
	Type                int       `xorm:"not null default 4"`
	Content             []byte    `xorm:"varbinary(8000)"`
	Created_at          time.Time `xorm:"not null created"`
	Updated_at          time.Time `xorm:"not null updated"`
}

var DocumentsFK = []string{"ALTER TABLE [dbo].[documents] ADD CONSTRAINT [FK_documents_moduser] FOREIGN KEY([id_user_modified_by]) REFERENCES [dbo].[user]([id]) ON UPDATE NO ACTION ON DELETE SET NULL"}

var UserData = []User{
	{
		Id_user_modified_by: 1,
		Display_name:        "user1",
		First_name:          "",
		Last_name:           "",
		Mail:                "",
		Phone_number:        "",
		Language:            "",
		Password:            "",
		Extended:            "",
		Is_active:           true,
		Is_locked:           false,
		Confirmation_token:  "",
		Reset_token:         "",
		Login_count:         0,
		Login_failed:        0,
		Last_login:          time.Time{},
		User_type:           "",
	},
	{
		Id_user_modified_by: 1,
		Display_name:        "user2",
		First_name:          "",
		Last_name:           "",
		Mail:                "",
		Phone_number:        "",
		Language:            "",
		Password:            "",
		Extended:            "",
		Is_active:           true,
		Is_locked:           false,
		Confirmation_token:  "",
		Reset_token:         "",
		Login_count:         0,
		Login_failed:        0,
		Last_login:          time.Time{},
		User_type:           "",
	},
	{
		Id_user_modified_by: 1,
		Display_name:        "user3",
		First_name:          "",
		Last_name:           "",
		Mail:                "",
		Phone_number:        "",
		Language:            "",
		Password:            "",
		Extended:            "",
		Is_active:           true,
		Is_locked:           false,
		Confirmation_token:  "",
		Reset_token:         "",
		Login_count:         0,
		Login_failed:        0,
		Last_login:          time.Time{},
		User_type:           "",
	}}
var RoleData = []Role{
	{
		Id_user_modified_by: 1,
		Name:                "Role1",
	},
	{
		Id_user_modified_by: 1,
		Name:                "Role2",
	},
	{
		Id_user_modified_by: 1,
		Name:                "Role3",
	},
	{
		Id_user_modified_by: 1,
		Name:                "Role4",
	}}
var UserRoleData = []UserRole{}
var DocumentsData = []Documents{}
