package qp_test
import (
	"github.com/darinpp/qp"
	"github.com/draffensperger/golp/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)


// min f(x) = 1/2 |x0, x1||4, -2||x0|+|6,0||x0|
//                        |-2, 4||x1|      |x1|
func TestSimple(t *testing.T) {
  //  [4, -2]
  //  [-2, 4]
  G := qp.NewMatrix()
  G.Set([]float64{4, -2, -2, 4},2,2)

  // [6, 0]
  g0 := qp.NewVector()
  g0.Set([]float64{6,0}, 2)

  // [1]
  // [1]
  CE := qp.NewMatrix()
  CE.Set([]float64{1,1},2,1)

  // [-3]
	ce0 := qp.NewVector()
	ce0.Set([]float64{-3},1)

	// [1,0,1]
	// [0,1,1]
  CI := qp.NewMatrix()
	CI.Set([]float64{1,0,1,0,1,1}, 2, 3)

  // [0,0,-2]
	ci0 := qp.NewVector()
	ci0.Set([]float64{0,0,-2}, 3)

	x := qp.NewVector()

	res := qp.Solve_quadprog(G, g0, CE, ce0, CI, ci0, x)

	assert.InEpsilon(t, 12.0, res, 1e-10)
	assert.InEpsilon(t, 1.0, x.At(0), 1e-10)
	assert.InEpsilon(t, 2.0, x.At(1), 1e-10)
}
