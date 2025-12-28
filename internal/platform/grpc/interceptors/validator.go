package interceptors

import (
	"context"
	"errors"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func MethodScopedValidatioUnaryInterceptor(
	validator protovalidate.Validator,
	enabled map[string]bool,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !enabled[info.FullMethod] {
			return handler(ctx, req)
		}

		msg, ok := req.(proto.Message)
		if !ok {
			return nil, status.Error(codes.Internal, "request does not implement proto.Message")
		}

		if err := validator.Validate(msg); err != nil {
			return nil, rpcValidationError(err)
		}

		return handler(ctx, req)
	}
}

func rpcValidationError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if err == nil {
			return st.Err()
		}
	}

	return status.Error(codes.Internal, err.Error())
}

func convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateVialationsToGoogleViolations(valErr.ToProto().GetViolations()),
	}
}

func protovalidateVialationsToGoogleViolations(violations []*validate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(violations))
	for i, violation := range violations {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       violation.GetField().String(),
			Description: violation.GetMessage(),
		}
	}

	return res
}
