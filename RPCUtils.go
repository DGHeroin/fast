package gf

import (
    "bytes"
    "context"
    "fmt"
    "log"
    "reflect"
    "runtime"
)

func MakeRPCFunc(name string, fn interface{}, r interface{}, w interface{}) (result func(ctx context.Context, req []byte) (resp []byte, err error), err error) {
    // 检查传入的函数是否符合格式要求
    var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
    f, ok := fn.(reflect.Value)
    if !ok {
        f = reflect.ValueOf(fn)
    }
    t := f.Type()
    if t.NumIn() != 3 { // context/request/response
        return nil, fmt.Errorf("func in param num error")
    }
    if t.NumOut() != 1 {
        return nil, fmt.Errorf("func out param num error")
    }
    if returnType := t.Out(0); returnType != typeOfError {
        return nil, fmt.Errorf("func out param type error")
    }
    // 生成 rpcx server 识别格式的函数
    result = func(ctx context.Context, req []byte) (resp []byte, err error) {
        defer func() {
            if p := recover(); p != nil {
                buffer := bytes.NewBufferString(fmt.Sprintf("%v\n", p))
                //打印调用栈信息
                buf := make([]byte, 2048)
                n := runtime.Stack(buf, false)
                stackInfo := fmt.Sprintf("%s", buf[:n])
                buffer.WriteString(fmt.Sprintf("panic stack info %s", stackInfo))
                log.Printf("RPC[%s] 请求失败:%v\n", name, buffer)
            }
        }()

        ctx = &RPCContext{
            values: map[string]interface{}{},
        }

        t0 := reflect.ValueOf(ctx)
        t1 := reflect.New(reflect.TypeOf(r))
        t2 := reflect.New(reflect.TypeOf(w))

        err = MSGUnpack(req, t1.Interface())
        if err != nil {
            return nil, err
        }
        in := []reflect.Value{
            t0, t1, t2,
        }
        rs := f.Call(in)
        r1 := rs[0]

        data, err := MSGPack(t2.Interface())
        if err != nil {
            return nil, err
        }

        if r1.Interface() == nil {
            return data, nil
        } else {
            err = r1.Interface().(error)
            return nil, err
        }
    }
    return result, nil
}
