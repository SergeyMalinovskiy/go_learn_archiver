package table

import (
	"awesome-archiver/lib/compression/vlc"
	"reflect"
	"testing"
)

func Test_encodingTable_DecodingTree(t *testing.T) {
	tests := []struct {
		name string
		et   vlc.encodingTable
		want decodingTree
	}{
		{
			name: "base tree test",
			et: vlc.encodingTable{
				'a': "11",
				'b': "1001",
				'z': "0101",
			},
			want: decodingTree{
				Zero: &decodingTree{
					One: &decodingTree{
						Zero: &decodingTree{
							One: &decodingTree{
								Value: "z",
							},
						},
					},
				},
				One: &decodingTree{
					Zero: &decodingTree{
						Zero: &decodingTree{
							One: &decodingTree{
								Value: "b",
							},
						},
					},
					One: &decodingTree{
						Value: "a",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.DecodingTree(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodingTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
