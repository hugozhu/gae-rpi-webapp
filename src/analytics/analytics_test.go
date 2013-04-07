package analytics

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestPV(t *testing.T) {
	visitors := []string{"a", "b", "c", "d", "e", "f"}

	analytics := NewAnalytics(2, 10) //count pv in 20 seconds
	count := 1
	go func() {
		for {
			// a, slots_a := uv.Sum()
			a, slots_a, b, slots_b := analytics.Sum()
			log.Println(count, "UV", a, slots_a)
			log.Println(count, "PV", b, slots_b)
			count++
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 2; i++ {
			now := time.Now()
			zcookie := visitors[rand.Intn(len(visitors))]
			analytics.AddOne(zcookie, now)
			time.Sleep(1 * time.Second)
		}
		time.Sleep(15 * time.Second)
		for i := 0; i < 10; i++ {
			now := time.Now()
			zcookie := visitors[rand.Intn(len(visitors))]
			analytics.AddOne(zcookie, now)
			time.Sleep(20 * time.Second)
		}
		time.Sleep(41 * time.Second)
		for i := 0; i < 15; i++ {
			now := time.Now()
			zcookie := visitors[rand.Intn(len(visitors))]
			analytics.AddOne(zcookie, now)
			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(2 * time.Minute)
}
