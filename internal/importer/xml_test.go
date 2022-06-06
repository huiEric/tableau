package importer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
)

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
		{
			name: "FeatureToggle",
			args: args{
				doc: `
<Conf>
	<Client EnvID="map<uint32,Client>">
		<Toggle ID="map<enum<.ToggleType>, Toggle>" WorldID="uint32"/>
	</Client>
</Conf>
`,
			},
			want: `
<Conf>
	<Client EnvID="map&lt;uint32,Client&gt;">
		<Toggle ID="map&lt;enum&lt;.ToggleType&gt;, Toggle&gt;" WorldID="uint32"/>
	</Client>
</Conf>
`,
		},
		{
			name: "Prop",
			args: args{
				doc: `
<Conf>
	<Client ID="map<uint32, Client>|{unique:true range:"1,~"}" OpenTime="datetime|{default:"2022-01-23 15:40:00"}"/>
</Conf>
`,
			},
			want: `
<Conf>
	<Client ID="map&lt;uint32, Client&gt;|{unique:true range:&#34;1,~&#34;}" OpenTime="datetime|{default:&#34;2022-01-23 15:40:00&#34;}"/>
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

func Test_isRepeated(t *testing.T) {
	doc := `
<?xml version='1.0' encoding='UTF-8'?>
<MatchCfg Open="true">
	<TeamRatingWeight>
		<Weight Num="1">
			<Param Value="100"/>
		</Weight>
		<Weight Num="2">
			<Param Value="30"/>
			<Param Value="70"/>
		</Weight>
	</TeamRatingWeight>
</MatchCfg>
`
	root, _ := xmlquery.Parse(strings.NewReader(doc))
	node1 := xmlquery.FindOne(root, "MatchCfg/TeamRatingWeight/Weight")
	node2 := xmlquery.FindOne(root, "MatchCfg/TeamRatingWeight/Weight/Param")
	node3 := xmlquery.FindOne(root, "MatchCfg/TeamRatingWeight")
	node4 := xmlquery.FindOne(root, "MatchCfg")
	type args struct {
		root, curr *xmlquery.Node
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "node1",
			args: args{
				root: root,
				curr: node1,
			},
			want: true,
		},
		{
			name: "node2",
			args: args{
				root: root,
				curr: node2,
			},
			want: true,
		},
		{
			name: "node3",
			args: args{
				root: root,
				curr: node3,
			},
			want: false,
		},
		{
			name: "sheet attr",
			args: args{
				root: root,
				curr: node4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRepeated(tt.args.root, tt.args.curr); got != tt.want {
				t.Errorf("isRepeated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchAttr(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "scalar type",
			args: args{
				s: `<AAA bb="bool" cc="int64" dd="enum<.EnumType>" >`,
			},
			want: []string{
				`bb="bool"`, `bb`, `bool`, ``,
			},
		},
		// TODO: Add test cases.
		{
			name: "Prop",
			args: args{
				s: `<Client OpenTime="datetime|{default:"2022-01-23 15:40:00"}" CloseTime="datetime|{default:"2022-01-23 15:40:00"}"/>`,
			},
			want: []string{
				`OpenTime="datetime|{default:"2022-01-23 15:40:00"}"`,
				`OpenTime`, `datetime`, `|{default:"2022-01-23 15:40:00"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchAttr(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFirstChild(t *testing.T) {
	doc := `
<?xml version='1.0' encoding='UTF-8'?>
<Server>
    <MapConf>
        <Weight Num="map&lt;uint32,Weight&gt;"/>
    </MapConf>
</Server>
`
	root, _ := xmlquery.Parse(strings.NewReader(doc))
	node1 := xmlquery.FindOne(root, "Server/MapConf/Weight")
	type args struct {
		curr *xmlquery.Node
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Weight",
			args: args{
				curr: node1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFirstChild(tt.args.curr); got != tt.want {
				t.Errorf("isFirstChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_correctType(t *testing.T) {
	doc := `
<?xml version='1.0' encoding='UTF-8'?>
<MatchCfg open="true">
	<MatchMode MissionType="map&lt;enum&lt;.MissionType&gt;,MatchMode&gt;">
		<MatchAI IsOpen="bool" PlayerOnlyOneCamp="bool">
			<AI Type="[AI]&lt;enum&lt;.ENMAIWarmType&gt;&gt;" IsOpen="bool" MinTime="duration" MaxTime="duration" />
		</MatchAI>
    </MatchMode>
	<MapConf>
        <Weight Num="map&lt;uint32,Weight&gt;"/>
    </MapConf>
	<Client EnvID="map&lt;uint32,Client&gt;">
		<Toggle ID="map&lt;enum&lt;.ToggleType&gt;, Toggle&gt;" WorldID="uint32"/>
	</Client>
</MatchCfg>

<MatchCfg>
	<StructConf>
		<Weight Num="1">
			<Param Value="100"/>
		</Weight>
	</StructConf>

	<ListConf>
        <Weight Num="1">
            <Param Value="100"/>
        </Weight>
        <Weight Num="2">
            <Param Value="30"/>
            <Param Value="70"/>
        </Weight>
    </ListConf>
</MatchCfg>
`
	root, _ := xmlquery.Parse(strings.NewReader(doc))
	node1 := xmlquery.FindOne(root, "MatchCfg/MatchMode/MatchAI/AI")
	node2 := xmlquery.FindOne(root, "MatchCfg/MapConf/Weight")
	node3 := xmlquery.FindOne(root, "MatchCfg/StructConf/Weight")
	node4 := xmlquery.FindOne(root, "MatchCfg/ListConf/Weight")
	node5 := xmlquery.FindOne(root, "MatchCfg/ListConf/Weight/Param")
	node6 := xmlquery.FindOne(root, "MatchCfg/Client/Toggle")
	node7 := xmlquery.FindOne(root, "MatchCfg")
	type args struct {
		root    *xmlquery.Node
		curr    *xmlquery.Node
		oriType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "MatchAI",
			args: args{
				root:    root,
				curr:    node1,
				oriType: `[AI]<enum<.ENMAIWarmType>>`,
			},
			want: `[AI]<enum<.ENMAIWarmType>>`,
		},
		{
			name: "MapConf",
			args: args{
				root:    root,
				curr:    node2,
				oriType: `map<uint32,Weight>`,
			},
			want: `{MapConf}map<uint32,Weight>`,
		},
		{
			name: "StructConf",
			args: args{
				root:    root,
				curr:    node3,
				oriType: `int32`,
			},
			want: `{StructConf}{Weight}int32`,
		},
		{
			name: "ListConf",
			args: args{
				root:    root,
				curr:    node4,
				oriType: `int32`,
			},
			want: `{ListConf}[Weight]<int32>`,
		},
		{
			name: "ListConf/Param",
			args: args{
				root:    root,
				curr:    node5,
				oriType: `int32`,
			},
			want: `[Param]<int32>`,
		},
		{
			name: "FeatureToggle",
			args: args{
				root:    root,
				curr:    node6,
				oriType: `map<enum<.ToggleType>, Toggle>`,
			},
			want: `map<enum<.ToggleType>, Toggle>`,
		},
		{
			name: "sheet attr",
			args: args{
				root:    root,
				curr:    node7,
				oriType: `bool`,
			},
			want: `bool`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := correctType(tt.args.root, tt.args.curr, tt.args.oriType); got != tt.want {
				t.Errorf("correctType() = %v, want %v", got, tt.want)
			}
		})
	}
}
