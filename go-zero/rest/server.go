// 文件功能: 定义 REST 服务器的构建、路由注册、启动与中间件装配，支持 CORS、静态文件、JWT、签名校验、SSE、超时等特性。
// 关键技术点:
// - 组合模式: `Server` 内部持有 `engine` 与 `router`，职责分离。
// - 选项模式(RunOption/RouteOption): 通过函数式选项定制运行期行为与路由特性。
// - 可插拔中间件链: 使用 `chain.Chain` 与 `Use` 方法灵活扩展。
// - CORS/静态文件包装: 通过自定义 `httpx.Router` 包装器透明增强。
// - 优雅退出: 默认启用，异常使用统一的 `handleError` 兜底处理。
// - TLS 支持: 通过 `WithTLSConfig` 注入自定义 TLS 配置。
// - 非侵入 API: 大量 `With*` 助手函数降低接入成本。
// 适用场景: 快速落地具备治理能力的 HTTP 服务。
// 包声明: 当前文件属于 `rest` 包，暴露对外的 REST Server API。
package rest

// 导入依赖包列表
import (
	// TLS、错误类型、HTTP、路径处理、时间
	"crypto/tls"
	"errors"
	"net/http"
	"path"
	"time"

	// go-zero 核心与子系统能力
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/chain"
	"github.com/zeromicro/go-zero/rest/handler"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/internal"
	"github.com/zeromicro/go-zero/rest/internal/cors"
	"github.com/zeromicro/go-zero/rest/internal/fileserver"
	"github.com/zeromicro/go-zero/rest/router"
)

type (
	// RunOption defines the method to customize a Server.
	// 运行期选项：用于调整 `Server` 的整体行为
	RunOption func(*Server)

	// StartOption defines the method to customize http server.
	// 启动选项：用于调整底层 HTTP 服务器行为（端口、超时、优雅退出等）
	StartOption = internal.StartOption

	// A Server is a http server.
	// 对外的 REST 服务器抽象，持有引擎与路由器
	Server struct {
		ngin   *engine
		router httpx.Router
	}
)

// MustNewServer returns a server with given config of c and options defined in opts.
// Be aware that later RunOption might overwrite previous one that write the same option.
// The process will exit if error occurs.
func MustNewServer(c RestConf, opts ...RunOption) *Server {
	// 构建服务端，如出错则立即 panic 暴露问题
	server, err := NewServer(c, opts...)
	if err != nil {
		logx.Must(err)
	}

	return server
}

// NewServer returns a server with given config of c and options defined in opts.
// Be aware that later RunOption might overwrite previous one that write the same option.
func NewServer(c RestConf, opts ...RunOption) (*Server, error) {
	// 配置落地（如日志、验证器、限流器等初始化）
	if err := c.SetUp(); err != nil {
		return nil, err
	}

	server := &Server{
		ngin:   newEngine(c),
		router: router.NewRouter(),
	}

	// 默认安装 404 处理器，调用方可覆盖
	opts = append([]RunOption{WithNotFoundHandler(nil)}, opts...)
	for _, opt := range opts {
		// 应用 RunOption，后定义的选项可覆盖先前设置
		opt(server)
	}

	return server, nil
}

// AddRoute adds given route into the Server.
func (s *Server) AddRoute(r Route, opts ...RouteOption) {
	// 复用批量添加函数，统一行为
	s.AddRoutes([]Route{r}, opts...)
}

// AddRoutes add given routes into the Server.
func (s *Server) AddRoutes(rs []Route, opts ...RouteOption) {
	r := featuredRoutes{
		routes: rs,
	}
	for _, opt := range opts {
		opt(&r)
	}
	// 将带特性的路由交给引擎统一管理
	s.ngin.addRoutes(r)
}

// PrintRoutes prints the added routes to stdout.
func (s *Server) PrintRoutes() {
	s.ngin.print()
}

// Routes returns the HTTP routers that registered in the server.
func (s *Server) Routes() []Route {
	// 从引擎内部聚合出所有已注册的原始路由
	routes := make([]Route, 0, len(s.ngin.routes))

	for _, r := range s.ngin.routes {
		routes = append(routes, r.routes...)
	}

	return routes
}

