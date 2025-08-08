// 文件功能: 定义 zrpc 客户端构建与使用入口，封装 gRPC 连接、鉴权、超时、Keepalive 等配置。
// 关键技术点:
// - 选项模式(Option Pattern): 通过 `ClientOption` 组合可选配置，灵活构建客户端。
// - gRPC 拦截器: 通过 `clientinterceptors` 实现日志、超时、熔断、耗时统计等横切能力。
// - 鉴权凭证: 使用 `PerRPCCredentials` 在每次 RPC 调用附带 App/Token。
// - 非阻塞拨号: 支持 `WithNonBlock` 在后台建立连接，提高启动速度。
// - Keepalive: 使用 `keepalive.ClientParameters` 维持长连接存活与探测。
// - 目标解析: 通过 `RpcClientConf.BuildTarget()` 支持直连/注册中心(如 etcd/discov) 等多种寻址方式。
// - 配置默认值: 使用 `conf.FillDefault` 快速生成默认配置实例。
// 适用场景: 快速创建具备治理能力（日志、熔断、超时等）的 gRPC 客户端。
// 包声明: 当前文件属于 `zrpc` 包，暴露对外的客户端 API。
package zrpc

// 导入依赖包列表
import (
	// 标准库 `time` 用于配置超时、阈值等基于时间的参数
	"time"

	// go-zero 核心配置与日志模块
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	// 内部客户端实现与鉴权模块
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/internal/auth"
	"github.com/zeromicro/go-zero/zrpc/internal/clientinterceptors"

	// gRPC 核心与 Keepalive 支持
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	// WithDialOption is an alias of internal.WithDialOption.
	// 技术点: 选项转发。对外导出与 `internal` 同名能力，保持 API 简洁统一。
	WithDialOption = internal.WithDialOption
	// WithNonBlock sets the dialing to be nonblock.
	// 技术点: 非阻塞连接建立，避免启动期被阻塞在网络连接。
	WithNonBlock = internal.WithNonBlock
	// WithStreamClientInterceptor is an alias of internal.WithStreamClientInterceptor.
	// 技术点: 流式拦截器，可统一处理 Stream RPC 的日志、指标等。
	WithStreamClientInterceptor = internal.WithStreamClientInterceptor
	// WithTimeout is an alias of internal.WithTimeout.
	// 技术点: 拨号/调用超时控制，提高健壮性。
	WithTimeout = internal.WithTimeout
	// WithTransportCredentials return a func to make the gRPC calls secured with given credentials.
	// 技术点: TLS/凭证配置，启用安全链路。
	WithTransportCredentials = internal.WithTransportCredentials
	// WithUnaryClientInterceptor is an alias of internal.WithUnaryClientInterceptor.
	// 技术点: 一元拦截器，可统一处理 Unary RPC 的日志、指标等。
	WithUnaryClientInterceptor = internal.WithUnaryClientInterceptor
)

type (
	// Client is an alias of internal.Client.
	// 对外暴露统一 Client 接口，具体实现位于 internal。
	Client = internal.Client
	// ClientOption is an alias of internal.ClientOption.
	// 对外暴露统一 Option 类型，兼容内部实现。
	ClientOption = internal.ClientOption

	// A RpcClient is a rpc client.
	// 封装内部 `Client`，对外提供更稳定的 API 表面。
	RpcClient struct {
		client Client
	}
)

// MustNewClient returns a Client, exits on any error.
func MustNewClient(c RpcClientConf, options ...ClientOption) Client {
	// 构建客户端，如出错则通过 logx.Must 触发 panic 以便立即暴露问题
	cli, err := NewClient(c, options...)
	logx.Must(err)
	return cli
}

// NewClient returns a Client.
func NewClient(c RpcClientConf, options ...ClientOption) (Client, error) {
	// 收集最终生效的选项，优先根据配置组装，再追加调用方传入的自定义选项
	var opts []ClientOption
	if c.HasCredential() {
		// 为每次 RPC 调用注入鉴权凭证(App/Token)
		opts = append(opts, WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
			App:   c.App,
			Token: c.Token,
		})))
	}
	if c.NonBlock {
		// 启用非阻塞拨号，提高启动速度
		opts = append(opts, WithNonBlock())
	}
	if c.Timeout > 0 {
		// 配置拨号超时，单位毫秒
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	if c.KeepaliveTime > 0 {
		// 配置客户端 keepalive，保持长连接活性
		opts = append(opts, WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: c.KeepaliveTime,
		})))
	}

	// 追加调用方自定义选项，后面的选项可以覆盖前面相同语义的设置
	opts = append(opts, options...)

	// 解析目标地址(直连/服务发现等)
	target, err := c.BuildTarget()
	if err != nil {
		return nil, err
	}

	// 根据目标与中间件配置创建内部客户端
	client, err := internal.NewClient(target, c.Middlewares, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

// NewClientWithTarget returns a Client with connecting to given target.
func NewClientWithTarget(target string, opts ...ClientOption) (Client, error) {
	// 生成带默认值的配置，便于快速通过目标地址创建客户端
	var config RpcClientConf
	if err := conf.FillDefault(&config); err != nil {
		return nil, err
	}

	config.Target = target

	return NewClient(config, opts...)
}

// Conn returns the underlying grpc.ClientConn.
func (rc *RpcClient) Conn() *grpc.ClientConn {
	// 暴露底层连接，便于外部直接构造 Stub 或进行底层操作
	return rc.client.Conn()
}

// DontLogClientContentForMethod disable logging content for given method.
func DontLogClientContentForMethod(method string) {
	// 关闭指定方法的请求/响应内容日志，避免敏感信息泄露或日志过大
	clientinterceptors.DontLogContentForMethod(method)
}

// SetClientSlowThreshold sets the slow threshold on client side.
func SetClientSlowThreshold(threshold time.Duration) {
	// 设置慢调用阈值，超过阈值的调用会以 slow 级别记录
	clientinterceptors.SetSlowThreshold(threshold)
}

// WithCallTimeout return a call option with given timeout to make a method call.
func WithCallTimeout(timeout time.Duration) grpc.CallOption {
	// 返回每次调用级别的超时设置，精细化控制单个 RPC 的调用超时
	return clientinterceptors.WithCallTimeout(timeout)
}
