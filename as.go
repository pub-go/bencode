package bencode

func AsStr(v Value) string {
	return string(AsString(v))
}

func AsInt(v Value) int64 {
	return int64(AsInteger(v))
}

func AsString(v Value) String {
	if l, ok := v.(String); ok {
		return l
	}
	return ""
}

func AsInteger(v Value) Integer {
	if l, ok := v.(Integer); ok {
		return l
	}
	return 0
}

func AsList(v Value) List {
	if l, ok := v.(List); ok {
		return l
	}
	return List{}
}

func AsDict(v Value) Dict {
	if l, ok := v.(Dict); ok {
		return l
	}
	return Dict{}
}
