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

	Context("Abs(expr)", func() {
		It("int", func() {
			actual := Abs(1)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("int", func() {
			actual := Abs(-1)
			expected := 1
			Expect(actual).To(Equal(expected))
		})

		It("int8", func() {
			actual := Abs(int8(-1))
			expected := int8(1)
			Expect(actual).To(Equal(expected))
		})

		It("int16", func() {
			actual := Abs(int16(-1))
			expected := int16(1)
			Expect(actual).To(Equal(expected))
		})

		It("int32", func() {
			actual := Abs(int32(-1))
			expected := int32(1)
			Expect(actual).To(Equal(expected))
		})

		It("int64", func() {
			actual := Abs(int64(-1))
			expected := int64(1)
			Expect(actual).To(Equal(expected))
		})

		It("uint8", func() {
			actual := Abs(uint8(1))
			expected := uint8(1)
			Expect(actual).To(Equal(expected))
		})

		It("uint16", func() {
			actual := Abs(uint16(1))
			expected := uint16(1)
			Expect(actual).To(Equal(expected))
		})

		It("uint32", func() {
			actual := Abs(uint32(1))
			expected := uint32(1)
			Expect(actual).To(Equal(expected))
		})

		It("uint64", func() {
			actual := Abs(uint64(1))
			expected := uint64(1)
			Expect(actual).To(Equal(expected))
		})

		It("complex64", func() {
			var x complex64 = 3 + 4i
			actual := Abs(x)
			expected := float64(5)
			Expect(actual).To(Equal(expected))
		})

		It("string", func() {
			Ω(func(){Abs("abc")}).Should(Panic())
		})
	})

	Context("Pow(x, y)", func() {
		It("int", func() {
			actual := Pow(3, 2)
			expected := 9
			Expect(actual).To(Equal(expected))
		})

		It("int", func() {
			actual := Pow(-3, 2)
			expected := 9
			Expect(actual).To(Equal(expected))
		})

		It("int8", func() {
			actual := Pow(int8(-3), int8(2))
			expected := int8(9)
			Expect(actual).To(Equal(expected))
		})

		It("int16", func() {
			actual := Pow(int16(-3), int8(2))
			expected := int16(9)
			Expect(actual).To(Equal(expected))
		})

		It("int32", func() {
			actual := Pow(int32(-3), int8(2))
			expected := int32(9)
			Expect(actual).To(Equal(expected))
		})

		It("int64", func() {
			actual := Pow(int64(-3), int8(2))
			expected := int64(9)
			Expect(actual).To(Equal(expected))
		})

		It("uint8", func() {
			actual := Pow(uint8(3), uint8(2))
			expected := uint8(9)
			Expect(actual).To(Equal(expected))
		})

		It("uint16", func() {
			actual := Pow(uint16(3), uint8(2))
			expected := uint16(9)
			Expect(actual).To(Equal(expected))
		})

		It("uint32", func() {
			actual := Pow(uint32(3), 2)
			expected := uint32(9)
			Expect(actual).To(Equal(expected))
		})

		It("uint64", func() {
			actual := Pow(uint64(3), 2)
			expected := uint64(9)
			Expect(actual).To(Equal(expected))
		})

		It("complex64", func() {
			var x complex64 = 3 + 4i
			var y complex64 = 2
			actual := Pow(x, y)
			expected := complex(float32(-7), float32(24))
			Expect(actual).To(Equal(expected))
		})

		It("complex64", func() {
			var x complex64 = 3 + 4i
			var y int = 2
			actual := Pow(x, y)
			expected := complex(float32(-7), float32(24))
			Expect(actual).To(Equal(expected))
		})

		It("string", func() {
			Ω(func(){Pow("abc", 2)}).Should(Panic())
		})
	})
})