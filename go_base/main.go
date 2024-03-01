package main

import (
	"fmt"
	"strconv"
)

type Person struct {
	username string
	password string
	age      int
}

func main() {
	var a string = "Runoob"
	// slice切片类型
	sl := []string{"d", "c"}
	var slm []string
	// 使用append添加数据需要替换原来的切换才可以
	sl = append(sl, "12")
	fmt.Println(sl)
	// 获取指定长度的切片内容：可以用来删除第一项数据或者获取数据之类操作
	sl = sl[1:3]
	fmt.Println(sl)
	sl = slm
	fmt.Println(len(sl))
	if sl == nil {
		fmt.Println("空切片")
	}
	// map类型
	ccn := map[string]string{
		"cn": "admin12",
	}
	// 获取map内的键值
	v1 := ccn["cn"]

	fmt.Println("v1:" + v1)

	// 获取map的len
	ccnLen := len(ccn)
	// 这里的ccnLen为int 需要转换类型，go没有类型的强制转换
	fmt.Println("ccnLen:" + strconv.Itoa(ccnLen))
	ccn["new"] = "admin123"
	// 遍历 Map

	// Go 语言中 range 关键字用于 for 循环中迭代数组(array)、切片(slice)、通道(channel)或集合(map)的元素。在数组和切片中它返回元素的索引和索引对应的值，在集合中返回 key-value 对。
	for k, v := range ccn {
		fmt.Printf("key=%s, value=%d\n", k, v)
	}
	// 删除键
	delete(ccn, "new")
	// 获取new的值
	new, ok := ccn["new"]
	if ok {
		fmt.Println("new的值", new)
	} else {
		fmt.Println("new已经删除")
	}
	var person Person
	person.username = "admin123"
	person.age = 12
	person.password = "admin123"
	cv := 123
	fmt.Println(ccn)
	fmt.Println(cv)
	fmt.Println(a)
	fmt.Println(person)

	var b, c int = 1, 2
	fmt.Println(b, c)
}
