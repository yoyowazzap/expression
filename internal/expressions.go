package internal

type (
	Path []interface{}

	PathParser interface {
		GetValue(Path) (interface{}, bool)
		GetNumber(Path) (float64, bool)
		GetBoolean(Path) (bool, bool)
		GetString(Path) (string, bool)
		GetArray(Path) ([]interface{}, bool)
	}

	Expression interface {
		Value(PathParser) interface{}
		Reduce() Expression
	}

	NumberExpression interface {
		Value(PathParser) float64
		Reduce() NumberExpression
	}

	BooleanExpression interface {
		Value(PathParser) bool
		Reduce() BooleanExpression
	}

	StringExpression interface {
		Value(PathParser) string
	}

	ArrayExpression interface {
		Value(PathParser) []interface{}
	}

	generic struct {
		n NumberExpression
		b BooleanExpression
		s StringExpression
		a ArrayExpression
	}

	genericPath struct {
		path []interface{}
	}

	genericPathWithDefault struct {
		path         []interface{}
		defaultValue interface{}
	}

	number struct {
		n float64
	}

	numberPath struct {
		path []interface{}
	}

	numberPathWithDefault struct {
		path         []interface{}
		defaultValue NumberExpression
	}

	inverseExpression struct {
		subExpression NumberExpression
	}

	sumExpression struct {
		subExpressions []NumberExpression
	}

	subtractExpression struct {
		e1 NumberExpression
		e2 NumberExpression
	}

	timesExpression struct {
		subExpressions []NumberExpression
	}

	divideExpression struct {
		e1 NumberExpression
		e2 NumberExpression
	}

	lengthExpression struct {
		ae ArrayExpression
	}

	boolean struct {
		b bool
	}

	booleanPath struct {
		path []interface{}
	}

	booleanPathWithDefault struct {
		path         []interface{}
		defaultValue BooleanExpression
	}

	notExpression struct {
		subExpression BooleanExpression
	}

	lessThanExpression struct {
		e1 NumberExpression
		e2 NumberExpression
	}

	lessThanOrEqualExpression struct {
		e1 NumberExpression
		e2 NumberExpression
	}

	greaterThanExpression struct {
		e1 NumberExpression
		e2 NumberExpression
	}

	greaterThanOrEqualExpression struct {
		e1 NumberExpression
		e2 NumberExpression
	}

	equalExpression struct {
		e1 Expression
		e2 Expression
	}

	andExpression struct {
		subExpressions []BooleanExpression
	}

	orExpression struct {
		subExpressions []BooleanExpression
	}

	str struct {
		s string
	}

	strPath struct {
		path []interface{}
	}

	strPathWithDefault struct {
		path         []interface{}
		defaultValue StringExpression
	}

	arrayPath struct {
		path []interface{}
	}
)

func (e *generic) Value(pp PathParser) interface{} {
	switch {
	case e.n != nil:
		return e.n.Value(pp)
	case e.b != nil:
		return e.b.Value(pp)
	case e.s != nil:
		return e.s.Value(pp)
	case e.a != nil:
		return e.a.Value(pp)
	}
	return nil
}

func (e *generic) Reduce() Expression {
	if e.n != nil {
		e.n = e.n.Reduce()
	}
	if e.b != nil {
		e.b = e.b.Reduce()
	}
	return e
}

func (gp *genericPath) Value(pp PathParser) interface{} {
	value, _ := pp.GetValue(gp.path)
	return value
}

func (gpd *genericPathWithDefault) Value(pp PathParser) interface{} {
	value, ok := pp.GetArray(gpd.path)
	if !ok {
		return gpd.defaultValue
	}
	return value
}

func (gpd *genericPathWithDefault) Reduce() Expression {
	return gpd
}

func (n *number) Value(PathParser) float64 {
	return n.n
}

func (n *number) Reduce() NumberExpression {
	return n
}

func (np *numberPath) Value(pp PathParser) float64 {
	num, _ := pp.GetNumber(np.path)
	return num
}

