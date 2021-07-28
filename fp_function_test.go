package fp_test

import (
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("function", func() {
	Context("Apply(f, expr)", func() {
		It("Range", func() {
			xs := Apply(Range, []interface{}{5})
			expected := []int{1, 2, 3, 4, 5}
			Expect(xs).To(Equal(expected))
		})

		It("reverse", func() {
			xs := Apply(Reverse, []interface{}{Range(5)})
			expected := Reverse([]int{1, 2, 3, 4, 5})
			Expect(xs).To(Equal(expected))
		})

		It("reverse", func() {
			Ω(Apply(Reverse, Range(5))).Should(Panic())
		})
	})
})
