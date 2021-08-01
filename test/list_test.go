package test

import (
	"fmt"
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"strings"
	"time"
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

	Context("ParallelMap(f, expr)", func() {
		It("applies f to each element in expr.", func() {
			add1 := func(x int) int {
				time.Sleep(10 * time.Millisecond)
				return x + 1
			}
			xs := Range(5)
			t, actual := Timing(func() interface{} {return ParallelMap(add1, xs)})
			expected := Range(2,6)
			Expect(actual).To(Equal(expected))
			fmt.Printf("t = %v", t)
		})

		It("applies f to each element in expr.", func() {
			add1 := func(x int) int {
				return x + 1
			}
			xs := [5]int{1, 2, 3, 4, 5}
			t, actual := Timing(func() interface{} {return ParallelMap(add1, xs)})
			expected := []int{2, 3, 4, 5, 6}
			Expect(actual).To(Equal(expected))
			fmt.Printf("t = %v", t)
		})
	})

	Context("MapThread(f, list...)", func() {
		It("applies f to each element in expr.", func() {
			add := func(x, y int) int {
				return x + y
			}
			xs := Range(5)
			ys := []int{5,4,3,2,1}
			actual := MapThread(add, xs, ys)
			expected := []int{6, 6, 6, 6, 6}
			Expect(actual).To(Equal(expected))
		})

		It("applies f to each element in expr.", func() {
			add := func(x, y int) int {
				return x + y
			}
			xs := Range(5)
			ys := []int{5,4,3,2,1,2}
			actual := MapThread(add, xs, ys)
			expected := []int{6, 6, 6, 6, 6}
			Expect(actual).To(Equal(expected))
		})

		It("applies f to each element in expr.", func() {
			add := func(x int, s string) string {
				return strconv.Itoa(x) + s
			}
			xs := Range(5)
			ys := []string{"a","b","c","d","e"}
			actual := MapThread(add, xs, ys)
			expected := []string{"1a","2b","3c","4d","5e"}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Do(f, expr)", func() {
		It("applies f to each element in expr.", func() {
			actual := make([]int, 5)
			update := func(x int) {
				actual[x-1] = x
			}

			xs := Range(5)
			expected := xs
			Do(update, xs)
			Expect(actual).To(Equal(expected))
		})

		It("applies f to each element in expr.", func() {
			actual := make([]int, 5)
			update := func(x int) {
				actual[x-1] = x*x
			}
			xs := Range(5)
			expected := Map(func(x int) int {return x*x}, xs)
			Do(update, xs)
			Expect(actual).To(Equal(expected))
		})
	})

	Context("ParallelDo(f, expr)", func() {
		It("applies f to each element in expr.", func() {
			idle := func(x int) {
				time.Sleep(10 * time.Millisecond)
			}
			xs := Range(5)
			t, actual := Timing(func() {ParallelDo(idle, xs)})
			Expect(actual).To(BeNil())
			fmt.Printf("t = %v", t)
		})

		It("applies f to each element in expr.", func() {
			doNothing := func(x int) {
				_ = x + 1
			}
			xs := [5]int{1, 2, 3, 4, 5}
			t, actual := Timing(func() {ParallelDo(doNothing, xs)})
			Expect(actual).To(BeNil())
			fmt.Printf("t = %v", t)
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
			expected := []int{1, 2}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Take(list, -n)", func() {
		It("gives the last n elements in list.", func() {
			xs := Range(5)
			actual := Take(xs, -2)
			expected := []int{4,5}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Drop(list, n)", func() {
		It("gives list with its first n elements dropped.", func() {
			xs := Range(5)
			actual := Drop(xs, 2)
			expected := []int{3,4,5}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Drop(list, -n)", func() {
		It("gives list with its last n elements dropped.", func() {
			xs := Range(5)
			actual := Drop(xs, -2)
			expected := []int{1,2,3}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Most(list)", func() {
		It("[]int", func() {
			xs := Range(5)
			actual := Most(xs)
			expected := []int{1, 2, 3, 4}
			Expect(actual).To(Equal(expected))
		})

		It("[]string", func() {
			xs := []string{"abc", "def"}
			actual := Most(xs)
			expected := []string{"abc"}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Rest(list)", func() {
		It("[]int", func() {
			xs := Range(5)
			actual := Rest(xs)
			expected := Range(2, 5)
			Expect(actual).To(Equal(expected))
		})

		It("[]string", func() {
			xs := []string{"abc", "def"}
			actual := Rest(xs)
			expected := []string{"def"}
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
			actual := Count(xs, Person{name: "def"})
			expected := 2
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

	Context("MemberQ(list, x)", func() {
		It("returns true if an element of list matches form, and false otherwise.", func() {
			xs := Range(10)
			actual := MemberQ(xs, 1)
			expected := true
			Expect(actual).To(Equal(expected))
		})

		It("returns true if an element of list matches form, and false otherwise.", func() {
			xs := Range(10)
			actual := MemberQ(xs, 11)
			expected := false
			Expect(actual).To(Equal(expected))
		})

		It("returns true if an element of list matches form, and false otherwise.", func() {
			xs := []string{"abc", "def"}
			actual := MemberQ(xs, "abc")
			expected := true
			Expect(actual).To(Equal(expected))
		})

		It("returns true if an element of list matches form, and false otherwise.", func() {
			xs := []string{"abc", "def"}
			actual := MemberQ(xs, "acd")
			expected := false
			Expect(actual).To(Equal(expected))
		})

		It("returns true if an element of list matches form, and false otherwise.", func() {
			type Person struct {
				name string
				age int
			}
			person := Person{name:"b1", age: 10}
			xs := []interface{}{1, "def", []int{1,2,3},Person{name:"b1", age: 10}}
			Expect(MemberQ(xs, 1)).To(Equal(true))
			Expect(MemberQ(xs, "def")).To(Equal(true))
			Expect(MemberQ(xs, []int{1,2,3})).To(Equal(true))
			Expect(MemberQ(xs, []int{1,2,2})).To(Equal(false))
			Expect(MemberQ(xs, person)).To(Equal(true))
			Expect(MemberQ(xs, Person{name: "b1", age: 11})).To(Equal(false))
		})
	})

	Context("KeyMemberQ(dict, key)", func() {
		It("yields true if a key in the association assoc matches key, and false otherwise.", func() {
			m := map[string]int{"a": 1, "b": 2}
			Expect(KeyMemberQ(m, "a")).To(Equal(true))
			Expect(KeyMemberQ(m, "b")).To(Equal(true))
			Expect(KeyMemberQ(m, "c")).To(Equal(false))
		})

		It("yields true if a key in the association assoc matches key, and false otherwise.", func() {
			m := map[interface{}]int{"a": 1, "b": 2, 3: 3}
			Expect(KeyMemberQ(m, "a")).To(Equal(true))
			Expect(KeyMemberQ(m, "b")).To(Equal(true))
			Expect(KeyMemberQ(m, 3)).To(Equal(true))
			Expect(KeyMemberQ(m, "c")).To(Equal(false))
		})
	})

	Context("Keys(dict)", func() {
		It("[]string", func() {
			m := map[string]int{"a": 1, "b": 2}
			Expect(Sort(Keys(m))).To(Equal([]string{"a", "b"}))
		})
	})

	Context("Values(dict)", func() {
		It("[]int", func() {
			m := map[string]int{"a": 1, "b": 2}
			Expect(Sort(Values(m))).To(Equal([]int{1, 2}))
		})
	})

	Context("Union(list1, list2...)", func() {
		It("[]int", func() {
			xs := []int{1,3,2,2}
			actual := Union(xs)
			expected := []int{1,3,2}
			Expect(actual).To(Equal(expected))
		})

		It("[]int, []int", func() {
			xs := []int{1,2,2,3}
			ys := []int{1,2,3,3,4}
			actual := Union(xs, ys)
			expected := []int{1,2,3,4}
			Expect(actual).To(Equal(expected))
		})

		It("[]int, []int, float32[]", func() {
			xs := []int{1,2,2,3}
			ys := []int{1,2,3,3,4}
			zs := []float32{1,2,3,3,4}
			Ω(func(){Union(xs,ys,zs)}).Should(Panic())
		})

		It("[]Person", func() {
			type Person struct {
				name string
				age int
			}
			xs := []Person{
				Person{name:"a1", age: 10},
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			ys := []Person{
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			actual := Union(xs)
			expected := ys
			Expect(actual).To(Equal(expected))
		})
	})

	Context("DeleteDuplicates(list, test)", func() {
		It("[]int", func() {
			xs := []int{1, 2, 2, 3}
			actual := DeleteDuplicates(xs)
			expected := []int{1, 2, 3}
			Expect(actual).To(Equal(expected))
		})

		It("[]int", func() {
			xs := []int{1, 1, 2, 2, 3}
			//lessThan3 := func(x interface{}, y interface{}) bool { return Abs(x.(int)-y.(int)) == 1}
			lessThan3 := func(x int, y int) bool { return Abs(x-y) == 1}
			actual := DeleteDuplicates(xs, lessThan3)
			expected := []int{1, 1, 3}
			Expect(actual).To(Equal(expected))
		})

		It("[]string", func() {
			xs := []string{"abc", "abc", "def", "def"}
			actual := DeleteDuplicates(xs)
			expected := []string{"abc", "def"}
			Expect(actual).To(Equal(expected))
		})

		It("[]Person", func() {
			type Person struct {
				name string
				age int
			}
			xs := []Person{
				Person{name:"a1", age: 10},
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			ys := []Person{
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			actual := DeleteDuplicates(xs)
			expected := ys
			Expect(actual).To(Equal(expected))
		})

		It("[]Person", func() {
			type Person struct {
				name string
				age int
			}
			xs := []Person{
				Person{name:"a1", age: 10},
				Person{name:"b1", age: 10},
				Person{name:"c1", age: 15},
			}

			ys := []Person{
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			SameAgeQ := func(p1, p2 Person) bool {return p1.age == p2.age}
			actual := DeleteDuplicates(xs, SameAgeQ)
			expected := ys
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Intersection(list1, list2...)", func() {
		It("[]int", func() {
			xs := []int{1, 2}
			ys := []int{2, 3}
			actual := Intersection(xs, ys)
			expected := []int{2}
			Expect(actual).To(Equal(expected))
		})

		It("[]int, []int", func() {
			xs := []int{1, 2, 2, 3}
			ys := []int{1, 2, 3, 3, 4}
			actual := Intersection(xs, ys)
			expected := []int{1, 2, 3}
			Expect(actual).To(Equal(expected))
		})

		It("[]int, []int, []int", func() {
			xs := []int{1, 2, 2, 3}
			ys := []int{1, 2, 3, 3, 4}
			zs := []int{3, 4, 5}
			actual := Intersection(xs, ys, zs)
			expected := []int{3}
			Expect(actual).To(Equal(expected))
		})

		It("[]Person, []Person", func() {
			type Person struct {
				name string
				age int
			}
			xs := []Person{
				Person{name:"a1", age: 10},
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			ys := []Person{
				Person{name:"a1", age: 10},
				Person{name:"c1", age: 15},
			}

			actual := Intersection(xs, ys)
			expected := ys
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Intersection(list1..., f)", func() {
		It("[]int", func() {
			xs := []int{1, 2, 3, 4}

			f := func(x int, y int) bool { return Abs(x-y) == 1}
			actual := Intersection(xs, f)
			expected := []int{4}
			Expect(actual).To(Equal(expected))
		})

		It("[]int, []int", func() {
			xs := []int{1, 2}
			ys := []int{2, 3}

			f := func(x int, y int) bool { return Abs(x-y) == 1}
			actual := Intersection(xs, ys, f)
			expected := []int{2}
			Expect(actual).To(Equal(expected))
		})
	})
})