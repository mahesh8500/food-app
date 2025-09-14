package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ApiKeyMiddleware validates the api_key header before calling the next handler
func ApiKeyMiddleware(expectedKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("api_key")
		fmt.Printf("DEBUG api_key header: expected='%s', got='%s'\n", expectedKey, got)
		if r.Header.Get("api_key") != expectedKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoadCoupons loads promo codes.
func LoadCoupons(paths []string) (map[string]int, error) {
	var wg sync.WaitGroup
	results := make(chan map[string]bool, len(paths))
	errs := make(chan error, len(paths))

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			file, err := os.Open(p)
			if err != nil {
				errs <- err
				return
			}
			defer file.Close()

			local := make(map[string]bool)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				code := strings.TrimSpace(scanner.Text())
				if len(code) >= 8 && len(code) <= 10 {
					local[code] = true
				}
			}
			if err := scanner.Err(); err != nil {
				errs <- err
				return
			}
			results <- local
		}(path)
	}

	wg.Wait()
	close(results)
	close(errs)

	counts := make(map[string]int)
	for res := range results {
		for code := range res {
			counts[code]++
		}
	}
	if len(errs) > 0 {
		return nil, <-errs
	}
	return counts, nil
}
