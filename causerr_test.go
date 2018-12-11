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
		Context("With error cause and message", func() {
			var err error

			errCause := stderr.New("err")
			errMessage := "error"

			BeforeEach(func() {
				err = New(errCause, errMessage)
			})

			It("should fullfill error interface", func() {
				Expect(err.Error()).To(ContainSubstring(
					"%v (%s)",
					errCause.Error(), errMessage,
				))
			})

			It("should fullfill fmt.Formatter interface", func() {
				say := `error`

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

			It("should work with Message", func() {
				Expect(Message(err)).To(Equal(errMessage))
			})

			It("should work with Cause", func() {
				Expect(Cause(err).Error()).To(Equal(errCause.Error()))
			})
		})

		Context("With string cause and message", func() {
			var err error

			errCause := "err"
			errMessage := "error"

			BeforeEach(func() {
				err = New(errCause, errMessage)
			})

			It("should fullfill error interface", func() {
				Expect(err.Error()).To(ContainSubstring(
					"%v (%s)",
					errCause, errMessage,
				))
			})

			It("should work with Message", func() {
				Expect(Message(err)).To(Equal(errMessage))
			})

			It("should work with Cause", func() {
				Expect(Cause(err).Error()).To(Equal(errCause))
			})
		})

		Context("Invalid type of cause (not error or string)", func() {
			It("should panic", func() {
				Expect(func() { New([]byte{0x0}, "") }).To(Panic())
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
