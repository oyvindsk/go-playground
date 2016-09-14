package foo

import (
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
)

func doStuff(ctx context.Context, n1, n2 int) (int, *fooError) {
	log.Infof(ctx, "doing stuff, data: %d %d", n1, n2)

	val := ctx.Value("KEY").(int) // FIXME

	log.Infof(ctx, "doing stuff, context: %d", val)

	// oh snap
	err := fmt.Errorf("OH SNAP")
	return 0, &fooError{err, "Could not doStuff()", 500}

	// return n1 + n2, nil

}
