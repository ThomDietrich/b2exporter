package main

import (
	"fmt"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/kurin/blazer/b2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/inhies/go-bytesize"
	"log"
	"net/http"
	"os"
	"time"
)

var addr = flag.String("listen", ":8080", "The address to listen on for HTTP requests.")
var period = flag.Duration("period", 60*time.Minute, "The update period.")
var bucketSizeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "b2_bucket_size_bytes",
	Help: "Size (in bytes) of the Backlaze B2 bucket.",
}, []string{"name"})
var bucketCountGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "b2_bucket_count",
	Help: "Count of objects in the Backblaze B2 bucket.",
}, []string{"name"})

func main() {
	log.Println("Welcome to b2exporter v1.0")
	flag.Parse()
	godotenv.Load()

	b2id := os.Getenv("B2_ACCOUNT_ID")
	b2key := os.Getenv("B2_ACCOUNT_KEY")
	if len(b2id) == 0 || len(b2key) == 0 {
		log.Println("Error. Environment variables for B2_ACCOUNT_ID and B2_ACCOUNT_KEY not provided. Please check your configuration and restart. Exiting ...")
		os.Exit(1)
	} else {
		log.Println("Environment variables B2_ACCOUNT_ID and B2_ACCOUNT_KEY available.")
	}
	log.Println("Update period:", period)

	r := prometheus.NewRegistry()
	r.MustRegister(bucketSizeGauge)
	r.MustRegister(bucketCountGauge)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	go func() {
		log.Println("Entering update loop and listening for requests ...")
		for {
			update(b2id, b2key)
		        log.Println("Sleeping for", period)
			time.Sleep(*period)
		}
	}()

	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func update(b2id string, b2key string) {
	ctx := context.Background()
	c, err := b2.NewClient(ctx, b2id, b2key)
	if err != nil {
		log.Println(err)
		log.Println("The b2 client ran into an error while connecting. Please check your configuration. Trying again in the next loop cycle ...")
		bucketSizeGauge.Reset()
		bucketCountGauge.Reset()
		return
	}
	buckets, err := c.ListBuckets(ctx)
	if err != nil {
		log.Println(err)
		log.Println("The b2 client ran into an error while listing your buckets. Please check your configuration. Trying again in the next loop cycle ...")
		bucketSizeGauge.Reset()
		bucketCountGauge.Reset()
		return
	}
	for _, bucket := range buckets {
		name := bucket.Name()
		log.Println("Scanning content of bucket", name)
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
		log.Println(fmt.Sprintf("Bucket '%s' with %d objects occupies %s", name, count, bytesize.New(float64(size))))
	}
}
