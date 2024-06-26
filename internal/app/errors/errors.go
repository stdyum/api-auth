package errors

import (
	"net/http"

	"github.com/stdyum/api-auth/internal/app/controllers"
	grpcErr "github.com/stdyum/api-common/grpc"
	httpErr "github.com/stdyum/api-common/http"
	"google.golang.org/grpc/codes"
)

var (
	HttpErrorsMap = map[error]any{
		controllers.ErrUnauthorized: http.StatusUnauthorized,
	}

	GRpcErrorsMap = map[error]any{
		controllers.ErrUnauthorized: codes.Unauthenticated,
	}
)

func Register() {
	httpErr.RegisterErrors(HttpErrorsMap)
	grpcErr.RegisterErrors(GRpcErrorsMap)
}
