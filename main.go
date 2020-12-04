package main

import (
	"fmt"
	"reflect"
	"time"
	_ "upper.io/db.v3/mssql"
	"xorm.io/xorm"
)

//https://github.com/go-xorm/xorm

//https://gobook.io/read/gitea.com/xorm/manual-en-US/chapter-15/index.html
//https://gitea.com/xorm/xorm
//https://godoc.org/github.com/go-xorm/xorm
//https://unknwon.io/posts/140502_xorm-go-orm-basic-guide/

// define DataSourceName
const dsn = "server=dev-sql-05;user id=signit;password=signit;database=wolfgang.SQLTest2"
const driverName = "mssql"

func main() {
	engine, err := xorm.NewEngine(driverName, dsn)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer engine.Close()

	// enable logging on console
	engine.ShowSQL(false)

	dbSchema := []tableSchema{
		{tableStruct: new(User), foreignKeys: UserFK, data: UserData},
		{tableStruct: new(Role), foreignKeys: RoleFK, data: RoleData},
		{tableStruct: new(UserRole), foreignKeys: UserRoleFK, data: UserRoleData},
		{tableStruct: new(Documents), foreignKeys: DocumentsFK, data: DocumentsData}}
	if err = deleteSchema(engine, dbSchema...); err != nil {
		fmt.Printf("Error delete db schema: %v\n", err)
	}
	if err = addSchema(engine, dbSchema...); err != nil {
		fmt.Printf("Error add deb schema: %v\n", err)
	}
	if err = insertData(engine, dbSchema...); err != nil {
		fmt.Printf("Error insert data into tables: %v\n", err)
	}

	//  returns all tables schema information.
	{
		t, _ := engine.DBMetas()
		fmt.Printf("DBMetas: %v\n", t)
	}

	// SELECT * FROM user WHERE display_name="user3"
	{
		var u User
		exists, err := engine.Where("display_name=?", "user3").Get(&u)
		fmt.Printf("Get: exists: %v, User: %v, error: %v\n", exists, u, err)
	}

	// INSERT INTO [documents] ([id_user_modified_by],[name],[type],[content],[created_at],[updated_at]) OUTPUT Inserted.id
	// VALUES (?,?,?,?,?,?) [3 Document 0 [104 97 108 108 111] 2020-12-04 20:44:33 2020-12-04 20:44:33]
	//
	// UPDATE [documents] SET [name] = ?, [type] = ?, [updated_at] = ? WHERE [id]=? [blabla 1 2020-12-04 16:41:06 18]
	{
		// get UserId
		var moduser User
		_, err := engine.Where("display_name=?", "user3").Get(&moduser)

		// INSERT INTO [documents] ([id_user_modified_by],[name],[type],[content],[created_at],[updated_at]) OUTPUT Inserted.id
		// VALUES (?,?,?,?,?,?) [3 Document 0 [104 97 108 108 111] 2020-12-04 20:44:33 2020-12-04 20:44:33]
		d := Documents{
			Id_user_modified_by: moduser.Id,
			Name:                "Document",
			Type:                0,
			Content:             []byte("hallo"),
		}
		lines, err := engine.Insert(&d)
		fmt.Printf("document insert: lines:%v, error: %v\n", lines, err)

		// After inserted, documents.Id will be filled with primary key column value.
		id := d.Id

		// UPDATE
		// update only column Type where Id=id
		// Caution: update_at will be updated, too (autoupdate!)
		// UPDATE [documents] SET [name] = ?, [type] = ?, [updated_at] = ? WHERE [id]=? [blabla 1 2020-12-04 16:41:06 18]
		d.Type = 1
		d.Content = []byte("ein ganz neuer text")
		d.Name = "blabla"
		lines, err = engine.ID(id).Cols("type", "name").Update(&d)
		fmt.Printf("document update: lines:%v, error: %v\n", lines, err)
	}

	// INSERT INTO [user_role] ([id_user_modified_by],[id_user],[id_role],[created_at],[updated_at])
	// VALUES (?,?,?,?,?) [3 2 4 2020-12-04 20:43:37 2020-12-04 20:43:37]
	{
		var user_id int
		var usermod_id int
		var role_id int

		if _, err = engine.Table("user").Where("display_name = ?", "user3").And("is_active = 1").Cols("Id").Get(&usermod_id); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if _, err = engine.Table("user").Where("display_name = ?", "user2").And("is_active = 1").Cols("Id").Get(&user_id); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if _, err = engine.Table("role").Where("name = ?", "role4").Cols("Id").Get(&role_id); err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		lines, err := engine.Insert(&UserRole{
			Id_user_modified_by: usermod_id,
			Id_user:             user_id,
			Id_role:             role_id,
		})
		fmt.Printf("user_role insert: lines:%v, error: %v\n", lines, err)
	}

	// SELECT [user].*, [role].*,[user_role].created_at
	// FROM [user]
	// LEFT JOIN [user_role] ON user_role.id_user = [user].id
	// INNER JOIN [role] ON role.id = user_role.id_role
	// WHERE ([user].Display_name Like ?) [User%]
	{
		type UserDetail struct {
			User       `xorm:"extends"`
			Role       `xorm:"extends"`
			Created_at time.Time
		}

		var users []UserDetail
		err = engine.Table("user").Select("[user].*, [role].*,[user_role].created_at").
			Join("LEFT", "user_role", "user_role.id_user = [user].id").
			Join("INNER", "role", "role.id = user_role.id_role").
			Where("[user].Display_name Like ?", "User%").
			Find(&users)
		fmt.Printf("Find: UserDetail: %v, error: %v\n", users, err)
	}
}