func (np *numberPath) Reduce() NumberExpression {
	return np
}

func (npd *numberPathWithDefault) Value(pp PathParser) float64 {
	num, ok := pp.GetNumber(npd.path)
	if !ok {
		return npd.defaultValue.Value(pp)
	}
	return num
}

func (npd *numberPathWithDefault) Reduce() NumberExpression {
	npd.defaultValue = npd.defaultValue.Reduce()
	return npd
}

func (ie *inverseExpression) Value(pp PathParser) float64 {
	return -ie.subExpression.Value(pp)
}

func (ie *inverseExpression) Reduce() NumberExpression {
	ie.subExpression = ie.subExpression.Reduce()
	if numExpr, ok := ie.subExpression.(*number); ok {
		return &number{n: -numExpr.n}
	}
	return ie
}

func (se *sumExpression) Value(pp PathParser) float64 {
	var sum float64
	for _, subExpression := range se.subExpressions {
		sum += subExpression.Value(pp)
	}
	return sum
}

func (se *sumExpression) Reduce() NumberExpression {
	var sum float64
	var subExpressions []NumberExpression
	for _, subExpression := range se.subExpressions {
		reducedSubExpression := subExpression.Reduce()
		if numExpr, ok := reducedSubExpression.(*number); ok {
			sum += numExpr.n
		} else if sumExpr, ok := reducedSubExpression.(*sumExpression); ok {
			for _, subReducedExpression := range sumExpr.subExpressions {
				if subNumExpr, ok := subReducedExpression.(*number); ok {
					sum += subNumExpr.n
				} else {
					subExpressions = append(subExpressions, subReducedExpression)
				}
			}
		} else {
			subExpressions = append(subExpressions, reducedSubExpression)
		}
	}
	if len(subExpressions) == 0 {
		return &number{n: sum}
	}
	if sum != 0 {
		subExpressions = append(subExpressions, &number{n: sum})
	}
	se.subExpressions = subExpressions
	return se
}

func (se *subtractExpression) Value(pp PathParser) float64 {
	return se.e1.Value(pp) - se.e2.Value(pp)
}

func (se *subtractExpression) Reduce() NumberExpression {
	se.e1 = se.e1.Reduce()
	se.e2 = se.e2.Reduce()

	numExpr1, ok1 := se.e1.(*number)
	numExpr2, ok2 := se.e2.(*number)
	if ok1 && ok2 {
		return &number{n: numExpr1.n - numExpr2.n}
	}

	return se
}

func (te *timesExpression) Value(pp PathParser) float64 {
	var product float64 = 1
	for _, subExpression := range te.subExpressions {
		product *= subExpression.Value(pp)
	}
	return product
}

func (te *timesExpression) Reduce() NumberExpression {
	var product float64 = 1
	var subExpressions []NumberExpression
	for _, subExpression := range te.subExpressions {
		reducedSubExpression := subExpression.Reduce()
		if numExpr, ok := reducedSubExpression.(*number); ok {
			product *= numExpr.n
		} else if timesExpr, ok := reducedSubExpression.(*timesExpression); ok {
			for _, subReducedExpression := range timesExpr.subExpressions {
				if subNumExpr, ok := subReducedExpression.(*number); ok {
					product *= subNumExpr.n
				} else {
					subExpressions = append(subExpressions, subReducedExpression)
				}
			}
		} else {
			subExpressions = append(subExpressions, reducedSubExpression)
		}
	}
	if len(subExpressions) == 0 {
		return &number{n: product}
	}
	if product != 1 {
		subExpressions = append(subExpressions, &number{n: product})
	}
	te.subExpressions = subExpressions
	return te
}

func (de *divideExpression) Value(pp PathParser) float64 {
	return de.e1.Value(pp) / de.e2.Value(pp)
}

