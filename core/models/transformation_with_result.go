package models

type TransformationWithResult struct {
	transformation Transformation
	result         Result
}

func NewTransformationWithResult(
	transformation Transformation,
	result Result,
) *TransformationWithResult {
	return &TransformationWithResult{transformation, result}
}

func (twr TransformationWithResult) GetTransformation() Transformation {
	return twr.transformation
}

func (twr TransformationWithResult) GetResult() Result {
	return twr.result
}