// Start starts the Server.
// Graceful shutdown is enabled by default.
// Use proc.SetTimeToForceQuit to customize the graceful shutdown period.
func (s *Server) Start() {
	// 启动引擎并绑定当前路由器，统一错误处理
	handleError(s.ngin.start(s.router))
}

// StartWithOpts starts the Server.
// Graceful shutdown is enabled by default.
// Use proc.SetTimeToForceQuit to customize the graceful shutdown period.
func (s *Server) StartWithOpts(opts ...StartOption) {
	// 允许按需传入底层 HTTP 服务器的启动选项
	handleError(s.ngin.start(s.router, opts...))
}

// Stop stops the Server.
func (s *Server) Stop() {
	// 关闭日志与资源，配合优雅退出
	logx.Close()
}

// Use adds the given middleware in the Server.
func (s *Server) Use(middleware Middleware) {
	// 将中间件追加到引擎链路中
	s.ngin.use(middleware)
}

// build builds the Server and binds the routes to the router.
func (s *Server) build() error {
	// 将引擎中的路由与中间件绑定到路由器
	return s.ngin.bindRoutes(s.router)
}

// serve serves the HTTP requests using the Server's router.
func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	// 直接交给路由器处理
	s.router.ServeHTTP(w, r)
}

// ToMiddleware converts the given handler to a Middleware.
func ToMiddleware(handler func(next http.Handler) http.Handler) Middleware {
	// 适配器: 将 `func(http.Handler) http.Handler` 转换为框架统一的 `Middleware` 类型
	return func(handle http.HandlerFunc) http.HandlerFunc {
		return handler(handle).ServeHTTP
	}
}

// WithChain returns a RunOption that uses the given chain to replace the default chain.
// JWT auth middleware and the middlewares that added by svr.Use() will be appended.
func WithChain(chn chain.Chain) RunOption {
	return func(svr *Server) {
		svr.ngin.chain = chn
	}
}

// WithCors returns a func to enable CORS for given origin, or default to all origins (*).
func WithCors(origin ...string) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(cors.NotAllowedHandler(nil, origin...))
		server.router = newCorsRouter(server.router, nil, origin...)
	}
}

// WithCorsHeaders returns a RunOption to enable CORS with given headers.
func WithCorsHeaders(headers ...string) RunOption {
	const allDomains = "*"

	return func(server *Server) {
		server.router.SetNotAllowedHandler(cors.NotAllowedHandler(nil, allDomains))
		server.router = newCorsRouter(server.router, func(header http.Header) {
			cors.AddAllowHeaders(header, headers...)
		}, allDomains)
	}
}

// WithCustomCors returns a func to enable CORS for given origin, or default to all origins (*),
// fn lets caller customizing the response.
func WithCustomCors(middlewareFn func(header http.Header), notAllowedFn func(http.ResponseWriter),
	origin ...string) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(cors.NotAllowedHandler(notAllowedFn, origin...))
		server.router = newCorsRouter(server.router, middlewareFn, origin...)
	}
}

// WithFileServer returns a RunOption to serve files from given dir with given path.
func WithFileServer(path string, fs http.FileSystem) RunOption {
	return func(server *Server) {
		server.router = newFileServingRouter(server.router, path, fs)
	}
}

// WithJwt returns a func to enable jwt authentication in given route.
func WithJwt(secret string) RouteOption {
	return func(r *featuredRoutes) {
		validateSecret(secret)
		r.jwt.enabled = true
		r.jwt.secret = secret
	}
}

// WithJwtTransition returns a func to enable jwt authentication as well as jwt secret transition.
// Which means old and new jwt secrets work together for a period.
func WithJwtTransition(secret, prevSecret string) RouteOption {
	return func(r *featuredRoutes) {
		// why not validate prevSecret, because prevSecret is an already used one,
		// even it not meet our requirement, we still need to allow the transition.
		validateSecret(secret)
		r.jwt.enabled = true
		r.jwt.secret = secret
		r.jwt.prevSecret = prevSecret
	}
}

// WithMaxBytes returns a RouteOption to set maxBytes with the given value.
func WithMaxBytes(maxBytes int64) RouteOption {
	return func(r *featuredRoutes) {
		r.maxBytes = maxBytes
	}
}

// WithMiddlewares adds given middlewares to given routes.
func WithMiddlewares(ms []Middleware, rs ...Route) []Route {
	for i := len(ms) - 1; i >= 0; i-- {
		rs = WithMiddleware(ms[i], rs...)
	}
	return rs
}

