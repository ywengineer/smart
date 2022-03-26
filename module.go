package mr_smart

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
	"strconv"
	"strings"
)

type SmartModule interface {
	Name() string
}

func RegisterModule(module SmartModule) error {
	mv := reflect.ValueOf(module)
	if mv.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("module must be a ptr and implements SmartModule")
	}
	methods := mv.NumMethod()
	if methods == 0 {
		return fmt.Errorf("module[%s] has not any method exposed", mv.Type().Name())
	}
	// generate handler definition
	for mIndex := 0; mIndex < methods; mIndex++ {
		method := mv.Method(mIndex)
		mName := mv.Type().Method(mIndex).Name
		mSignature := method.Type()
		if !strings.HasPrefix(mName, HandlerPrefix) {
			srvLogger.Debug("not a request handler", zap.String("method", mName))
			continue
		}
		if mSignature.NumIn() != 2 {
			return createMethodSignatureError(mName)
		}
		// *SocketChannel
		in0 := mSignature.In(0)
		// error
		in1 := mSignature.In(1)
		if in0.Kind() != reflect.Ptr || in0 != TypeSocketChannel || in1.Kind() != reflect.Ptr {
			return createMethodSignatureError(mName)
		}
		// 1th out must be ptr if num out greater than 0
		if mSignature.NumOut() > 0 {
			if mSignature.Out(0).Kind() != reflect.Ptr {
				return createMethodSignatureError(mName)
			}
		}
		s := mName[HandlerPrefixLength:]
		if code, err := strconv.Atoi(s); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("handler name must be %s[1-9][0-9]*", HandlerPrefix))
		} else {
			addHandlerDefinition(&handlerDefinition{
				messageCode: code,
				name:        mName,
				inType:      in1,
				method:      method,
			})
		}
	}
	return nil
}

func createMethodSignatureError(mName string) error {
	return fmt.Errorf("handler[%s] signature must be Handler[1-9][0-9]*(*SocketChannel, *RequestType) [*ResponseType]", mName)
}
