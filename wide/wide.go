// This package implements a logging and tracing system based on "wide events".
// Instead of having different logging/tracing/metrics systems there is only
// an "event" which contains a very large amount of context data.
// Different exporters can be used, the default is slog.
//
// In the tradition of this library, this module expands on the std library's
// expvar package, exporting values to it and using some of its types.

package wide

import (
	"cmp"
	"context"
	"expvar"

	"github.com/rushsteve1/fp"
	"golang.org/x/sync/errgroup"
)

var base = Event{
	vars: expvar.NewMap("base"),
}

func SlogExporter(ctx context.Context, m *expvar.Map) error {
	// slog.LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

var exporter func(context.Context, *expvar.Map) error = SlogExporter

func Publish(name string, v expvar.Var) {
	expvar.Publish(name, v)
	base.vars.Set(name, v)
}

const MAIN_WIDE_EVENT_KEY = "MAIN_WIDE_EVENT"

func MainEvent(ctx context.Context) Event {
	return ctx.Value(MAIN_WIDE_EVENT_KEY).(Event)
}

func SetMainEvent(ctx context.Context, evt Event) context.Context {
	return context.WithValue(ctx, MAIN_WIDE_EVENT_KEY, evt)
}

func GetMainEvent(ctx context.Context) Event {
	return context.GetValue(ctx, MAIN_WIDE_EVENT_KEY)
}

type Event struct {
	ctx  context.Context
	vars *expvar.Map
}

func NewEvent(ctx context.Context, name string) Event {
	return Event{
		ctx:  ctx,
		vars: expvar.NewMap(name),
	}
}

func WithContext(evt Event, ctx context.Context) Event {
	evt.ctx = ctx
	return evt
}

func EventGroup(ctx context.Context, name string, fs ...func(context.Context) error) error {
	evt := NewEvent(name)
	defer evt.Close()
	ctx = SetMainEvent(ctx, evt)
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(len(fs))

	for _, f := range fs {
		eg.Go(func() error {
			return f(ctx)
		})
	}

	err := eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (e Event) Close() error {
	return exporter(e.ctx, e.vars)
}
