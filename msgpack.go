package fast

import (
    "github.com/vmihailenco/msgpack"
)

func MSGPack(obj interface{}) ([]byte, error) {
    return msgpack.Marshal(obj)
}

func MSGUnpack(data []byte, obj interface{}) error {
    return msgpack.Unmarshal(data, obj)
}
