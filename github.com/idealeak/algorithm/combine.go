package algorithm

import (
	"container/list"
	"fmt"
)

var LogEnable = false

func log(a ...interface{}) {
	if LogEnable {
		fmt.Println(a...)
	}
}

func CombNumber(n, m uint64) uint64 {
	a := uint64(1)
	b := uint64(1)
	for i := m; i > 0; i-- {
		a *= n
		n -= 1
		b *= i
	}
	return a / b
}

func CombinerSelect(result *list.List, data, workspace []byte, m, n int) {
	if len(workspace) == n {
		result.PushBack(workspace)
		//log(workspace)
		return
	}

	for i := 0; i < len(data); i++ {
		copyWorkspace := make([]byte, len(workspace))
		copy(copyWorkspace, workspace)
		copyWorkspace = append(copyWorkspace, data[i])
		copyData := data[i+1:]
		if len(copyData)+len(copyWorkspace) < n {
			continue
		}
		CombinerSelect(result, copyData, copyWorkspace, m, n)
	}
}

type CombineTask struct {
	Data      []byte
	Workspace []byte
	N         int
	M         int
	Result    chan []byte
}

func ConcurrentCombinerSelect(result chan []byte, taskQue chan *CombineTask, data, workspace []byte, m, n int) {
	for i := 0; i < len(data); i++ {
		copyWorkspace := make([]byte, len(workspace))
		copy(copyWorkspace, workspace)
		copyWorkspace = append(copyWorkspace, data[i])
		if len(copyWorkspace) == n {
			//log(copyWorkspace)
			result <- copyWorkspace
			continue
		}
		copyData := data[i+1:]
		if len(copyData)+len(copyWorkspace) < n {
			continue
		}
		task := &CombineTask{
			Data:      copyData,
			Workspace: copyWorkspace,
			N:         n,
			M:         m,
			Result:    result,
		}
		taskQue <- task
	}
}

func CombinerSelectNoRecursion(atable []int32, n int) [][]int32 {
	cnt := len(atable)
	if n > len(atable) {
		return nil
	}

	meta := make([]int32, cnt)
	//init meta data
	for i := 0; i < cnt; i++ {
		if i < n {
			meta[i] = 1
		} else {
			meta[i] = 0
		}
	}

	var result [][]int32

	//记录一次组合
	var tmp []int32
	for i := 0; i < cnt; i++ {
		if meta[i] == 1 {
			tmp = append(tmp, atable[i])
		}
	}
	result = append(result, tmp)

	for {
		//前面连续的0
		zero_count := 0
		for i := 0; i < cnt-n; i++ {
			if meta[i] == 0 {
				zero_count = zero_count + 1
			} else {
				break
			}
		}
		// 前m-n位都是0，说明处理结束
		if zero_count == cnt-n {
			break
		}

		var idx int
		for j := 0; j < cnt-1; j++ {
			// 10 交换为 01
			if meta[j] == 1 && meta[j+1] == 0 {
				meta[j], meta[j+1] = meta[j+1], meta[j]
				idx = j
				break
			}
		}
		// 将idx左边所有的1移到最左边
		var k = idx
		var count = 0
		for count <= k {
			for i := k; i >= 1; i-- {
				if meta[i] == 1 {
					meta[i], meta[i-1] = meta[i-1], meta[i]
				}
			}
			count = count + 1
		}

		// 记录一次组合
		var tmp []int32
		for i := 0; i < cnt; i++ {
			if meta[i] == 1 {
				tmp = append(tmp, atable[i])
			}
		}
		result = append(result, tmp)
	}

	return result
}

func CombinerSelectUseRecursion(atable []int32, n int) [][]int32 {
	var r [][]int32

	var f func(t, a []int32, num int)
	f = func(t, a []int32, num int) {
		if num == 0 {
			r = append(r, t)
			return
		}

		l := len(a)
		for i := 0; i <= l-num; i++ {
			tt := make([]int32, len(t), n)
			copy(tt, t)
			tt = append(tt, a[i])
			f(tt, a[i+1:], num-1)
		}
	}

	f([]int32{}, atable, n)
	return r
}

//全排列
func FullPermutation(atable []int32) [][]int32 {
	var r [][]int32

	var f func(a []int32, idx int)
	f = func(a []int32, idx int) {
		if idx >= len(a) {
			dst := make([]int32, len(atable))
			copy(dst, a)
			r = append(r, dst)
			return
		}

		for i := idx; i < len(a); i++ {
			a[i], a[idx] = a[idx], a[i]
			f(a, idx+1)
			a[i], a[idx] = a[idx], a[i]
		}
	}

	f(atable, 0)
	return r
}
