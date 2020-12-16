This is a really bare bones version that provides only the B2 buckets sizes and file counts.
The B2 API doe snot seem to have an efficient way of calculating the sizes of buckets that doesn't involve iterating through the whole set of objects.
Consequently it takes over ten minutes for my B2 account to be counted.
I hope it doesn't run foul of any rate limits.
As long at it doesn't the update time is not a great problem as long as the update period is >> update time.


docker-compose.yaml:

```yaml
version: "3.7"

services:
  backblaze-b2-metrics-exporter:
    build:
      context: .
    image: b2exporter:local
    environment:
      B2_ACCOUNT_ID: ${B2_ACCOUNT_ID_READ}
      B2_ACCOUNT_KEY: ${B2_ACCOUNT_KEY_READ}
      TZ: Europe/Berlin
    command: -period 6h00m
    expose:
      - 8080
```

