/**
 * Created by PhpStorm.
 * @file   matchmaker.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 1106
 * @desc   matchmaker.go
 */

package params

type AddMatchMakerParams struct {
	MinCount          int
	MaxCount          int
	Query             string
	StringProperties  map[string]string
	NumericProperties map[string]float64
}
