package responsecode

type ResponseCode string

const (

	// General
	CodeSuccess ResponseCode = "00"
	CodeFailed  ResponseCode = "XX"

	// Expected Error
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
	CodeOsError       ResponseCode = "OSERR"
	CodeTmuxError     ResponseCode = "TMXERR"
	CodeCentralError  ResponseCode = "CENTERR"
	CodeDockerError   ResponseCode = "DOCKERERR"

	// DB Error
	CodeTblServerTemplateError      ResponseCode = "TBLSTE"
	CodeTblUserError                ResponseCode = "TBLUSR"
	CodeTblInstanceError            ResponseCode = "TBLINT"
	CodeTblInstanceTypeError        ResponseCode = "TBLITT"
	CodeTblPersistentNodeError      ResponseCode = "TBLPND"
	CodeTblPersistentNodeImageError ResponseCode = "TBLPNI"
	CodeTblScriptError              ResponseCode = "TBLSCP"
	CodeTblStorageError             ResponseCode = "TBLSTO"
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
