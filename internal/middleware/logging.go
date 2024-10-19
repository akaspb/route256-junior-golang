package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"gitlab.ozon.dev/siralexpeter/Homework/internal/event_logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LocalLogging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	var logStr strings.Builder
	logStr.WriteString(fmt.Sprintf("[interceptor.Logging] method: %v", info.FullMethod))

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logStr.WriteString(fmt.Sprintf("; metadata: %v", md))
	}

	logStr.WriteString(fmt.Sprintf("; request: %v", req))

	res, err := handler(ctx, req)
	if err != nil {
		logStr.WriteString(fmt.Sprintf("; error: %v", err))
		log.Print(logStr.String())
		return
	}

	logStr.WriteString(fmt.Sprintf("; response: %v", res))
	log.Print(logStr.String())

	return res, nil
}

type EventDetails struct {
	Request  any    `json:"request"`
	Response any    `json:"response,omitempty"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

func GetRemoteLogging(
	eventLogger event_logger.EventLogger,
	eventFactory event_logger.EventFactory,
	chosenMethods []string,
) grpc.UnaryServerInterceptor {
	methods := make(map[string]struct{}, len(chosenMethods))
	for _, method := range chosenMethods {
		methods[method] = struct{}{}
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, methodErr error) {
		method := getMethodName(info.FullMethod)
		resp, methodErr = handler(ctx, req)

		if _, ok := methods[method]; !ok {
			return
		}

		var eventDetails EventDetails
		if methodErr == nil {
			eventDetails = EventDetails{
				Request:  req,
				Response: resp,
				Status:   "success",
			}
		} else {
			errString, err := responseErrorToString(methodErr)
			if err != nil {
				handleRemoteLoggerError(err)
				return
			}

			eventDetails = EventDetails{
				Request: req,
				Status:  "error",
				Error:   errString,
			}
		}

		event, err := eventFactory.Create(event_logger.EventType(method), eventDetails)
		if err != nil {
			handleRemoteLoggerError(err)
			return
		}

		err = eventLogger.Send(event)
		if err != nil {
			handleRemoteLoggerError(err)
			return
		}

		return
	}
}

func getMethodName(fullMethod string) string {
	return filepath.Base(fullMethod)
}

func responseErrorToString(err error) (string, error) {
	errStatus, ok := status.FromError(err)
	if !ok {
		return "", errors.New("responseErrorToString function should be used with status errors only")
	}

	return fmt.Sprintf("%v: %v", errStatus.Code(), errStatus.Message()), nil
}

func handleRemoteLoggerError(err error) {
	log.Printf("[kafka producer] error: %v", err)
}
