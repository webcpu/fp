package test

import (
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = Describe("function", func() {
	Context("Apply(f, expr)", func() {
		It("1 parameter", func() {
			xs := Apply(Range, []interface{}{5})
			expected := []int{1, 2, 3, 4, 5}
			Expect(xs).To(Equal(expected))
		})

		It("return 0 value", func() {
			x := 1
			f := func() { x += 1 }
			Apply(f, []interface{}{})
			Expect(x).To(Equal(2))
		})

		It("return 1 value", func() {
			actual := Apply(Reverse, []interface{}{Range(5)})
			expected := []interface{}{5,4,3,2,1}
			Expect(actual).To(Equal(expected))
		})

		It("return 2 values", func() {
			f := func() int {return 3}
			xs := Apply(Timing, []interface{}{f})
			expected := 3
			Expect(Last(xs)).To(Equal(expected))
		})

		It("return 3 values", func() {
			f := func() (int, int, int) {return 1,2,3}
			xs := Apply(f, []interface{}{})
			expected := []interface{}{1,2,3}
			Expect(xs).To(Equal(expected))
		})
	})

	Context("Construct(f, expr)", func() {
		It("Range", func() {
			xs := Construct(Range, 5)
			expected := []int{1, 2, 3, 4, 5}
			Expect(xs).To(Equal(expected))
		})

		It("sum", func() {
			sum := func(nums ...int) (int, int, int) {
				total := 0
				for _, num := range nums {
					total += num
				}
				n := len(nums)
				mean := total / n
				return n, mean, total
			}
			xs := Construct(sum, 1,2,3,4,5)
			expected := []interface{}{5,3,15}
			Expect(xs).To(Equal(expected))
		})
	})

	Context("Composition(f, g, h...)", func() {
		It("0 function", func() {
			f := Composition()
			actual := f(4)
			expected := 4
			Expect(actual).To(Equal(expected))
		})

		It("1 function", func() {
			s := 1
			g := func(x int) { s = x}
			f := Composition(g)
			f(4)
			expected := 4
			Expect(s).To(Equal(expected))
		})

		It("1 function", func() {
			f := Composition(Range)
			actual := f(4)
			expected := []int{1,2,3,4}
			Expect(actual).To(Equal(expected))
		})

		It("2 functions", func() {
			f := Composition(Reverse, Range)
			actual := f(4)
			expected := []interface{}{4,3,2,1}
			Expect(actual).To(Equal(expected))
		})

		It("3 functions", func() {
			f := Composition(Reverse, Reverse, Range)
			actual := f(4)
			expected := []interface{}{1,2,3,4}
			Expect(actual).To(Equal(expected))
		})
	})

	Context("Bind(f, g, h...)", func() {
		It("functions", func() {
			f1 := func(a, b int) (int, int, error) {
				return a, a + b, nil
			}

			f2 := func(n int) (string, error) {
				return strconv.Itoa(n * n), nil
			}

			f3 := func(s string) (int, int, error) {
				n, _ := strconv.Atoi(s)
				return n, n * n, nil
			}

			f4 := func(a, b int) (int, error) {
				return a + b, nil
			}

			f := Bind(f1, f1, f4, f2, f3, f4)
			actual, err := f(1,2)
			expected := 650
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expected))
		})

	})
})
