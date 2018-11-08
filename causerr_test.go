package causerr_test

import (
	stderr "errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/willeponken/causerr"
)

type fakeState struct {
	*gbytes.Buffer
	flag int
}

func (f *fakeState) Precision() (int, bool) {
	return -1, false
}

func (f *fakeState) Flag(c int) bool {
	return (f.flag == c)
}

func (f *fakeState) Width() (int, bool) {
	return -1, false
}

func newFakeState(flag int) *fakeState {
	state := &fakeState{
		Buffer: gbytes.NewBuffer(),
		flag:   flag,
	}
	return state
}

var _ = Describe("Error", func() {
	Describe("error.New", func() {
		Context("With valid ID (>=0), error cause and message", func() {
			var err error

			errCause := stderr.New("err")
			errValidID := 0
			errMessage := "error"

			BeforeEach(func() {
				err = New(errValidID, errCause, errMessage)
			})

			It("should fullfill error interface", func() {
				Expect(err.Error()).To(ContainSubstring(
					"%v (%d: %s)",
					errCause.Error(), errValidID, errMessage,
				))
			})

			It("should fullfill fmt.Formatter interface", func() {
				say := `#0: error`

				formatter, ok := err.(fmt.Formatter)
				Expect(ok).To(Equal(true))

				statePlus := newFakeState('+')
				state := newFakeState('-')

				tbl := []struct {
					state *fakeState
					verb  rune
					say   string
				}{
					{
						state: statePlus,
						verb:  'v',
						say:   say,
					},
					{
						state: state,
						verb:  'v',
						say:   say,
					},
					{
						state: state,
						verb:  's',
						say:   say,
					},
					{
						state: state,
						verb:  'q',
						say:   say,
					},
				}

				for _, test := range tbl {
					formatter.Format(test.state, test.verb)
					Eventually(test.state.Buffer).Should(gbytes.Say(test.say))
				}
			})

			It("should work with ID", func() {
				Expect(ID(err)).To(Equal(errValidID))
			})

			It("should work with Message", func() {
				Expect(Message(err)).To(Equal(errMessage))
			})

			It("should work with Cause", func() {
				Expect(Cause(err).Error()).To(Equal(errCause.Error()))
			})
		})

		Context("With valid ID (>=0), string cause and message", func() {
			var err error

			errCause := "err"
			errValidID := 0
			errMessage := "error"

			BeforeEach(func() {
				err = New(errValidID, errCause, errMessage)
			})

			It("should fullfill error interface", func() {
				Expect(err.Error()).To(ContainSubstring(
					"%v (%d: %s)",
					errCause, errValidID, errMessage,
				))
			})

			It("should work with ID", func() {
				Expect(ID(err)).To(Equal(errValidID))
			})

			It("should work with Message", func() {
				Expect(Message(err)).To(Equal(errMessage))
			})

			It("should work with Cause", func() {
				Expect(Cause(err).Error()).To(Equal(errCause))
			})
		})

		Context("Invalid ID (<0)", func() {
			It("should panic", func() {
				Expect(func() { New(-1, "", "") }).To(Panic())
			})
		})

		Context("Invalid type of cause (not error or string)", func() {
			It("should panic", func() {
				Expect(func() { New(0, []byte{0x0}, "") }).To(Panic())
			})
		})
	})

	Describe("error.ID", func() {
		Context("With a standard error", func() {
			Specify("-1 integer is returned", func() {
				Expect(ID(stderr.New("invalid error"))).To(Equal(-1))
			})
		})
	})

	Describe("error.Cause", func() {
		Context("With a standard error", func() {
			Specify("nil is returned", func() {
				Expect(Cause(stderr.New("invalid error"))).To(BeNil())
			})
		})
	})

	Describe("error.Message", func() {
		Context("With a standard error", func() {
			Specify("empty string is returned", func() {
				Expect(Message(stderr.New("invalid error"))).To(Equal(""))
			})
		})
	})
})