// WithMiddleware adds given middleware to given route.
func WithMiddleware(middleware Middleware, rs ...Route) []Route {
	routes := make([]Route, len(rs))

	for i := range rs {
		route := rs[i]
		routes[i] = Route{
			Method:  route.Method,
			Path:    route.Path,
			Handler: middleware(route.Handler),
		}
	}

	return routes
}

// WithNotFoundHandler returns a RunOption with not found handler set to given handler.
func WithNotFoundHandler(handler http.Handler) RunOption {
	return func(server *Server) {
		notFoundHandler := server.ngin.notFoundHandler(handler)
		server.router.SetNotFoundHandler(notFoundHandler)
	}
}

// WithNotAllowedHandler returns a RunOption with not allowed handler set to given handler.
func WithNotAllowedHandler(handler http.Handler) RunOption {
	return func(server *Server) {
		server.router.SetNotAllowedHandler(handler)
	}
}

// WithPrefix adds group as a prefix to the route paths.
func WithPrefix(group string) RouteOption {
	return func(r *featuredRoutes) {
		routes := make([]Route, 0, len(r.routes))
		for _, rt := range r.routes {
			p := path.Join(group, rt.Path)
			routes = append(routes, Route{
				Method:  rt.Method,
				Path:    p,
				Handler: rt.Handler,
			})
		}
		r.routes = routes
	}
}

// WithPriority returns a RunOption with priority.
func WithPriority() RouteOption {
	return func(r *featuredRoutes) {
		r.priority = true
	}
}

// WithRouter returns a RunOption that make server run with given router.
func WithRouter(router httpx.Router) RunOption {
	return func(server *Server) {
		server.router = router
	}
}

// WithSignature returns a RouteOption to enable signature verification.
func WithSignature(signature SignatureConf) RouteOption {
	return func(r *featuredRoutes) {
		r.signature.enabled = true
		r.signature.Strict = signature.Strict
		r.signature.Expiry = signature.Expiry
		r.signature.PrivateKeys = signature.PrivateKeys
	}
}

// WithSSE returns a RouteOption to enable server-sent events.
func WithSSE() RouteOption {
	return func(r *featuredRoutes) {
		r.sse = true
	}
}

// WithTimeout returns a RouteOption to set timeout with given value.
func WithTimeout(timeout time.Duration) RouteOption {
	return func(r *featuredRoutes) {
		r.timeout = &timeout
	}
}

// WithTLSConfig returns a RunOption that with given tls config.
func WithTLSConfig(cfg *tls.Config) RunOption {
	return func(svr *Server) {
		svr.ngin.setTlsConfig(cfg)
	}
}

// WithUnauthorizedCallback returns a RunOption that with given unauthorized callback set.
func WithUnauthorizedCallback(callback handler.UnauthorizedCallback) RunOption {
	return func(svr *Server) {
		svr.ngin.setUnauthorizedCallback(callback)
	}
}

// WithUnsignedCallback returns a RunOption that with given unsigned callback set.
func WithUnsignedCallback(callback handler.UnsignedCallback) RunOption {
	return func(svr *Server) {
		svr.ngin.setUnsignedCallback(callback)
	}
}

func handleError(err error) {
	// ErrServerClosed means the server is closed manually
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	logx.Error(err)
	panic(err)
}

func validateSecret(secret string) {
	if len(secret) < 8 {
		panic("secret's length can't be less than 8")
	}
}

type corsRouter struct {
	httpx.Router
	middleware Middleware
}

func newCorsRouter(router httpx.Router, headerFn func(http.Header), origins ...string) httpx.Router {
	return &corsRouter{
		Router:     router,
		middleware: cors.Middleware(headerFn, origins...),
	}
}

func (c *corsRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.middleware(c.Router.ServeHTTP)(w, r)
}

type fileServingRouter struct {
	httpx.Router
	middleware Middleware
}

func newFileServingRouter(router httpx.Router, path string, fs http.FileSystem) httpx.Router {
	return &fileServingRouter{
		Router:     router,
		middleware: fileserver.Middleware(path, fs),
	}
}

func (f *fileServingRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.middleware(f.Router.ServeHTTP)(w, r)
}
