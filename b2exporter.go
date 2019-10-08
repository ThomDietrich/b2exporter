package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/kurin/blazer/b2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

var addr = flag.String("listen", ":8080", "The address to listen on for HTTP requests.")
var period = flag.Duration("period", 60*time.Minute, "The update period.")
var bucketSizeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "b2_bucket_size_bytes",
	Help: "B2 bucket size.",
}, []string{"name"})
var bucketCountGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "b2_bucket_count",
	Help: "B2 bucket object count.",
}, []string{"name"})

func init() {
	prometheus.MustRegister(bucketSizeGauge)
	prometheus.MustRegister(bucketCountGauge)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func main() {
	flag.Parse()

	log.Println("b2exporter")
	log.Println(period)

	godotenv.Load()

	go func() {
		for {
			update()
			time.Sleep(*period)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func update() {
	b2id := os.Getenv("B2_ACCOUNT_ID")
	b2key := os.Getenv("B2_ACCOUNT_KEY")
	ctx := context.Background()
	c, err := b2.NewClient(ctx, b2id, b2key)
	if err != nil {
		log.Fatal(err)
		return
	}
	buckets, err := c.ListBuckets(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, bucket := range buckets {
		name := bucket.Name()
		log.Println(name)
		iterator := bucket.List(ctx, b2.ListHidden())
		var size int64 = 0
		var count int64 = 0
		for iterator.Next() {
			attrs, err := iterator.Object().Attrs(ctx)
			if err != nil {
				log.Fatal(err)
				return
			}
			count = count + 1
			size = size + attrs.Size
		}
		bucketSizeGauge.WithLabelValues(name).Set(float64(size))
		bucketCountGauge.WithLabelValues(name).Set(float64(count))
		log.Println(count)
		log.Println(size)
	}
}
