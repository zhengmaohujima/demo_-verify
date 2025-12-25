package homework02

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// 指针1. 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。
func AddTen(n *int) {
	if n != nil {
		*n += 10
	}
}

// 指针2. DoubleSlice 接收一个整数切片的指针，将切片中的每个元素乘以2
func DoubleSlice(nums *[]int) {
	if nums == nil {
		return
	}
	for i := range *nums {
		(*nums)[i] *= 2
	}
}

// Goroutine 题目一：奇数和偶数打印
func PrintOddEven() {
	done := make(chan bool)

	// 打印奇数
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Println("Odd:", i)
		}
		done <- true
	}()

	// 打印偶数
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Println("Even:", i)
		}
		done <- true
	}()

	// 等待两个协程完成
	<-done
	<-done
}

// Goroutine 题目二：任务调度器
func TaskScheduler(tasks []func()) []time.Duration {
	var wg sync.WaitGroup
	durations := make([]time.Duration, len(tasks))

	for i, task := range tasks {
		wg.Add(1)
		go func(i int, task func()) {
			defer wg.Done()
			start := time.Now()
			task()
			durations[i] = time.Since(start)
		}(i, task)
	}

	wg.Wait()
	return durations
}

// 面向对象题目一：Shape 接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

//面向对象题目二：组合 Person 和 Employee

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("Name: %s, Age: %d, EmployeeID: %s\n", e.Name, e.Age, e.EmployeeID)
}
