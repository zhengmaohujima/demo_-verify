package homework05

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var (
	failedQuestions []string
	totalQuestions  int
	mu              sync.Mutex
)

func recordResult(t *testing.T, name string) {
	mu.Lock()
	defer mu.Unlock()
	totalQuestions++
	if t.Failed() {
		failedQuestions = append(failedQuestions, name)
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Print summary
	if totalQuestions > 0 {
		fmt.Println("\n---------------------------------------------------")
		fmt.Printf("Total Questions: %d\n", totalQuestions)
		fmt.Printf("Passed: %d\n", totalQuestions-len(failedQuestions))
		fmt.Printf("Failed: %d\n", len(failedQuestions))

		score := float64(totalQuestions-len(failedQuestions)) / float64(totalQuestions) * 100
		fmt.Printf("Score: %.2f%%\n", score)

		if len(failedQuestions) > 0 {
			fmt.Println("Failed Questions:")
			for _, q := range failedQuestions {
				fmt.Printf("- %s\n", q)
			}
		}
		fmt.Println("---------------------------------------------------")
	}

	os.Exit(code)
}
