package fp_test

import (
	"fmt"
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			add1 := func(x interface{}) interface{} {
				return x.(int) + 1
			}
			xs := Range(5)
			actual := Map(add1, xs)
			expected := []interface{}{2, 3, 4, 5, 6}
			Expect(actual).To(Equal(expected))
		})

		It("applies f to each element in expr.", func() {
			double := func(x interface{}) interface{} {
				return fmt.Sprintf("%v%v", x, x)
			}
			xs := Range(5)
			actual := Map(double, xs)
			expected := []interface{}{"11", "22", "33", "44", "55"}
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
			evenQ := func(x interface{}) bool {
				return x.(int) % 2 == 0
			}
			xs := Range(6)
			actual := Filter(evenQ, xs)
			expected := []interface{}{2, 4, 6}
			Expect(actual).To(Equal(expected))
		})

		It("picks out all elements of list for which crit(e) is true.", func() {
			oddQ := func(x interface{}) bool {
				return x.(int) % 2 == 1
			}
			xs := Range(6)
			actual := Filter(oddQ, xs)
			expected := []interface{}{1,3,5}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Fold(f, r, expr)", func() {
		It("combine elements of expr using function f", func() {
			add := func(r interface{}, x interface{}) interface{} {
				return r.(int) + x.(int)
			}
			xs := Range(5)
			actual := Fold(add, 0, xs)
			expected := 15
			Expect(actual).To(Equal(expected))
		})

		It("combine elements of expr using function f", func() {
			times := func(r interface{}, x interface{}) interface{} {
				return r.(int) * x.(int)
			}
			xs := Range(5)
			actual := Fold(times, 1, xs)
			expected := 120
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
	})

	Context("Last(expr)", func() {
		It("gives the last element in expr.", func() {
			xs := Range(5)
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
})
