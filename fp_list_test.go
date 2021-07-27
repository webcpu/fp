package fp_test

import (
	"fmt"
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"strings"
)

var _ = Describe("list", func() {
	Context("Range(max)", func() {
		It("generates the list [1, 2, ..., max].", func() {
			xs := Range(5)
			expected := []int{1, 2, 3, 4, 5}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [1, 2, ..., max].", func() {
			xs := Range(3)
			expected := []int{1, 2, 3}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [1, 2, ..., max].", func() {
			xs := Range(0)
			expected := []int{}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [1, 2, ..., max].", func() {
			xs := Range(-3)
			expected := []int{}
			Expect(xs).To(Equal(expected))
		})
	})
	Context("Range(min, max)", func() {
		It("generates the list [min, ..., max].", func() {
			xs := Range(0, 3)
			expected := []int{0,1,2,3}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [min, ..., max].", func() {
			xs := Range(3, 3)
			expected := []int{3}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [min, ..., max].", func() {
			xs := Range(3, 0)
			expected := []int{}
			Expect(xs).To(Equal(expected))
		})
	})

	Context("Range(min, max, step)", func() {
		It("generates the list using step.", func() {
			xs := Range(0, 3, 2)
			expected := []int{0,2}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [min, ..., max].", func() {
			xs := Range(3, 3, 2)
			expected := []int{3}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [min, ..., max].", func() {
			xs := Range(3, 0, 2)
			expected := []int{}
			Expect(xs).To(Equal(expected))
		})

		It("generates the list [min, ..., max].", func() {
			xs := Range(3, 0, -1)
			expected := []int{3, 2, 1, 0}
			Expect(xs).To(Equal(expected))
		})
	})

	Context("Range panic", func(){
		It("generates the list [min, ..., max].", func() {
			Ω(func(){Range()}).Should(Panic())
		})

		It("generates the list [min, ..., max].", func() {
			Ω(func(){Range(3, 0, 0)}).Should(Panic())
		})

		It("generates the list using step.", func() {
			Ω(func(){Range(1, 3, 1, 1)}).Should(Panic())
		})
	})

	Context("Map(f, expr)", func() {
		It("applies f to each element in expr.", func() {
			add1 := func(x int) int {
				return x + 1
			}
			xs := Range(5)
			actual := Map(add1, xs)
			expected := []int{2, 3, 4, 5, 6}
			Expect(actual).To(Equal(expected))
		})

		It("applies f to each element in expr.", func() {
			add1 := func(x int) int {
				return x + 1
			}
			xs := [5]int{1,2,3,4,5}
			actual := Map(add1, xs)
			expected := []int{2, 3, 4, 5, 6}
			Expect(actual).To(Equal(expected))
		})

		It("applies f to each element in expr.", func() {
			double := func(x int) string {
				return fmt.Sprintf("%v%v", x, x)
			}
			xs := Range(5)
			actual := Map(double, xs)
			expected := []string{"11", "22", "33", "44", "55"}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Map panic", func() {
		It("applies f to each element in expr.", func() {
			Ω(func() { Map(Identity, 3) }).Should(Panic())
		})
	})

	Context("Filter(f, expr)", func() {
		It("picks out all elements of list for which crit(e) is true.", func() {
			evenQ := func(x int) bool {
				return x % 2 == 0
			}
			xs := Range(6)
			actual := Filter(evenQ, xs)
			expected := []int{2, 4, 6}
			Expect(actual).To(Equal(expected))
		})

		It("picks out all elements of list for which crit(e) is true.", func() {
			oddQ := func(x int) bool {
				return x % 2 == 1
			}
			xs := Range(6)
			actual := Filter(oddQ, xs)
			expected := []int{1,3,5}
			Expect(actual).To(Equal(expected))
		})

		It("picks out all elements of list for which crit(e) is true.", func() {
			uppercaseQ := func(x string) bool {
				return strings.ToUpper(x) == x
			}
			xs := []string{"abc", "DEF"}
			actual := Filter(uppercaseQ, xs)
			expected := []string{"DEF"}
			Expect(actual).To(Equal(expected))
		})

		It("picks out all elements of list for which crit(e) is true.", func() {
			uppercaseQ := func(x string) bool {
				return strings.ToUpper(x) == x
			}
			xs := [2]string{"abc", "DEF"}
			actual := Filter(uppercaseQ, xs)
			expected := []string{"DEF"}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Fold(f, r, expr)", func() {
		It("combine elements of expr using function f", func() {
			add := func(r int, x int) int {
				return r + x
			}
			xs := Range(5)
			actual := Fold(add, 0, xs)
			expected := 15
			Expect(actual).To(Equal(expected))
		})

		It("combine elements of expr using function f", func() {
			times := func(r int, x int) int {
				return r*x
			}
			xs := Range(5)
			actual := Fold(times, 1, xs)
			expected := 120
			Expect(actual).To(Equal(expected))
		})

		It("combine elements of expr using function f", func() {
			times := func(r string, x int) string {
				return r + strconv.Itoa(x)
			}
			xs := Range(5)
			actual := Fold(times, "0", xs)
			expected := "012345"
			Expect(actual).To(Equal(expected))
		})
	})

	Context("MapIndexed(f, expr)", func() {
		It("MapIndexed[f,expr] applies f to the elements of expr, giving the part specification of each element as a second argument to f.", func() {
			combine := func(x string, index int) string {
				return fmt.Sprintf("%v -> %v", index, x)
			}
			xs := []string{"ab", "cd", "ef"}
			actual := MapIndexed(combine, xs)
			expected := []string{"0 -> ab", "1 -> cd", "2 -> ef"}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("First(expr)", func() {
		It("gives the first element in expr.", func() {
			xs := Range(5)
			actual := First(xs)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("gives the first element in expr.", func() {
			xs := [5]int{1,2,3,4,5}
			actual := First(xs)
			expected := 1
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Last(expr)", func() {
		It("gives the last element in expr.", func() {
			xs := Range(5)
			actual := Last(xs)
			expected := 5
			Expect(actual).To(Equal(expected))
		})

		It("gives the last element in expr.", func() {
			xs := [5]int{1,2,3,4,5}
			actual := Last(xs)
			expected := 5
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Take(list, n)", func() {
		It("gives the first n elements in list.", func() {
			xs := Range(5)
			actual := Take(xs, 2)
			expected := []interface{}{1, 2}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Take(list, -n)", func() {
		It("gives the last n elements in list.", func() {
			xs := Range(5)
			actual := Take(xs, -2)
			expected := []interface{}{4,5}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Drop(list, n)", func() {
		It("gives list with its first n elements dropped.", func() {
			xs := Range(5)
			actual := Drop(xs, 2)
			expected := []interface{}{3,4,5}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Drop(list, -n)", func() {
		It("gives list with its last n elements dropped.", func() {
			xs := Range(5)
			actual := Drop(xs, -2)
			expected := []interface{}{1,2,3}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Length(expr)", func() {
		It("gives the numbers of elements in expr.", func() {
			xs := Range(5)
			actual := Length(xs)
			expected := 5
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Position(expr)", func() {
		It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
			xs := Range(0, 9)
			actual := Position(xs, 2)
			expected := [][]interface{}{{2}}
			Expect(actual).To(Equal(expected))
		})

		It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
			xs := [][]int{{0,1,2,3}}
			actual := Position(xs, []int{0,1,2,3})
			expected := [][]interface{}{{0}}
			Expect(actual).To(Equal(expected))
		})

		It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
			xs := [][]int{{0,1,2,3}}
			actual := Position(xs, []int{0,1,2,5})
			expected := [][]interface{}{}
			Expect(actual).To(Equal(expected))
		})

		It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
			xs := []string{"abc", "def"}
			actual := Position(xs, "def")
			expected := [][]interface{}{{1}}
			Expect(actual).To(Equal(expected))
		})

		It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
			xs := []string{"abc", "def"}
			actual := Position(xs, "adc")
			expected := [][]interface{}{}
			Expect(actual).To(Equal(expected))
		})

		It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
			xs := [2]string{"abc", "def"}
			actual := Position(xs, "adc")
			expected := [][]interface{}{}
			Expect(actual).To(Equal(expected))
		})

		//It("gives a list of the positions at which objects matching pattern appear in expr.", func() {
		//	xs := map[int]string{1:"abc", 2:"def", 7: "def"}
		//	var indices [][]interface{} = Position(xs, "def")
		//	var actual []interface{} = Sort(indices, func(a,b interface{}) bool {
		//		return a.([]interface{})[0].(int) < b.([]interface{})[0].(int)
		//	})
		//	var expected []interface{} = []interface{}{[]int{2},[]int{7}}
		//	Expect(fmt.Sprintf("%v", actual)).To(Equal(fmt.Sprintf("%v", expected)))
		//})
	})

	Context("Count(expr)", func() {
		It("gives the number of elements in list that match pattern.", func() {
			xs := Range(0, 9)
			actual := Count(xs, 2)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("gives the number of elements in list that match pattern.", func() {
			xs := [][]int{{0,1,2,3}}
			actual := Count(xs, []int{0,1,2,3})
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("gives the number of elements in list that match pattern.", func() {
			xs := [][]int{{0,1,2,3}}
			actual := Count(xs, []int{0,1,2,5})
			expected := 0
			Expect(actual).To(Equal(expected))
		})

		It("gives the number of elements in list that match pattern.", func() {
			xs := []string{"abc", "def"}
			actual := Count(xs, "def")
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("gives the number of elements in list that match pattern.", func() {
			xs := [2]string{"abc", "def"}
			actual := Count(xs, "adc")
			expected := 0
			Expect(actual).To(Equal(expected))
		})

		It("gives the number of elements in list that match pattern.", func() {
			xs := map[int]string{1:"abc", 2:"def", 7: "def"}
			actual := Count(xs, "def")
			expected := 2
			Expect(actual).To(Equal(expected))
		})

		It("gives the number of elements in list that match pattern.", func() {
			type Person struct {
				name string
			}
			xs := map[int]Person{1:Person{name:"abc"}, 2:Person{name:"def"}, 7:Person{name:"def"}}
			actual := Count(xs, Person{name:"def"})
			expected := 2
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Sort(expr)", func() {
		It("sorts the elements of list into canonical order.", func() {
			xs := []int{3, 1, 2, 5}
			actual := Sort(xs)
			expected := []int{1, 2, 3, 5}
			Expect(actual).To(Equal(expected))
		})

		XIt("sorts the elements of list into canonical order.", func() {
			xs := []int{3, 1, 2, 5}
			actual := Sort(xs, Greater)
			expected := []int{5, 3, 2, 1}
			Expect(actual).To(Equal(expected))
		})

		It("sorts the elements of list into canonical order.", func() {
			xs := []int{3, 1, 2, 5}
			less := func(a, b interface{}) bool { return a.(int) < b.(int) }
			actual := Sort(xs, less)
			expected := []int{1, 2, 3, 5}
			Expect(actual).To(Equal(expected))
		})

		It("sorts the elements of list into canonical order.", func() {
			xs := []string{"bc", "cd", "ab"}
			actual := Sort(xs)
			expected := []string{"ab", "bc", "cd"}
			Expect(actual).To(Equal(expected))
		})

		It("sorts the elements of list into canonical order.", func() {
			xs := []string{"bc", "cd", "ab"}
			less := func(a, b interface{}) bool { return a.(string) < b.(string) }
			actual := Sort(xs, less)
			expected := []string{"ab", "bc", "cd"}
			Expect(actual).To(Equal(expected))
		})
	})

	XContext("Sort(expr)", func() {
		type Person struct {
			name string
			age int
		}
		xs := []Person{
			Person{name:"b1", age: 10},
			Person{name:"a1", age: 20},
			Person{name:"c1", age: 15},
		}

		It("sorts the elements of list into canonical order.", func() {
			Ω(func() { Sort(xs) }).Should(Panic())
		})

		It("sorts the elements of list into canonical order.", func() {
			less := func(a , b interface{}) bool { return a.(Person).name < b.(Person).name }
			actual := Sort(xs, less)
			expected := []interface{}{
				Person{name:"a1", age: 20},
				Person{name:"b1", age: 10},
				Person{name:"c1", age: 15},
			}
			Expect(actual).To(Equal(expected))
		})

		It("sorts the elements of list into canonical order.", func() {
			less := func(a , b interface{}) bool { return a.(Person).age < b.(Person).age}
			actual := Sort(xs, less)
			expected := []interface{}{
				Person{name:"b1", age: 10},
				Person{name:"c1", age: 15},
				Person{name:"a1", age: 20},
			}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Reverse(expr)", func() {
		It("array", func() {
			xs := [10]int{0,1,2,3,4,5,6,7,8,9}
			actual := Reverse(xs)
			expected := []interface{}{9,8,7,6,5,4,3,2,1,0}
			Expect(actual).To(Equal(expected))
		})

		It("slice", func() {
			xs := Range(0, 9)
			actual := Reverse(xs)
			expected := []interface{}{9,8,7,6,5,4,3,2,1,0}
			Expect(actual).To(Equal(expected))
		})

		It("slice", func() {
			xs := []string{"def", "abc", "ghi"}
			actual := Reverse(xs)
			expected := []interface{}{"ghi", "abc", "def"}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Max(expr)", func() {
		It("array", func() {
			xs := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
			actual := Max(xs)
			expected := 9
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			xs := [1]int{9}
			actual := Max(xs)
			expected := 9
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			xs := [3]string{"abc", "def", "ghi"}
			actual := Max(xs)
			expected := "ghi"
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			actual := Max(1,2,3)
			expected := 3
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			actual := Max("abc", "def")
			expected := "def"
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Min(expr)", func() {
		It("array", func() {
			xs := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
			actual := Min(xs)
			expected := 0
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			xs := [1]int{1}
			actual := Min(xs)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			xs := [3]string{"abc", "def", "ghi"}
			actual := Min(xs)
			expected := "abc"
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			actual := Min(1,2,3)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			actual := Min(1)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("array", func() {
			actual := Min("abc", "def")
			expected := "abc"
			Expect(actual).To(Equal(expected))
		})
	})
})