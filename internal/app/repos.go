package app

import (
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/repository/database"
	"cynxhostagent/internal/repository/database/tblinstance"
	"cynxhostagent/internal/repository/database/tblinstancetype"
	"cynxhostagent/internal/repository/database/tblparameter"
	"cynxhostagent/internal/repository/database/tblpersistentnode"
	"cynxhostagent/internal/repository/database/tblpersistentnodeimage"
	"cynxhostagent/internal/repository/database/tblscript"
	"cynxhostagent/internal/repository/database/tblservertemplate"
	"cynxhostagent/internal/repository/database/tblstorage"
	"cynxhostagent/internal/repository/database/tbluser"
)

type Repos struct {
	TblUser                database.TblUser
	TblScript              database.TblScript
	TblServerTemplate      database.TblServerTemplate
	TblInstance            database.TblInstance
	TblInstanceType        database.TblInstanceType
	TblPersistentNode      database.TblPersistentNode
	TblPersistentNodeImage database.TblPersistentNodeImage
	TblStorage             database.TblStorage
	TblParameter           database.TblParameter
	JWTManager             *dependencies.JWTManager
}

func NewRepos(dependencies *Dependencies) *Repos {

	return &Repos{
		TblUser:                tbluser.New(dependencies.DatabaseClient.Db),
		TblScript:              tblscript.New(dependencies.DatabaseClient.Db),
		TblServerTemplate:      tblservertemplate.New(dependencies.DatabaseClient.Db),
		TblInstance:            tblinstance.New(dependencies.DatabaseClient.Db),
		TblInstanceType:        tblinstancetype.New(dependencies.DatabaseClient.Db),
		TblPersistentNode:      tblpersistentnode.New(dependencies.DatabaseClient.Db),
		TblPersistentNodeImage: tblpersistentnodeimage.New(dependencies.DatabaseClient.Db),
		TblStorage:             tblstorage.New(dependencies.DatabaseClient.Db),
		TblParameter:           tblparameter.New(dependencies.DatabaseClient.Db),
		JWTManager:             dependencies.JWTManager,
	}
}