func (de *divideExpression) Reduce() NumberExpression {
	de.e1 = de.e1.Reduce()
	de.e2 = de.e2.Reduce()

	numExpr1, ok1 := de.e1.(*number)
	numExpr2, ok2 := de.e2.(*number)
	if ok1 && ok2 {
		return &number{n: numExpr1.n / numExpr2.n}
	}

	return de
}

func (le *lengthExpression) Value(pp PathParser) float64 {
	return float64(len(le.ae.Value(pp)))
}

func (le *lengthExpression) Reduce() NumberExpression {
	return le
}

func (b *boolean) Value(PathParser) bool {
	return b.b
}

func (b *boolean) Reduce() BooleanExpression {
	return b
}

func (bp *booleanPath) Value(pp PathParser) bool {
	value, _ := pp.GetBoolean(bp.path)
	return value
}

func (bp *booleanPath) Reduce() BooleanExpression {
	return bp
}

func (bpd *booleanPathWithDefault) Value(pp PathParser) bool {
	value, ok := pp.GetBoolean(bpd.path)
	if !ok {
		return bpd.defaultValue.Value(pp)
	}
	return value
}

func (bpd *booleanPathWithDefault) Reduce() BooleanExpression {
	bpd.defaultValue = bpd.defaultValue.Reduce()
	return bpd
}

func (ne *notExpression) Value(pp PathParser) bool {
	return !ne.Value(pp)
}

func (ne *notExpression) Reduce() BooleanExpression {
	ne.subExpression = ne.subExpression.Reduce()
	if boolExpr, ok := ne.subExpression.(*boolean); ok {
		return &boolean{b: !boolExpr.b}
	}
	return ne
}

func (e *lessThanExpression) Value(pp PathParser) bool {
	return e.e1.Value(pp) < e.e2.Value(pp)
}

func (e *lessThanExpression) Reduce() BooleanExpression {
	e.e1 = e.e1.Reduce()
	e.e2 = e.e2.Reduce()

	numExpr1, ok1 := e.e1.(*number)
	numExpr2, ok2 := e.e2.(*number)
	if ok1 && ok2 {
		return &boolean{b: numExpr1.n < numExpr2.n}
	}

	return e
}

func (e *lessThanOrEqualExpression) Value(pp PathParser) bool {
	return e.e1.Value(pp) <= e.e2.Value(pp)
}

func (e *lessThanOrEqualExpression) Reduce() BooleanExpression {
	e.e1 = e.e1.Reduce()
	e.e2 = e.e2.Reduce()

	numExpr1, ok1 := e.e1.(*number)
	numExpr2, ok2 := e.e2.(*number)
	if ok1 && ok2 {
		return &boolean{b: numExpr1.n <= numExpr2.n}
	}

	return e
}

func (e *greaterThanExpression) Value(pp PathParser) bool {
	return e.e1.Value(pp) > e.e2.Value(pp)
}

func (e *greaterThanExpression) Reduce() BooleanExpression {
	e.e1 = e.e1.Reduce()
	e.e2 = e.e2.Reduce()

	numExpr1, ok1 := e.e1.(*number)
	numExpr2, ok2 := e.e2.(*number)
	if ok1 && ok2 {
		return &boolean{b: numExpr1.n > numExpr2.n}
	}

	return e
}

func (e *greaterThanOrEqualExpression) Value(pp PathParser) bool {
	return e.e1.Value(pp) >= e.e2.Value(pp)
}

func (e *greaterThanOrEqualExpression) Reduce() BooleanExpression {
	e.e1 = e.e1.Reduce()
	e.e2 = e.e2.Reduce()

	numExpr1, ok1 := e.e1.(*number)
	numExpr2, ok2 := e.e2.(*number)
	if ok1 && ok2 {
		return &boolean{b: numExpr1.n >= numExpr2.n}
	}

	return e
}

