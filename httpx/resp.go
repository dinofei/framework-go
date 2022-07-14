package httpx

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/dinofei/framework-go/errorx"
	"github.com/dinofei/framework-go/validatorx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gotomicro/ego/core/elog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	OK                int = 200
	GlobalInternalMsg     = "系统内部异常"
)

type (
	Body struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
	Option func(*Body)
)

func Code(code int) Option {
	return func(body *Body) {
		body.Code = code
	}
}

func Response(ctx *gin.Context, data interface{}, err error, opts ...Option) {
	b := &Body{}
	if err == nil {
		b.Code = OK
		if data == nil || reflect.ValueOf(data).IsZero() {
			data = []struct{}{}
		}
		b.Data = data
	} else {
		switch e := err.(type) {
		case validator.ValidationErrors:
			lan := GetClientLanguage(ctx.Request)
			b.Code = http.StatusBadRequest
			for _, v := range e {
				b.Message = v.Translate(validatorx.GetTranslator(lan))
				break
			}
		case *errorx.BizError:
			b.Code = e.Code
			b.Message = e.Message
		case interface{ GRPCStatus() *status.Status }:
			stats := e.GRPCStatus()
			b.Code = GrpcToHTTPStatusCode(stats.Code())
			if stats.Code() == codes.Unknown || stats.Code() == codes.Internal {
				b.Message = GlobalInternalMsg
				elog.Error("call rpc error", elog.FieldCustomKeyValue("path", ctx.Request.URL.Path), elog.FieldErr(stats.Err()))
			} else {
				b.Message = stats.Message()
			}
		default:
			b.Code = http.StatusInternalServerError
			b.Message = GlobalInternalMsg
			elog.Error("sys error", elog.FieldCustomKeyValue("path", ctx.Request.URL.Path), elog.FieldErr(e))
		}
	}

	for _, opt := range opts {
		opt(b)
	}

	ctx.JSON(b.Code, b)
}

// GrpcToHTTPStatusCode gRPC转HTTP Code
func GrpcToHTTPStatusCode(statusCode codes.Code) int {
	switch statusCode {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusServiceUnavailable
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func GetClientLanguage(r *http.Request) string {
	lan := r.Header.Get("Accept-Language")
	if strings.HasPrefix(lan, "zh") {
		lan = validatorx.ZH
	} else {
		lan = validatorx.EN
	}
	return lan
}
