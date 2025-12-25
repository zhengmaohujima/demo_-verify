package homework02

import (
	"fmt"
	"math"
	"os"
	"sync"
	"testing"
	"time"
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

func TestAddTen(t *testing.T) {
	defer recordResult(t, "AddTen") // 使用自动批改框架记录结果

	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"Add 10 to 0", 0, 10},
		{"Add 10 to 5", 5, 15},
		{"Add 10 to -3", -3, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.input
			AddTen(&n)
			if n != tt.want {
				t.Errorf("AddTen() = %v, want %v", n, tt.want)
			}
		})
	}
}

func TestDoubleSlice(t *testing.T) {
	defer recordResult(t, "DoubleSlice")

	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"Empty slice", []int{}, []int{}},
		{"Positive numbers", []int{1, 2, 3}, []int{2, 4, 6}},
		{"Negative numbers", []int{-1, -2}, []int{-2, -4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.input))
			copy(slice, tt.input)
			DoubleSlice(&slice)
			for i := range slice {
				if slice[i] != tt.want[i] {
					t.Errorf("DoubleSlice() = %v, want %v", slice, tt.want)
				}
			}
		})
	}
}
func TestPrintOddEven(t *testing.T) {
	defer recordResult(t, "PrintOddEven")
	PrintOddEven()
	// 这里只是运行输出，不做结果验证
}

func TestTaskScheduler(t *testing.T) {
	defer recordResult(t, "TaskScheduler")

	tasks := []func(){
		func() { time.Sleep(50 * time.Millisecond) },
		func() { time.Sleep(30 * time.Millisecond) },
	}

	durations := TaskScheduler(tasks)
	for i, d := range durations {
		if d <= 0 {
			t.Errorf("Task %d duration invalid: %v", i, d)
		}
	}
}

func TestShapes(t *testing.T) {
	defer recordResult(t, "Shapes")

	rect := Rectangle{Width: 3, Height: 4}
	if rect.Area() != 12 || rect.Perimeter() != 14 {
		t.Errorf("Rectangle calculation error")
	}

	circ := Circle{Radius: 1}
	if math.Abs(circ.Area()-math.Pi) > 1e-6 || math.Abs(circ.Perimeter()-2*math.Pi) > 1e-6 {
		t.Errorf("Circle calculation error")
	}
}

func TestEmployee(t *testing.T) {
	defer recordResult(t, "Employee")

	e := Employee{
		Person:     Person{Name: "Alice", Age: 30},
		EmployeeID: "E001",
	}

	e.PrintInfo() // 输出信息，不验证结果
}
