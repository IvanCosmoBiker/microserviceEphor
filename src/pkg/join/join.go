package join

import (
 connectionPostgresql "ephorservices/src/pkg/connectDb"
 "context"
 "strings"
 "github.com/jackc/pgx/v4"
 "fmt"
 "log"
)

type Join struct {
    ConnectionDb connectionPostgresql.DatabaseInstance
    RowsData pgx.Rows
}

// func (j *Join) MakeMapResult() {

//     // fields := rows.FieldDescriptions()
//     // fmt.Printf("%s",fields[0].Name)
//     // data,_ := rows.Values()
//     // fmt.Printf("%+v",data)
// }

func (j *Join) GetJoin(accountId interface{},sql string) {
    sqlReplace := ""
    lookForSchema := "/schema/"
    schema := strings.Contains(sql, lookForSchema)
    accountSchema := fmt.Sprintf("account%v",accountId) 
    lookForMain := "/main/"
    main := strings.Contains(sql, lookForMain)
    if schema == true {
        sqlReplace = strings.Replace(sql, lookForSchema, accountSchema, -1)
    }else if main == true {
        sqlReplace = strings.Replace(sql, lookForMain, "main", -1)
    }
    ctx := context.Background()
    rows,err := j.ConnectionDb.Conn.Query(ctx,sqlReplace)
    if err != nil {
        log.Println(err) 
    }
    j.RowsData = rows
    // defer rows.Close()
}