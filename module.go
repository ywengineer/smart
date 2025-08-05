package smart

import (
	"fmt"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"github.com/gookit/event"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"reflect"
	"strconv"
)

// Module game logic module interface define
// logic method signature MethodFormat[(MessageName:string)(MessageCode:int)](context.Context, Channel, *Any) [*ResponseType]
// can get CtxKeySeq and CtxKeyHeader from logic method parameter context.Context
type Module interface {
	event.Listener
	Name() string
	Events() []string
}

var exclusiveMethod = []string{"Name", "Events", "Handle"}

func isExclusiveMethod(method string) bool {
	return lo.Contains(exclusiveMethod, method)
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
		mOutType := HandlerOutTypeNil
		// skip Module interface Name
		if isExclusiveMethod(mName) {
			continue
		}
		//
		if len(handlerMatched) < 2 { // method signature must contains route code
			logk.Debugf("not a request handler. method[%s] signature must be match regexp[%s]", mName, handlerSignatureRegexp.String())
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
		// outs
		nOut := mSignature.NumOut()
		if nOut > 0 {
			// 1st out must be *message.ProtocolMessage or proto.Message if num out equals 0
			ot := mSignature.Out(0)
			//
			if nOut == 1 {
				if ot == TypeSmartMessage {
					mOutType = HandlerOutTypeSmart
				} else {
					return createMethodSignatureError(mName)
				}
			} else {
				ot1 := mSignature.Out(1)
				//|| (!ot1.Implements(TypeProtoMessage) && ot1.Kind() != reflect.Slice)
				if ot.Kind() != reflect.Int {
					return createMethodSignatureError(mName)
				} else if ot1.Implements(TypeProtoMessage) {
					mOutType = HandlerOutTypeProtoMessage
				} else if ot1.Kind() == reflect.Slice {
					mOutType = HandlerOutTypeByteSlice
				} else {
					return createMethodSignatureError(mName)
				}
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
				outType:     mOutType,
			})
		}
	}
	//
	for _, e := range module.Events() {
		event.Listen(e, module)
	}
	return nil
}

func createMethodSignatureError(mName string) error {
	return fmt.Errorf("handler signature must be %s(context.Context, Channel, *Any) [ (int, []byte) | proto.Message ]", mName)
}
