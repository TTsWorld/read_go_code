// 文件功能: 定义 zrpc 服务端的构建与启动流程，包装内部 gRPC Server 实现与各类中间件（拦截器）。
// 关键技术点:
// - 服务治理: 统一设置链路追踪、恢复、熔断、统计、Prometheus、限流/负载剪枝、超时等拦截器。
// - 选项模式(Option Pattern): 通过 `ServerOption`、`Add*Interceptors` 等方式动态扩展能力。
// - 健康检查: 通过 `WithRpcHealth` 启用健康探针。
// - 多注册方式: etcd 发布/直连两种 Server 形态，兼容服务发现或直连部署。
// - 认证鉴权: 集成 Redis 支持的 Token/JWT 鉴权中间件。
// - 指标统计: 使用 `stat.Metrics` 统一指标收集与命名。
// 适用场景: 快速创建具备完善治理能力的 gRPC 服务端。
// 包声明: 当前文件属于 `zrpc` 包，暴露对外的服务端 API。
package zrpc

// 导入依赖包列表
import (
	// 标准库时间处理
	"time"

	// go-zero 核心: 负载剪枝、日志、指标、Redis 存储
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/redis"

	// 内部 server 抽象与鉴权、服务端拦截器
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/internal/auth"
	"github.com/zeromicro/go-zero/zrpc/internal/serverinterceptors"

	// gRPC 核心库
	"google.golang.org/grpc"
)

// A RpcServer is a rpc server.
type RpcServer struct {
	// `server` 为内部 Server 实现，可能是直连或发布型（含注册中心）
	server internal.Server
	// `register` 为用户传入的服务注册函数，用于将 gRPC 服务实现注册到 server
	register internal.RegisterFn
}

// MustNewServer returns a RpcSever, exits on any error.
func MustNewServer(c RpcServerConf, register internal.RegisterFn) *RpcServer {
	// 构建服务端，如出错则立即 panic 暴露问题
	server, err := NewServer(c, register)
	logx.Must(err)
	return server
}

// NewServer returns a RpcServer.
func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
	// 配置校验，确保关键参数合法
	var err error
	if err = c.Validate(); err != nil {
		return nil, err
	}

	// 根据是否配置 etcd 决定创建发布型或直连型 server
	var server internal.Server
	metrics := stat.NewMetrics(c.ListenOn)
	serverOptions := []internal.ServerOption{
		internal.WithRpcHealth(c.Health),
	}

	if c.HasEtcd() {
		// etcd 发布模式，供客户端通过注册中心发现
		server, err = internal.NewRpcPubServer(c.Etcd, c.ListenOn, serverOptions...)
		if err != nil {
			return nil, err
		}
	} else {
		// 直连模式，仅监听地址提供服务
		server = internal.NewRpcServer(c.ListenOn, serverOptions...)
	}

	// 统一设置服务名与指标名，便于观测与排障
	server.SetName(c.Name)
	metrics.SetName(c.Name)
	// 安装流式拦截器
	setupStreamInterceptors(server, c)
	// 安装一元拦截器（含统计/Prometheus/熔断/超时/自适应限流等）
	setupUnaryInterceptors(server, c, metrics)
	// 安装鉴权拦截器（可选）
	if err = setupAuthInterceptors(server, c); err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		server:   server,
		register: register,
	}
	// 执行额外自定义初始化逻辑（如日志、链路、认证等外部依赖）
	if err = c.SetUp(); err != nil {
		return nil, err
	}

	return rpcServer, nil
}

// AddOptions adds given options.
func (rs *RpcServer) AddOptions(options ...grpc.ServerOption) {
	// 透传 gRPC 原生 ServerOption，满足高级定制需求
	rs.server.AddOptions(options...)
}

// AddStreamInterceptors adds given stream interceptors.
func (rs *RpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	// 追加流式拦截器，按添加顺序生效
	rs.server.AddStreamInterceptors(interceptors...)
}

// AddUnaryInterceptors adds given unary interceptors.
func (rs *RpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	// 追加一元拦截器，按添加顺序生效
	rs.server.AddUnaryInterceptors(interceptors...)
}

// Start starts the RpcServer.
// Graceful shutdown is enabled by default.
// Use proc.SetTimeToForceQuit to customize the graceful shutdown period.
func (rs *RpcServer) Start() {
	// 启动服务并注册业务实现，内部已启用优雅退出
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}

// Stop stops the RpcServer.
func (rs *RpcServer) Stop() {
	// 关闭日志与资源，配合优雅退出
	logx.Close()
}

// DontLogContentForMethod disable logging content for given method.
// Deprecated: use ServerMiddlewaresConf.IgnoreContentMethods instead.
func DontLogContentForMethod(method string) {
	// 禁止对指定方法记录请求/响应内容（建议使用配置项替代）
	serverinterceptors.DontLogContentForMethod(method)
}

// SetServerSlowThreshold sets the slow threshold on server side.
// Deprecated: use ServerMiddlewaresConf.SlowThreshold instead.
func SetServerSlowThreshold(threshold time.Duration) {
	// 设置慢调用阈值（建议使用配置项替代）
	serverinterceptors.SetSlowThreshold(threshold)
}

func setupAuthInterceptors(svr internal.Server, c RpcServerConf) error {
	// 未开启鉴权直接返回
	if !c.Auth {
		return nil
	}
	// 构建 Redis 客户端，鉴权需要依赖存储
	rds, err := redis.NewRedis(c.Redis.RedisConf)
	if err != nil {
		return err
	}

	// 基于 Redis 与密钥、严格模式创建鉴权器
	authenticator, err := auth.NewAuthenticator(rds, c.Redis.Key, c.StrictControl)
	if err != nil {
		return err
	}

	// 安装鉴权拦截器（流式与一元）
	svr.AddStreamInterceptors(serverinterceptors.StreamAuthorizeInterceptor(authenticator))
	svr.AddUnaryInterceptors(serverinterceptors.UnaryAuthorizeInterceptor(authenticator))

	return nil
}

func setupStreamInterceptors(svr internal.Server, c RpcServerConf) {
	// 链路追踪（可选）
	if c.Middlewares.Trace {
		svr.AddStreamInterceptors(serverinterceptors.StreamTracingInterceptor)
	}
	// panic 恢复（可选）
	if c.Middlewares.Recover {
		svr.AddStreamInterceptors(serverinterceptors.StreamRecoverInterceptor)
	}
	// 熔断（可选）
	if c.Middlewares.Breaker {
		svr.AddStreamInterceptors(serverinterceptors.StreamBreakerInterceptor)
	}
}

func setupUnaryInterceptors(svr internal.Server, c RpcServerConf, metrics *stat.Metrics) {
	// 链路追踪（可选）
	if c.Middlewares.Trace {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryTracingInterceptor)
	}
	// panic 恢复（可选）
	if c.Middlewares.Recover {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryRecoverInterceptor)
	}
	// 指标统计（可选）
	if c.Middlewares.Stat {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryStatInterceptor(metrics, c.Middlewares.StatConf))
	}
	// Prometheus 指标（可选）
	if c.Middlewares.Prometheus {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryPrometheusInterceptor)
	}
	// 熔断（可选）
	if c.Middlewares.Breaker {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryBreakerInterceptor)
	}
	// 自适应负载剪枝（可选），基于 CPU 阈值
	if c.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
		svr.AddUnaryInterceptors(serverinterceptors.UnarySheddingInterceptor(shedder, metrics))
	}
	// 调用超时（可选），支持方法级覆盖
	if c.Timeout > 0 {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout)*time.Millisecond, c.MethodTimeouts...))
	}
}
