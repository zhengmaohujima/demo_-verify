package homework01

import "sort"

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	result := 0
	for _, num := range nums {
		result ^= num // 异或抵消成对数字
	}
	return result
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	// 负数不是回文数，末尾为0的非0数字也不是回文数
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	reversed := 0
	for x > reversed {
		reversed = reversed*10 + x%10
		x /= 10
	}

	// 当数字长度为奇数时，reversed/10 去掉中间位
	return x == reversed || x == reversed/10
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	stack := []rune{}

	// 定义匹配关系
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, ch := range s {
		switch ch {
		case '(', '{', '[':
			stack = append(stack, ch) // 压栈
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
				return false
			}
			stack = stack[:len(stack)-1] // 弹栈
		}
	}

	return len(stack) == 0
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	for i := 0; i < len(strs[0]); i++ {
		char := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || strs[j][i] != char {
				return strs[0][:i]
			}
		}
	}

	return strs[0]
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	n := len(digits)

	// 从末尾开始加 1
	for i := n - 1; i >= 0; i-- {
		digits[i]++
		if digits[i] < 10 {
			// 不需要进位，直接返回
			return digits
		}
		// 需要进位
		digits[i] = 0
	}

	// 如果走到这里，说明最高位也进位了，比如 999 -> 1000
	newDigits := make([]int, n+1)
	newDigits[0] = 1
	return newDigits
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	i := 0 // i 指向最后一个不重复元素的位置
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j] // 把新的不重复元素放到 i+1 的位置
		}
	}

	return i + 1 // 长度是 i+1
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。

func Merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return nil
	}

	// 按起点排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		current := intervals[i]

		if current[0] <= last[1] { // 有重叠
			if current[1] > last[1] {
				last[1] = current[1] // 合并区间
			}
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	// key: 数字 value: 索引
	indexMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if j, ok := indexMap[complement]; ok {
			// 找到目标组合，返回两个索引
			return []int{j, i}
		}
		// 保存当前数字和它的索引
		indexMap[num] = i
	}

	// 如果没有符合条件的组合，返回空数组
	return nil
}
