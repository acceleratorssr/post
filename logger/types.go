package logger

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	With(args ...Field) Logger
}

type Field struct {
	Key   string
	Value any
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Int32(key string, val int32) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Int64(key string, val int64) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Bool(key string, b bool) Field {
	return Field{
		Key:   key,
		Value: b,
	}
}
