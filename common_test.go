package commontools

import (
	"github.com/jonluo94/commontools/log"
	"testing"
	"github.com/jonluo94/commontools/xorm"
	"github.com/jonluo94/commontools/password"
	"fmt"
)


func TestLog(test *testing.T) {
	logger := log.GetLogger("test",log.DEBUG)
	logger.Debugf("debug %s", log.Secret("secret"))
	logger.Info("info")
	logger.Notice("notice")
	logger.Warning("warning")
	logger.Error("err")
	logger.Critical("crit")
}

func TestXorm(test *testing.T) {

	//docker run -p 3306:3306 --name some-mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
	type User struct{
		id string `json:"id" xorm:"varchar(64) pk not null"`
		name string `json:"name" xorm:"varchar(64) not null"`
	}

	user := &User{}
	engine := xorm.GetEngine("xorm/config.yaml")
	engine.CreateTables(user)
}


func TestPassword(test *testing.T) {

	cry := password.Encode("admin",12,"default")
    fmt.Println(cry)
	b := password.Validate("admin",cry)
	fmt.Println(b)
}