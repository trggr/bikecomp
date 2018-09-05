package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    lat := 41.813318 + r.Float32()/1000.0
    long := -72.7435836 + r.Float32()/1000.0
    alt := 0.0 + r.Float32()
    accuracy := 19.839000701 + r.Float32()
    bearing := 67.5 + r.Float32()

    fmt.Printf("{\n");
    fmt.Printf("  \"latitude\": %f,\n", lat);
    fmt.Printf("  \"longitude\": %f,\n", long);
    fmt.Printf("  \"altitude\": %f,\n", alt);
    fmt.Printf("  \"accuracy\": %f,\n", accuracy);
    fmt.Printf("  \"bearing\": %f\n", bearing);
    fmt.Printf("}\n");
}
