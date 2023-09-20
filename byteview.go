package simplecache

//定义存储数据，get时返回slice的拷贝，防止修改底层数据

type byteview struct {
	b []byte
}

func cloneBytes(bytes []byte) []byte {
	cloneByte := make([]byte, len(bytes))
	copy(cloneByte, bytes)
	return cloneByte
}

func (v byteview) Len() int {
	return len(v.b)
}

func (v byteview) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v byteview) String() string {
	return string(v.b)
}