func addSchema(engine *xorm.Engine, tables ...tableSchema) (err error) {
	// create Tables and extend table schema
	for _, t := range tables {
		if err = engine.Sync2(t.tableStruct); err != nil {
			return
		}
	}

	// Currently xorm hasn't support foreign key :(
	for _, t := range tables {
		for _, fk := range t.foreignKeys {
			if _, err = engine.Query(fk); err != nil {
				return
			}
		}
	}
	return
}

func deleteSchema(engine *xorm.Engine, tables ...tableSchema) (err error) {
	// the first table is the last to be deleted
	for i := len(tables) - 1; i >= 0; i-- {
		if ok, _ := engine.IsTableExist(tables[i].tableStruct); ok {
			if err = engine.DropTables(tables[i].tableStruct); err != nil {
				return
			}
		}
	}
	return nil
}

func insertData(engine *xorm.Engine, tables ...tableSchema) (err error) {
	// Insert Values into Tables
	for _, t := range tables {
		sliceValue := reflect.Indirect(reflect.ValueOf(t.data))
		if ok, _ := engine.IsTableExist(t.tableStruct); ok && sliceValue.Len() > 0 {
			if _, err = engine.Insert(t.data); err != nil {
				return
			}
		}
	}
	return nil
}

func examples(engine *xorm.Engine) {
	// query a SQL string, the returned results is []map[string][]byte
	results, err := engine.Query("select * from [user]")
	fmt.Printf("Query: results: %v, error: %v\n", results, err)

	affected, err := engine.Exec("select * from [user]")
	fmt.Printf("Query: affected: %v, error: %v\n", affected, err)

	// SELECT * FROM user LIMIT 1
	var u User
	has, err := engine.Get(&u)
	fmt.Printf("Get: has: %v, error: %v, user: %v\n", has, err, u)

	// SELECT * FROM user WHERE id=1
	u = User{}
	has, err = engine.ID(12).Get(&u)
	fmt.Printf("Get: has: %v, error: %v, user: %v\n", has, err, u)

	// SELECT * FROM user WHERE display_name="user1" AND id=1
	var list []User
	err = engine.Table("user").Where("display_name=? AND id=?", "user1", 1).Find(&list)
	fmt.Printf("Find: list: %v, error: %v\n", has, err)

}
