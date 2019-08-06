package processing

import "testing"

func Test_commonTest(t *testing.T) {
	q := taskQ{make([]string, 0)}

	if !q.empty() {
		t.Fatalf("should be empty")
	}

	q.push("1")
	q.push("2")
	q.push("3")

	if q.empty() {
		t.Fatalf("should not be empty")
	}

	el,_ := q.pop()

	if el != "3"{
		t.Fatalf("should be 3")
	}

	el,_ = q.pop()

	if el != "2" {
		t.Fatalf("should be 2")
	}

	el,_ = q.pop()

	if el != "1" {
		t.Fatalf("should be 1")
	}

	if !q.empty() {
		t.Fatalf("should  be empty")
	}
}

func Test_sandbox(t *testing.T){

	e := -1
	arr := toArr(e)
	t.Log(arr)

}


func toArr(el int) []int {
	arr := make([]int, 0)

	temp := el

	for temp > 0{
		arr = append(arr, temp % 10)
		temp = temp / 10
	}

	return arr
}