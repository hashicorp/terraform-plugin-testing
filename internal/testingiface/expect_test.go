package testingiface_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/internal/testingiface"
)

func TestExpectFail(t *testing.T) {
	t.Parallel()

	testingiface.ExpectPass(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectFail(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Fatal("test fatal")
		})
	})

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectFail(mockT1, func(_ *testingiface.MockT) {
			// intentionally no test error or test skip
		})
	})

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectFail(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Skip("test skip")
		})
	})
}

func TestExpectParallel(t *testing.T) {
	t.Parallel()

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectParallel(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Fatal("test fatal")
		})
	})

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectParallel(mockT1, func(_ *testingiface.MockT) {
			// intentionally no test error or test skip
		})
	})

	testingiface.ExpectPass(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectParallel(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Parallel()
		})
	})

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectParallel(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Skip("test skip")
		})
	})
}

func TestExpectPass(t *testing.T) {
	t.Parallel()

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectPass(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Fatal("test fatal")
		})
	})

	testingiface.ExpectPass(t, func(_ *testingiface.MockT) {
		// intentionally no test error or test skip
	})

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectPass(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Skip("test skip")
		})
	})
}

func TestExpectSkip(t *testing.T) {
	t.Parallel()

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectSkip(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Fatal("test fatal")
		})
	})

	testingiface.ExpectFail(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectSkip(mockT1, func(_ *testingiface.MockT) {
			// intentionally no test error or test skip
		})
	})

	testingiface.ExpectPass(t, func(mockT1 *testingiface.MockT) {
		testingiface.ExpectSkip(mockT1, func(mockT2 *testingiface.MockT) {
			mockT2.Skip("test skip")
		})
	})
}
