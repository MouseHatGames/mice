package auth

import "context"

type keyType int

const keyUserID keyType = 1

func WithUserID(ctx context.Context, u uint32) context.Context {
	return context.WithValue(ctx, keyUserID, u)
}

func GetUserID(ctx context.Context) (uint32, bool) {
	v := ctx.Value(keyUserID)
	if v == nil {
		return 0, false
	}

	return v.(uint32), true
}

func IsAuthed(ctx context.Context) bool {
	id, ok := GetUserID(ctx)

	return ok && id > 0
}
