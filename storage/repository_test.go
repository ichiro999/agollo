package storage

import (
	"strings"
	"testing"
	"time"

	. "github.com/tevid/gohamcrest"
	_ "github.com/zouyx/agollo/v3/agcache/memory"
	"github.com/zouyx/agollo/v3/env"
	_ "github.com/zouyx/agollo/v3/env/file/json"

	_ "github.com/zouyx/agollo/v3/utils/parse/normal"
	_ "github.com/zouyx/agollo/v3/utils/parse/properties"
)

//init param
func init() {
}

func creatTestApolloConfig(configurations map[string]interface{}, namespace string) {
	apolloConfig := &env.ApolloConfig{}
	apolloConfig.NamespaceName = namespace
	apolloConfig.AppID = "test"
	apolloConfig.Cluster = "dev"
	apolloConfig.Configurations = configurations
	UpdateApolloConfig(apolloConfig, true)

}

func TestUpdateApolloConfigNull(t *testing.T) {
	time.Sleep(1 * time.Second)

	configurations := make(map[string]interface{})
	configurations["string"] = "string"
	configurations["int"] = "1"
	configurations["float"] = "1"
	configurations["bool"] = "true"
	configurations["slice"] = []int{1, 2}

	apolloConfig := &env.ApolloConfig{}
	apolloConfig.NamespaceName = defaultNamespace
	apolloConfig.AppID = "test"
	apolloConfig.Cluster = "dev"
	apolloConfig.Configurations = configurations
	UpdateApolloConfig(apolloConfig, true)

	currentConnApolloConfig := env.GetCurrentApolloConfig()
	config := currentConnApolloConfig[defaultNamespace]

	Assert(t, config, NotNilVal())
	Assert(t, defaultNamespace, Equal(config.NamespaceName))
	Assert(t, apolloConfig.AppID, Equal(config.AppID))
	Assert(t, apolloConfig.Cluster, Equal(config.Cluster))
	Assert(t, "", Equal(config.ReleaseKey))
	Assert(t, len(apolloConfig.Configurations), Equal(5))

}

func TestGetApolloConfigCache(t *testing.T) {
	cache := GetApolloConfigCache()
	Assert(t, cache, NotNilVal())
}

func TestGetDefaultNamespace(t *testing.T) {
	namespace := GetDefaultNamespace()
	Assert(t, namespace, Equal(defaultNamespace))
}

func TestGetConfig(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["string"] = "string2"
	configurations["int"] = "2"
	configurations["float"] = "1"
	configurations["bool"] = "false"
	configurations["sliceString"] = []string{"1", "2", "3"}
	configurations["sliceInt"] = []int{1, 2, 3}
	configurations["sliceInter"] = []interface{}{1, "2", 3}
	creatTestApolloConfig(configurations, "test")
	config := GetConfig("test")
	Assert(t, config, NotNilVal())

	//string
	s := config.GetStringValue("string", "s")
	Assert(t, s, Equal(configurations["string"]))

	s = config.GetStringValue("s", "s")
	Assert(t, s, Equal("s"))

	//int
	i := config.GetIntValue("int", 3)
	Assert(t, i, Equal(2))
	i = config.GetIntValue("i", 3)
	Assert(t, i, Equal(3))

	//float
	f := config.GetFloatValue("float", 2)
	Assert(t, f, Equal(float64(1)))
	f = config.GetFloatValue("f", 2)
	Assert(t, f, Equal(float64(2)))

	//bool
	b := config.GetBoolValue("bool", true)
	Assert(t, b, Equal(false))

	b = config.GetBoolValue("b", false)
	Assert(t, b, Equal(false))

	slice := config.GetStringSliceValue("sliceString")
	Assert(t, slice, Equal([]string{"1", "2", "3"}))

	sliceInt := config.GetIntSliceValue("sliceInt")
	Assert(t, sliceInt, Equal([]int{1, 2, 3}))

	sliceInter := config.GetSliceValue("sliceInter")
	Assert(t, sliceInter, Equal([]interface{}{1, "2", 3}))

	//content
	content := config.GetContent()
	hasFloat := strings.Contains(content, "float=1")
	Assert(t, hasFloat, Equal(true))

	hasInt := strings.Contains(content, "int=2")
	Assert(t, hasInt, Equal(true))

	hasString := strings.Contains(content, "string=string2")
	Assert(t, hasString, Equal(true))

	hasBool := strings.Contains(content, "bool=false")
	Assert(t, hasBool, Equal(true))

	hasSlice := strings.Contains(content, "sliceString=[1 2 3]")
	Assert(t, hasSlice, Equal(true))
	hasSlice = strings.Contains(content, "sliceInt=[1 2 3]")
	Assert(t, hasSlice, Equal(true))
}
