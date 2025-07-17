package app

import (
	"context"
	"time"

	"k8s.io/klog/v2"
)

func Run(ctx context.Context, identity string) {
	klog.Infof("%s: Became leader, starting main app logic", identity)

	for {
		select {
		case <-ctx.Done():
			klog.Infof("%s: Leadership cancelled", identity)
			return
		case <-time.After(1 * time.Second): // Was 5s before
			klog.Infof("%s: Acting as leader", identity)
			// Here you can invoke Spring Boot call or other logic
		}
	}
}
