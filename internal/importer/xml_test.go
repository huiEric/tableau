package importer

import (
	"testing"
)

func Test_isMetaBeginning(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "standard",
			args: args{
				s: "<!-- @TABLEAU\n",
			},
			want: true,
		},
		{
			name: "with spaces",
			args: args{
				s: "   <!--      @TABLEAU           \n",
			},
			want: true,
		},
		{
			name: "none space",
			args: args{
				s: "<!--@TABLEAU\n",
			},
			want: true,
		},
		{
			name: "non-space heading",
			args: args{
				s: "test <!-- @TABLEAU\n",
			},
			want: false,
		},
		{
			name: "non-space tailing",
			args: args{
				s: "<!-- @TABLEAU <conf>\n",
			},
			want: false,
		},
		{
			name: "empty metaSheet with more spaces",
			args: args{
				s: "<!-- @TABLEAU -->\n",
			},
			want: true,
		},
		{
			name: "empty metaSheet with less spaces",
			args: args{
				s: " <!--@TABLEAU -->   \n",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMetaBeginning(tt.args.s); got != tt.want {
				t.Errorf("isMetaBeginning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isMetaEnding(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "standard",
			args: args{
				s: "-->\n",
			},
			want: true,
		},
		{
			name: "more spaces",
			args: args{
				s: "  --> \n",
			},
			want: true,
		},
		{
			name: "empty metaSheet with more spaces",
			args: args{
				s: "<!-- @TABLEAU -->\n",
			},
			want: true,
		},
		{
			name: "empty metaSheet with less spaces",
			args: args{
				s: " <!--@TABLEAU -->   \n",
			},
			want: true,
		},
		{
			name: "non-space characters heading",
			args: args{
				s: "</Conf> -->\n",
			},
			want: false,
		},
		{
			name: "non-space characters tailing",
			args: args{
				s: " --> <Conf> \n",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMetaEnding(tt.args.s); got != tt.want {
				t.Errorf("isMetaEnding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMetaDoc(t *testing.T) {
	type args struct {
		doc string
	}
	tests := []struct {
		name        string
		args        args
		wantMetaDoc string
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "standard",
			args: args{
				doc: `
<!-- @TABLEAU
<Conf>
    <Server Value="int32"/>
</Conf>
-->

<Conf>
    <Server Value="100"/>
</Conf>
`,
			},
			wantMetaDoc: `<Conf>
    <Server Value="int32"/>
</Conf>
`,
			wantErr: false,
		},
		{
			name: "empty metaSheet",
			args: args{
				doc: `
<!--@TABLEAU -->

<Conf>
    <Server Value="100"/>
</Conf>
`,
			},
			wantMetaDoc: ``,
			wantErr:     false,
		},
		{
			name: "none metaSheet",
			args: args{
				doc: `

<Conf>
    <Server Value="100"/>
</Conf>
`,
			},
			wantMetaDoc: ``,
			wantErr:     true,
		},
		{
			name: "standard",
			args: args{
				doc: `
<!-- @TABLEAU <Conf>
    <Server Value="int32"/>
</Conf>
-->

<Conf>
    <Server Value="100"/>
</Conf>
`,
			},
			wantMetaDoc: ``,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMetaDoc, err := getMetaDoc(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMetaDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMetaDoc != tt.wantMetaDoc {
				t.Errorf("getMetaDoc() = %v, want %v", gotMetaDoc, tt.wantMetaDoc)
			}
		})
	}
}

func Test_escapeAttrs(t *testing.T) {
	type args struct {
		doc string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "standard",
			args: args{
				doc: `
<Conf>
    <Server Type="map<enum<.ServerType>, Server>" Value="int32"/>
</Conf>
`,
			},
			want: `
<Conf>
    <Server Type="map&lt;enum&lt;.ServerType&gt;, Server&gt;" Value="int32"/>
</Conf>
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeAttrs(tt.args.doc); got != tt.want {
				t.Errorf("escapeAttrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
