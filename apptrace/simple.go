package apptrace

import "context"

func Ignore(ctx context.Context) {
	FromContext(ctx).WithIgnore(true)
}

func Log(ctx context.Context) {
	FromContext(ctx).WithLog(true)
}

func DontLog(ctx context.Context) {
	FromContext(ctx).WithLog(false)
}

func Notice(ctx context.Context, err error) {
	FromContext(ctx).WithError(err)
}

func Attributes(ctx context.Context, keysAndValues ...interface{}) {
	FromContext(ctx).WithAttributes(keysAndValues...)
}

func AttributesMap(ctx context.Context, attributes map[string]interface{}) {
	FromContext(ctx).WithAttributesMap(attributes)
}