func (e *equalExpression) Value(pp PathParser) bool {
	v1 := e.e1.Value(pp)
	v2 := e.e2.Value(pp)

	str1, ok1 := v1.(string)
	str2, ok2 := v2.(string)
	if ok1 && ok2 {
		return str1 == str2
	}

	num1, ok1 := v1.(float64)
	num2, ok2 := v2.(float64)
	if ok1 && ok2 {
		return num1 == num2
	}

	b1, ok1 := v1.(bool)
	b2, ok2 := v2.(bool)
	if ok1 == ok2 {
		return b1 == b2
	}

	return false
}

func (e *equalExpression) Reduce() BooleanExpression {
	e.e1 = e.e1.Reduce()
	e.e2 = e.e2.Reduce()

	expr1, ok1 := e.e1.(*generic)
	expr2, ok2 := e.e2.(*generic)
	if ok1 && ok2 {
		if expr1.s != nil && expr2.s != nil {
			strExpr1, ok1 := expr1.s.(*str)
			strExpr2, ok2 := expr2.s.(*str)
			if ok1 && ok2 {
				return &boolean{b: strExpr1.s == strExpr2.s}
			}
		}

		if expr1.n != nil && expr2.n != nil {
			numExpr1, ok1 := expr1.n.(*number)
			numExpr2, ok2 := expr2.n.(*number)
			if ok1 && ok2 {
				return &boolean{b: numExpr1.n == numExpr2.n}
			}
		}

		if expr1.b != nil && expr2.b != nil {
			boolExpr1, ok1 := expr1.b.(*boolean)
			boolExpr2, ok2 := expr2.b.(*boolean)
			if ok1 && ok2 {
				return &boolean{b: boolExpr1.b == boolExpr2.b}
			}
		}
	}

	return e
}

func (e *andExpression) Value(pp PathParser) bool {
	for _, subExpression := range e.subExpressions {
		if !subExpression.Value(pp) {
			return false
		}
	}
	return true
}

func (e *andExpression) Reduce() BooleanExpression {
	var subExpressions []BooleanExpression
	for _, subExpression := range e.subExpressions {
		reducedSubExpression := subExpression.Reduce()
		if boolExpr, ok := reducedSubExpression.(*boolean); ok {
			if !boolExpr.b {
				return boolExpr
			}
		} else if andExpr, ok := reducedSubExpression.(*andExpression); ok {
			subExpressions = append(subExpressions, andExpr.subExpressions...)
		} else {
			subExpressions = append(subExpressions, reducedSubExpression)
		}
	}
	if len(subExpressions) == 0 {
		return &boolean{b: true}
	}
	e.subExpressions = subExpressions
	return e
}

func (e *orExpression) Value(pp PathParser) bool {
	for _, subExpression := range e.subExpressions {
		if subExpression.Value(pp) {
			return true
		}
	}
	return false
}

func (e *orExpression) Reduce() BooleanExpression {
	var subExpressions []BooleanExpression
	for _, subExpression := range e.subExpressions {
		reducedSubExpression := subExpression.Reduce()
		if boolExpr, ok := reducedSubExpression.(*boolean); ok {
			if boolExpr.b {
				return boolExpr
			}
		} else if orExpr, ok := reducedSubExpression.(*orExpression); ok {
			subExpressions = append(subExpressions, orExpr.subExpressions...)
		} else {
			subExpressions = append(subExpressions, reducedSubExpression)
		}
	}
	if len(subExpressions) == 0 {
		return &boolean{b: false}
	}
	e.subExpressions = subExpressions
	return e
}

func (e *str) Value(PathParser) string {
	return e.s
}

func (e *str) Reduce() StringExpression {
	return e
}

func (e *strPath) Value(pp PathParser) string {
	value, _ := pp.GetString(e.path)
	return value
}

func (e *strPath) Reduce() StringExpression {
	return e
}

func (e *strPathWithDefault) Value(pp PathParser) string {
	value, ok := pp.GetString(e.path)
	if !ok {
		return e.defaultValue.Value(pp)
	}
	return value
}

func (e *arrayPath) Value(pp PathParser) []interface{} {
	value, _ := pp.GetArray(e.path)
	return value
}
