package test

import (
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("list", func() {
	Context("Sort(expr)", func() {
		It("sorts the elements of list into canonical order.", func() {
			xs := []int{3, 1, 2, 5}
			actual := Sort(xs)
			expected := []int{1, 2, 3, 5}
			Expect(actual).To(Equal(expected))
		})

		It("sorts the elements of list into canonical order.", func() {
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

	Context("Sort(expr)", func() {
		type Person struct {
			name string
			age  int
		}
		xs := []Person{
			Person{name: "b1", age: 10},
			Person{name: "a1", age: 20},
			Person{name: "c1", age: 15},
		}

		It("sorts the elements of list into canonical order.", func() {
			Ω(func() { Sort(xs) }).Should(Panic())
		})

		It("sorts the elements of list into canonical order.", func() {
			less := func(a, b interface{}) bool { return a.(Person).name < b.(Person).name }
			actual := Sort(xs, less)
			expected := []Person{
				Person{name: "a1", age: 20},
				Person{name: "b1", age: 10},
				Person{name: "c1", age: 15},
			}
			Expect(actual).To(Equal(expected))
		})

		It("sorts the elements of list into canonical order.", func() {
			less := func(a, b interface{}) bool { return a.(Person).age < b.(Person).age }
			actual := Sort(xs, less)
			expected := []Person{
				Person{name: "b1", age: 10},
				Person{name: "c1", age: 15},
				Person{name: "a1", age: 20},
			}
			Expect(actual).To(Equal(expected))
		})
	})
})