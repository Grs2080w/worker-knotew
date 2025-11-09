package memory

import (
	"runtime"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/prometheus"
)

func MonitorResources() {
    var m runtime.MemStats
    previousGC := uint32(0)

    for {
        runtime.ReadMemStats(&m)

        prometheus.MemoryUsageBytes.Set(float64(m.Alloc))
        prometheus.Goroutines.Set(float64(runtime.NumGoroutine()))

        if m.NumGC > previousGC {
			prometheus.GCCount.Add(float64(m.NumGC - previousGC))
            previousGC = m.NumGC
        }

        time.Sleep(1 * time.Second)
    }
}
