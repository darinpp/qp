%module Array
%{
#include "Array.hh"
#include "QuadProg++.hh"
%}

%include <typemaps.i>
%apply double *INOUT { const double* a };

namespace quadprogpp {
  template <itypename T> class Matrix {
  public:
    Matrix();
    ~Matrix();  
  
    inline void set(const T* a, unsigned int n, unsigned int m);

    inline unsigned int nrows() const { return n; } // number of rows
    inline unsigned int ncols() const { return m; } // number of columns
  };

  template <typename T> class Vector {
  public:
    Vector();
    ~Vector();

    inline void set(const T* a, const unsigned int n);
   
    inline unsigned int size() const;
  };

  double solve_quadprog(Matrix<double>& G, Vector<double>& g0,
                      const Matrix<double>& CE, const Vector<double>& ce0,
                      const Matrix<double>& CI, const Vector<double>& ci0,
                      Vector<double>& x);
};
%template(matrix) quadprogpp::Matrix<double>;
%template(vector) quadprogpp::Vector<double>;
%extend quadprogpp::Vector<double> {
  double at(unsigned int i) {
    return (*$self)[i];
  }
}
