package middleware

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/model"
	"github.com/TencentBlueKing/bk-apigateway-sdks/gin_fram/util"
)

// middlewareConfigs 用于存储中间件配置
var middlewareConfigs = make(map[uintptr]model.APIGatewayResourceConfig)
var funcCounter uint64

// WithGatewayResourceConfig 网关route配置中间件
func WithGatewayResourceConfig(config model.APIGatewayResourceConfig) gin.HandlerFunc {
	// 获取闭包函数的唯一标识
	wrapper := genHandler()
	funcPtr := reflect.ValueOf(wrapper).Pointer()
	// 注册中间件配置
	middlewareConfigs[funcPtr] = config
	return wrapper
}

// GetRouteConfigMap 获取路由网关配置
func GetRouteConfigMap(engine *gin.Engine) map[string]*model.APIGatewayResourceConfig {
	var configs []*model.RouteConfig
	engineVal := reflect.ValueOf(engine).Elem()

	// 获取方法树列表（注意methodTrees的类型定义）
	treesField := engineVal.FieldByName("trees")
	if !treesField.IsValid() {
		log.Println("路由树字段验证失败，请确认Gin版本兼容性")
		return nil
	}
	// 遍历所有HTTP方法树
	for i := 0; i < treesField.Len(); i++ {
		tree := treesField.Index(i) // methodTree是结构体非指针
		methodField := tree.FieldByName("method")
		rootField := tree.FieldByName("root")

		// 跳过无效节点
		if !methodField.IsValid() || rootField.IsNil() {
			continue
		}

		method := methodField.String()
		root := rootField.Elem() // 解引用node指针

		// 使用队列进行非递归遍历
		queue := []struct {
			node reflect.Value
			path string
		}{{node: root, path: ""}}

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			// 获取节点路径信息
			path := current.node.FieldByName("path").String()
			fullPath := current.node.FieldByName("fullPath").String()
			combinedPath := current.path + path

			// 优先使用fullPath（如果有的话）
			finalPath := fullPath
			if fullPath == "" {
				finalPath = combinedPath
			}

			// 提取处理链
			handlers := current.node.FieldByName("handlers")
			if handlers.IsValid() && handlers.Len() > 0 {
				for j := 0; j < handlers.Len(); j++ {
					if cfg, ok := extractHandlerConfig(handlers.Index(j)); ok {
						configs = append(configs, &model.RouteConfig{
							Method: method,
							Path:   finalPath,
							Config: cfg,
						})
					}
				}
			}

			// 处理子节点
			children := current.node.FieldByName("children")
			for j := 0; j < children.Len(); j++ {
				child := children.Index(j)
				if child.Kind() == reflect.Ptr && !child.IsNil() {
					queue = append(queue, struct {
						node reflect.Value
						path string
					}{node: child.Elem(), path: combinedPath})
				}
			}
		}
	}

	routeMap := make(map[string]*model.APIGatewayResourceConfig)
	for _, rc := range configs {
		key := fmt.Sprintf("%s:%s", util.ConvertExpressPathToSwagger(rc.Path), strings.ToLower(rc.Method))
		routeMap[key] = rc.Config
	}

	return routeMap
}

// extractHandlerConfig 提取处理链中的中间件配置
func extractHandlerConfig(handlerVal reflect.Value) (*model.APIGatewayResourceConfig, bool) {
	// 验证函数指针有效性
	if handlerVal.Kind() != reflect.Func || handlerVal.IsNil() {
		return nil, false
	}
	// 通过指针地址匹配注册配置
	funcPtr := handlerVal.Pointer()
	if cfg, exists := middlewareConfigs[funcPtr]; exists {
		return &cfg, true
	}
	return nil, false
}

//go:noinline
func genHandler() gin.HandlerFunc {
	id := generateUUID()
	counter := atomic.AddUint64(&funcCounter, 8)
	handler := func(c *gin.Context) {
		_ = counter // 阻止编译器优化
		c.Set("id", id)
		c.Next()
		runtime.KeepAlive(counter)
	}
	runtime.KeepAlive(counter)
	ptr := (*uintptr)(unsafe.Pointer(&handler))
	*ptr += uintptr(counter) // 人为修改指针，防止优化
	return handler
}

func generateUUID() string {
	return fmt.Sprintf("mw-%d-%d-%d",
		time.Now().UnixNano(),
		os.Getpid(),
		rand.Intn(1000),
	)
}
