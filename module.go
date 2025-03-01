package smart

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"reflect"
	"strconv"
)

// Module game logic module interface define
// logic method signature MethodFormat[(MessageName:string)(MessageCode:int)](context.Context, Channel, *Any) [*ResponseType]
// can get CtxKeySeq and CtxKeyHeader from logic method parameter context.Context
type Module interface {
	Name() string
}

func RegisterModule(module Module) error {
	mv := reflect.ValueOf(module)
	if mv.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("module must be a ptr and implements Module")
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
		handlerMatched := handlerSignatureRegexp.FindStringSubmatch(mName)
		if len(handlerMatched) < 2 { // method signature must contains route code
			utility.DefaultLogger().Debug("not a request handler", zap.String("method", mName), zap.String("regexp", handlerRegexp))
			continue
		}
		if mSignature.NumIn() != 3 {
			return createMethodSignatureError(mName)
		}
		// context.Context
		in0 := mSignature.In(0)
		// Channel
		in1 := mSignature.In(1)
		// error
		in2 := mSignature.In(2)
		if in0.Kind() != reflect.Interface || in0 != TypeContext ||
			in1.Kind() != reflect.Interface || in1 != TypeSocketChannel ||
			in2.Kind() != reflect.Ptr {
			return createMethodSignatureError(mName)
		}
		// 1st out must be ptr if num out greater than 0
		if mSignature.NumOut() > 0 {
			if mSignature.Out(0).Kind() != reflect.Ptr {
				return createMethodSignatureError(mName)
			}
		}
		s := handlerMatched[1]
		if code, err := strconv.Atoi(s); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("handler name must be match regexp[%s]", handlerRegexp))
		} else {
			hManager.addHandlerDefinition(&handlerDefinition{
				messageCode: code,
				name:        mName,
				inType:      in2,
				method:      method,
			})
		}
	}
	return nil
}

func createMethodSignatureError(mName string) error {
	return fmt.Errorf("handler signature must be %s(context.Context, Channel, *Any) [*ResponseType]", mName)
}
