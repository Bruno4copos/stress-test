package main

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/urfave/cli/v2"
)

type Result struct {
    StatusCode int
    Duration   time.Duration
}

func main() {
    app := &cli.App{
        Name:  "Load Tester",
        Usage: "A simple load testing CLI",
        Flags: []cli.Flag{
            &cli.StringFlag{Name: "url", Usage: "Target URL", Required: true},
            &cli.IntFlag{Name: "requests", Usage: "Total number of requests", Required: true},
            &cli.IntFlag{Name: "concurrency", Usage: "Number of concurrent requests", Required: true},
            &cli.StringFlag{Name: "method", Usage: "HTTP method to use (GET, POST, PUT, etc.)", Value: "GET"},
            &cli.StringSliceFlag{Name: "header", Usage: "Custom headers (key:value)"},
            &cli.StringFlag{Name: "body", Usage: "Body for POST/PUT requests"},
            &cli.StringFlag{Name: "csv", Usage: "Export results to CSV file (e.g., report.csv)"},
        },
        Action: func(c *cli.Context) error {
            headers := parseHeaders(c.StringSlice("header"))
            return runTest(
                c.String("url"),
                c.Int("requests"),
                c.Int("concurrency"),
                c.String("method"),
                headers,
                c.String("body"),
                c.String("csv"),
            )
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func parseHeaders(raw []string) map[string]string {
    headers := make(map[string]string)
    for _, h := range raw {
        parts := strings.SplitN(h, ":", 2)
        if len(parts) == 2 {
            headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
        }
    }
    return headers
}

func runTest(url string, totalRequests, concurrency int, method string, headers map[string]string, body, csvFile string) error {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var results []Result
    statusCodes := make(map[int]int)
    successCount := 0
    client := &http.Client{}

    startTime := time.Now()
    requestsPerWorker := totalRequests / concurrency
    remainder := totalRequests % concurrency

    for i := 0; i < concurrency; i++ {
        reqs := requestsPerWorker
        if i < remainder {
            reqs++
        }

        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            for j := 0; j < n; j++ {
                var reqBody io.Reader
                if body != "" {
                    reqBody = bytes.NewBufferString(body)
                }

                req, err := http.NewRequest(method, url, reqBody)
                if err != nil {
                    log.Printf("Failed to create request: %v", err)
                    continue
                }

                for k, v := range headers {
                    req.Header.Set(k, v)
                }

                reqStart := time.Now()
                resp, err := client.Do(req)
                reqDuration := time.Since(reqStart)

                if err != nil {
                    log.Printf("Request failed: %v", err)
                    continue
                }

                mu.Lock()
                statusCodes[resp.StatusCode]++
                if resp.StatusCode == 200 {
                    successCount++
                }
                results = append(results, Result{StatusCode: resp.StatusCode, Duration: reqDuration})
                mu.Unlock()
                resp.Body.Close()
            }
        }(reqs)
    }

    wg.Wait()
    duration := time.Since(startTime)

    fmt.Println("
--- Load Test Report ---")
    fmt.Printf("Total time taken: %v
", duration)
    fmt.Printf("Total requests: %d
", totalRequests)
    fmt.Printf("Successful (200) responses: %d
", successCount)
    fmt.Println("Other status codes:")
    for code, count := range statusCodes {
        if code != 200 {
            fmt.Printf("  %d: %d
", code, count)
        }
    }

    if csvFile != "" {
        err := exportCSV(results, csvFile)
        if err != nil {
            return fmt.Errorf("failed to export CSV: %v", err)
        }
        fmt.Printf("Results exported to %s
", csvFile)
    }

    return nil
}

func exportCSV(results []Result, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    writer.Write([]string{"Status Code", "Duration (ms)"})
    for _, r := range results {
        writer.Write([]string{
            strconv.Itoa(r.StatusCode),
            fmt.Sprintf("%.2f", r.Duration.Seconds()*1000),
        })
    }
    return nil
}
