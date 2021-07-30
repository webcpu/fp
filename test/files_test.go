package test

import (
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("files", func() {
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
	})
})
