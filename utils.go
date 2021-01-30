package hkit

/****************
Author: yuansudong
Date: 2018.06.23
*****************/

// _MergeHeader 用于合并一个header
func _MergeHeader(
	iHeaders ...[]HTTPOption,
) []HTTPOption {
	sTarget := []HTTPOption{}
	for _, item := range iHeaders {
		sTarget = append(sTarget, item...)
	}
	return sTarget
}
