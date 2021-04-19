package middleware

import (
	"ginblog/utils"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricsPath = "/metrics"
	faviconPath = "/favicon.ico"
)

var application = "application"

var (
	// httpHistogram prometheus 模型
	httpHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "http_server",
		Subsystem:   "",
		Name:        "histogram",
		Help:        "Histogram of response latency (seconds) of http handlers.",
		ConstLabels: nil,
		//Buckets: []float64{},
		//默认分组时间{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
		Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30}, //自定义分组时间
	}, []string{"method", "code", "uri", application})

	httpCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "http_server",
		Subsystem:   "",
		Name:        "Counter",
		Help:        "Counter of response latency (seconds) of http handlers.",
		ConstLabels: nil,
	}, []string{"method", "code", "uri", application})

	httpSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "http_server",
		Subsystem: "",
		Name:      "Summary",
		Help:      "Summary of response latency (seconds) of http handlers.",
		// Objectives:  map[float64]float64{},
		// 默认的话就是sum
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //定义分位数
		ConstLabels: nil,
	}, []string{"method", "code", "uri", application})
)

// init 初始化prometheus模型
func init() {
	prometheus.MustRegister(httpHistogram, httpCounter, httpSummary)
}

// handlerPath 定义采样路由struct
type handlerPath struct {
	sync.Map
}

// get 获取path
func (hp *handlerPath) get(handler string) string {
	v, ok := hp.Load(handler)
	if !ok {
		return ""
	}
	return v.(string)
}

// set 保存path到sync.Map
func (hp *handlerPath) set(ri gin.RouteInfo) {
	hp.Store(ri.Handler, ri.Path)
}

// GinPrometheus gin调用Prometheus的struct
type GinPrometheus struct {
	engine      *gin.Engine
	ignored     map[string]bool
	pathMap     *handlerPath
	updated     bool
	application string
}

type Option func(*GinPrometheus)

// Ignore 添加忽略的路径
func Ignore(path ...string) Option {
	return func(gp *GinPrometheus) {
		for _, p := range path {
			gp.ignored[p] = true
		}
	}
}

// New new gin prometheus
func New(e *gin.Engine, options ...Option) *GinPrometheus {
	if e == nil {
		return nil
	}

	gp := &GinPrometheus{
		engine: e,
		ignored: map[string]bool{
			metricsPath: true,
			faviconPath: true,
		},
		pathMap:     &handlerPath{},
		application: utils.AppName,
	}

	for _, o := range options {
		o(gp)
	}
	return gp
}

// updatePath 更新path
func (gp *GinPrometheus) updatePath() {
	gp.updated = true
	for _, ri := range gp.engine.Routes() {
		gp.pathMap.set(ri)
	}
}

// Middleware set gin middleware
func (gp *GinPrometheus) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gp.updated {
			gp.updatePath()
		}
		// 过滤请求
		if gp.ignored[c.Request.URL.String()] {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		httpHistogram.WithLabelValues(
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
			gp.pathMap.get(c.HandlerName()),
			gp.application,
		).Observe(time.Since(start).Seconds())

		httpCounter.WithLabelValues(
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
			gp.pathMap.get(c.HandlerName()),
			gp.application,
		).Inc()

		httpSummary.WithLabelValues(
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
			gp.pathMap.get(c.HandlerName()),
			gp.application,
		).Observe(time.Since(start).Seconds())
	}
}
