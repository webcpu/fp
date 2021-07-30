package test


import (
	"fmt"
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("files", func() {
	Context("Timing(expr)", func() {
		It("function with return", func() {
			f := func() []int {
				xs := Range(0, 1000)
				var result []int = Map(func(x int) int { return x + 1 }, xs).([]int)
				return result
			}
			t1, actual := Timing(f)
			_, expected := 0.003, Range(1, 1001)
			Expect(t1 > 0).To(BeTrue())
			Expect(actual).To(Equal(expected))
		})

		It("function with return", func() {
			f := func() {
				xs := Range(0, 10000)
				_ = Map(func(x int) int { return x + 1 }, xs).([]int)
			}
			t1, actual := Timing(f)
			Expect(t1 > 0).To(BeTrue())
			Expect(actual).To(BeNil())
		})

		It("function with return", func() {
			f := func(y int) {
				fmt.Println(y)
				xs := Range(0, 10000)
				_ = Map(func(x int) int { return x + 1 }, xs).([]int)
			}
			Ω(func() { Timing(f) }).Should(Panic())
		})
	})
})