package statsdclient

import (
	"strconv"
	"time"

	cactus "github.com/cactus/go-statsd-client/v5/statsd"
)

type CactusStatsdClientAdapter struct {
	cactusClient cactus.Statter
}

func (a *CactusStatsdClientAdapter) Open() {

}

func (a *CactusStatsdClientAdapter) Close() {
	a.cactusClient.Close()
}

func (a *CactusStatsdClientAdapter) Count(key string, value int, sampleRate float32) {
	a.cactusClient.Inc(key, int64(value), sampleRate)
}

func (a *CactusStatsdClientAdapter) Timing(key string, time int64, sampleRate float32) {
	a.cactusClient.Timing(key, time, sampleRate)
}

func (a *CactusStatsdClientAdapter) Gauge(key string, value int) {
	a.cactusClient.Gauge(key, int64(value), 1)
}

func (a *CactusStatsdClientAdapter) GaugeShift(key string, value int) {
	a.cactusClient.GaugeDelta(key, int64(value), 1)
}

func (a *CactusStatsdClientAdapter) Set(key string, value int) {
	a.cactusClient.SetInt(key, int64(value), 1)
}

func NewCactusClient(
	statsdHost string,
	statsdPort int,
) StatsdClientInterface {
	cactusClient, _ := cactus.NewBufferedClient(
		statsdHost+":"+strconv.Itoa(statsdPort),
		"",
		100*time.Microsecond,
		1432,
	)

	return &CactusStatsdClientAdapter{
		cactusClient: cactusClient,
	}
}
