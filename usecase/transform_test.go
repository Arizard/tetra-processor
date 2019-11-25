package usecase

import (
	"testing"

	"github.com/arizard/tetra"
	"github.com/magiconair/properties/assert"
)

var sampleJSON = `{
    "comma": 44,
    "fields_per_record": -1,
    "transforms": [
        {
            "operation": "slice_rows",
            "kw_args": {
                "start": 1,
                "end": -1
            }
        }
    ]
}`

func TestTransformCSVJsonConfig(t *testing.T) {
	removeFirstRowCfg := tetra.Config{}
	removeFirstRowCfg.LoadFromJSON([]byte(sampleJSON))
	type args struct {
		tetraCfgGetter func() tetra.Config
	}
	tests := []struct {
		name string
		args args
		give string
		want string
	}{
		{
			"Remove first row.",
			args{
				func() tetra.Config {
					return removeFirstRowCfg
				},
			},
			"a,b,c,\nd,e,f,\ng,h,i,\n",
			"d,e,f,\ng,h,i,\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resultCollector string
			_ = resultCollector

			testStringCSVLoader := func() string {
				return tt.give
			}

			testStringCSVSaver := func(s string) {
				resultCollector = s
			}

			TransformCSV(tt.args.tetraCfgGetter, testStringCSVLoader, testStringCSVSaver)

			assert.Equal(t, resultCollector, tt.want)
		})
	}
}
