package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
)

func RunInParallel(ctx context.Context, n int, parallel int, f func(idx int) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultChan := make(chan error, n)
	indexChan := make(chan int, n)
	for i := 0; i < n; i++ {
		indexChan <- i
	}
	close(indexChan)

	wg := sync.WaitGroup{}
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(idx int) {
			log.Printf("[%d]: started", idx)
			defer wg.Done()
			defer func() {
				log.Printf("[%d]: exiting", idx)
			}()

			for {
				select {
				case jobIdx, ok := <-indexChan:
					if !ok {
						log.Printf("[%d]: closed", idx)
						return
					}

					if err := f(jobIdx); err != nil {
						resultChan <- err
						log.Printf("[%d]: got error: %v", idx, err)

						cancel()
						return
					}
				case <-ctx.Done():
					log.Printf("[%d]: done", idx)

					return
				}
			}
		}(i)
	}

	log.Printf("waiting")
	wg.Wait()
	close(resultChan)

	var errors []string
	log.Printf("checking errors")
	for err := range resultChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}
