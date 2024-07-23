package sl

import "log/slog"

func Op(op string) slog.Attr {
	return slog.String("op", op)
}

func Err(err error) slog.Attr {
	return slog.String("err", err.Error())
}
