package tool

import (
	"CloudRestaurant/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DbEngine *Orm

type Orm struct {
	*xorm.Engine
}

func OrmEngine(cfg *Config) (*Orm, error) {
	dbCfg := cfg.Database
	//user:password@tcp(host:port)/dbname?charset= balaba
	source := dbCfg.User + ":" + dbCfg.Password + "@tcp(" + dbCfg.Host + ":" + dbCfg.Port + ")/" + dbCfg.DbName + "?charset=" + dbCfg.Charset
	engine, err := xorm.NewEngine(dbCfg.Driver, source)
	if err != nil {
		return nil, err
	}

	err = engine.Sync2(new(model.SmsCode), new(model.Member))
	if err != nil {
		return nil, err
	}

	orm := new(Orm)
	engine.ShowSQL(dbCfg.ShowSql)
	orm.Engine = engine

	DbEngine = orm

	return orm, nil
}
