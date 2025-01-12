package responsecode

type ResponseCode string

const (

	// Expected Error
	CodeSuccess             ResponseCode = "SU"
	CodeValidationError     ResponseCode = "VE"
	CodeAuthenticationError ResponseCode = "AU"
	CodeNotAllowed          ResponseCode = "NA"
	CodeNotFound            ResponseCode = "NF"
	CodeInvalidCredentials  ResponseCode = "IC"
	CodeForbidden           ResponseCode = "FB"

	// Internal
	CodeJwtError      ResponseCode = "JWTERR"
	CodeInternalError ResponseCode = "IE"
	CodeAWSError      ResponseCode = "AWSERR"
	CodeEC2Error      ResponseCode = "EC2ERR"
	CodeRCONError     ResponseCode = "RCONERR"

	// DB Error
	CodeTblServerTemplateError ResponseCode = "TBLSTE"
	CodeTblUserError           ResponseCode = "TBLUSR"
	CodeTblInstanceError       ResponseCode = "TBLINT"
	CodeTblInstanceTypeError   ResponseCode = "TBLITT"
	CodeTblPersistentNodeError ResponseCode = "TBLPND"
	CodeTblScriptError         ResponseCode = "TBLSCP"
	CodeTblStorageError        ResponseCode = "TBLSTO"
)

var ResponseCodeNames = map[ResponseCode]string{
	CodeSuccess:             "Success",
	CodeValidationError:     "Validation Error",
	CodeAuthenticationError: "Authentication Error",
	CodeInternalError:       "Internal Error",
	CodeNotAllowed:          "Not Allowed",
	CodeNotFound:            "Not Found",

	CodeTblServerTemplateError: "TblServerTemplate Error",
	CodeTblUserError:           "TblUser Error",
	CodeTblInstanceError:       "TblInstance Error",
	CodeTblInstanceTypeError:   "TblInstanceType Error",
	CodeTblPersistentNodeError: "TblPersistentNode Error",
	CodeTblScriptError:         "TblScript Error",
	CodeTblStorageError:        "TblStorage Error",
}
