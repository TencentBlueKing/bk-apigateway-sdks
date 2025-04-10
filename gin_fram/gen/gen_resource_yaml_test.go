package gen

import "testing"

func TestGenResourceYamlFromSwaggerJson(t *testing.T) {
	type args struct {
		docPath string
		engine  *gin.Engine
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenResourceYamlFromSwaggerJson(tt.args.docPath, tt.args.engine); got != tt.want {
				t.Errorf("GenResourceYamlFromSwaggerJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
