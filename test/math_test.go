package test

import (
	. "fp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("files", func() {
	Context("Max(expr)", func() {
		It("...int", func() {
			actual := Max(1,2)
			expected := 2
			Expect(actual).To(Equal(expected))
		})

		It("...int", func() {
			actual := Max(1,2,3)
			expected := 3
			Expect(actual).To(Equal(expected))
		})

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
			actual := Max(1, 2, 3)
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
		It("...int", func() {
			actual := Min(1,2)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("...int", func() {
			actual := Min(1,2,3)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

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
			actual := Min(1, 2, 3)
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