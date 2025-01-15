package param

import (
	"context"
	"cynxhostagent/internal/repository/database"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type params struct {
}

var StaticParam params

type paramDetail struct {
	IsObject             bool
	ParamObjectReference interface{}
}

func getParamDetailList(staticParam *params) map[string]paramDetail {
	return map[string]paramDetail{}
}

func SetupStaticParam(tblParam database.TblParameter, log *logrus.Logger) {
	for {
		staticParam := &params{}

		ctx := context.Background()

		paramDetailList := getParamDetailList(staticParam)

		var idNames []string
		for key := range paramDetailList {
			idNames = append(idNames, key)
		}

		_, tblParams, err := tblParam.SelectTblParameters(ctx, idNames)
		if err != nil {
			log.Infoln("Error getting parameter: " + err.Error())
			continue
		}

		for _, tblParam := range tblParams {
			paramDetail := paramDetailList[tblParam.Id]
			if paramDetail.IsObject {
				err = json.Unmarshal([]byte(tblParam.Value), &paramDetail.ParamObjectReference)
				if err != nil {
					log.Infoln("Error unmarshalling parameter " + tblParam.Id + " : " + err.Error())
				}
				continue
			}

			if pObject, ok := paramDetail.ParamObjectReference.(*string); ok {
				*pObject = tblParam.Value
			}
		}

		StaticParam = *staticParam

		time.Sleep(30 * time.Minute)
	}
}
