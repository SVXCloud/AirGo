package initialize

import (
	"AirGo/api"
	"AirGo/global"
	"AirGo/middleware"
	"AirGo/web"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const apiPrefix = "/api"

//	type Resource struct {
//		fs   embed.FS
//		path string
//	}
//
//	func NewResource() *Resource {
//		return &Resource{
//			fs:   f,
//			path: "web",
//		}
//	}
//
//	func (r *Resource) Open(name string) (fs.File, error) {
//		//if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
//		//	return nil, errors.New("http: invalid character in file path")
//		//}
//		fullName := filepath.Join(r.path, filepath.FromSlash(path.Clean("/"+name)))
//		file, err := r.fs.Open(fullName)
//		return file, err
//	}
//
// ////go:embed all:web/*

// 初始化总路由
func InitRouter() {
	// 正式发布模式
	gin.SetMode(gin.ReleaseMode) //ReleaseMode TestMode DebugMode
	Router := gin.Default()

	//Router.Use(static.Serve("/", static.LocalFile("./web", true))) //静态资源(不嵌入)，可解决/问题。项目目录下web文件夹
	//Router.Static("/static", "static") //静态资源(不嵌入)
	//Router.StaticFS("/web", http.FS(NewResource())) //静态资源(嵌入)

	Router.Use(middleware.Serve("/", middleware.EmbedFolder(web.Static, "web"))) // targetPtah=web 是embed和web文件夹的相对路径

	Router.Use(middleware.Cors(), middleware.Recovery()) //不开启跨域验证码出错

	RouterGroup := Router.Group(apiPrefix)
	//public
	publicRouter := RouterGroup.Group("/public").Use(middleware.RateLimitIP())
	{
		publicRouter.POST("/getEmailCode", api.GetMailCode)         //获取验证码
		publicRouter.GET("/getBase64Captcha", api.GetBase64Captcha) //获取base64Captcha

	}

	//user
	userRouter := RouterGroup.Group("/user").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		userRouter.POST("/changeSubHost", api.ChangeSubHost)           //修改混淆
		userRouter.GET("/getUserInfo", api.GetUserInfo)                //获取自身信息
		userRouter.POST("/changeUserPassword", api.ChangeUserPassword) //修改密码
		userRouter.GET("/resetSub", api.ResetSub)                      //重置订阅
	}
	userAdminRouter := RouterGroup.Group("/user").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		userAdminRouter.POST("/getUserList", api.GetUserlist) //获取用户列表
		userAdminRouter.POST("/newUser", api.NewUser)         //新建用户
		userAdminRouter.POST("/updateUser", api.UpdateUser)   //修改用户
		userAdminRouter.POST("/deleteUser", api.DeleteUser)   //删除用户
	}
	userRouterNoVerify := RouterGroup.Group("/user").Use(middleware.RateLimitIP())
	{
		userRouterNoVerify.POST("/register", api.Register)                   //用户注册
		userRouterNoVerify.POST("/login", api.Login)                         //用户登录
		userRouterNoVerify.GET("/getSub", api.GetSub)                        //获取订阅
		userRouterNoVerify.POST("/resetUserPassword", api.ResetUserPassword) //重置密码
	}

	//菜单
	menuRouter := RouterGroup.Group("/menu").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		menuRouter.GET("/getRouteList", api.GetRouteList) //获取当前角色动态路由
	}
	menuAdminRouter := RouterGroup.Group("/menu").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		menuAdminRouter.GET("/getRouteTree", api.GetRouteTree)              //获取当前角色动态路由tree
		menuAdminRouter.GET("/getAllRouteList", api.GetAllRouteList)        //获取全部动态路由
		menuAdminRouter.GET("/getAllRouteTree", api.GetAllRouteTree)        //获取全部动态路由tree
		menuAdminRouter.POST("/newDynamicRoute", api.NewDynamicRoute)       //新建动态路由
		menuAdminRouter.POST("/delDynamicRoute", api.DelDynamicRoute)       //删除动态路由
		menuAdminRouter.POST("/updateDynamicRoute", api.UpdateDynamicRoute) //修改动态路由
		menuAdminRouter.POST("/findDynamicRoute", api.FindDynamicRoute)     //查询单条动态路由 by meta.title
	}

	//角色
	roleAdminRouter := RouterGroup.Group("/role").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		roleAdminRouter.POST("/getRoleList", api.GetRoleList)       //获取role list
		roleAdminRouter.POST("/modifyRoleInfo", api.ModifyRoleInfo) //更新role
		roleAdminRouter.POST("/addRole", api.AddRole)               //添加role
		roleAdminRouter.DELETE("/delRole", api.DelRole)             //删除role
	}

	//系统设置
	systemAdminRouter := RouterGroup.Group("/system").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		systemAdminRouter.POST("/updateThemeConfig", api.UpdateThemeConfig) //设置主题
		systemAdminRouter.GET("/getSetting", api.GetSetting)                //获取系统设置
		systemAdminRouter.POST("/updateSetting", api.UpdateSetting)         //修改系统设置
	}
	systemRouter := RouterGroup.Group("/system").Use(middleware.RateLimitIP())
	{
		systemRouter.GET("/getThemeConfig", api.GetThemeConfig)     //获取主题配置
		systemRouter.GET("/getPublicSetting", api.GetPublicSetting) //获取公共系统设置
	}

	//节点
	nodeAdminRouter := RouterGroup.Group("/node").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		nodeAdminRouter.GET("/getAllNode", api.GetAllNode)      //查询全部节点
		nodeAdminRouter.POST("/newNode", api.NewNode)           //新建节点
		nodeAdminRouter.POST("/deleteNode", api.DeleteNode)     //删除节点
		nodeAdminRouter.POST("/updateNode", api.UpdateNode)     //更新节点
		nodeAdminRouter.POST("/getTraffic", api.GetNodeTraffic) //获取节点 with Traffic,分页
		nodeAdminRouter.POST("/nodeSort", api.NodeSort)         //节点排序

		nodeAdminRouter.POST("/newNodeShared", api.NewNodeShared)        //新增节点
		nodeAdminRouter.GET("/getNodeSharedList", api.GetNodeSharedList) //获取节点列表
		nodeAdminRouter.POST("/deleteNodeShared", api.DeleteNodeShared)  //删除节点

	}

	//sspqnel
	sspanelRouter := RouterGroup.Group("/mod_mu")
	{
		sspanelRouter.GET("/nodes/:nodeID/info", api.SSNodeInfo) //获取节点信息
		sspanelRouter.GET("/users", api.SSUsers)                 //获取当前节点可连接的用户
		sspanelRouter.POST("/users/traffic", api.SSUsersTraffic) //上报用户的流量使用情况
		sspanelRouter.POST("/users/aliveip", api.SSUsersAliveIP) //上报用户的当前在线IP
	}

	//商店
	shopRouter := RouterGroup.Group("/shop").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		shopRouter.POST("/preCreatePay", api.PreCreateOrder) //alipay,统一收单线下交易预创建
		shopRouter.POST("/purchase", api.Purchase)           //支付
	}
	shopAdminRouter := RouterGroup.Group("/shop").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		shopAdminRouter.GET("/getAllEnabledGoods", api.GetAllEnabledGoods) // 查询全部已启用商品
		shopAdminRouter.GET("/getAllGoods", api.GetAllGoods)               //查询全部商品
		shopAdminRouter.POST("/newGoods", api.NewGoods)                    //新建商品
		shopAdminRouter.POST("/deleteGoods", api.DeleteGoods)              //删除商品
		shopAdminRouter.POST("/updateGoods", api.UpdateGoods)              //更新商品
		shopAdminRouter.POST("/goodsSort", api.GoodsSort)                  //排序
	}
	shopRouterNoVerify := RouterGroup.Group("/shop")
	{
		shopRouterNoVerify.POST("/alipayNotify", api.AlipayNotify) //异步验证支付结果
	}
	//订单
	orderRouter := RouterGroup.Group("/order").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		orderRouter.POST("/getOrderInfo", api.GetOrderInfo)         //获取订单详情(下单时）
		orderRouter.POST("/getOrderByUserID", api.GetOrderByUserID) //获取订单，分页获取
	}
	orderAdminRouter := RouterGroup.Group("/order").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		orderAdminRouter.POST("/getAllOrder", api.GetAllOrder)                         //获取全部订单，分页获取
		orderAdminRouter.POST("/completedOrder", api.CompletedOrder)                   //完成订单
		orderAdminRouter.POST("/getMonthOrderStatistics", api.GetMonthOrderStatistics) //获取时间范围内订单统计
	}
	//casbin
	casbinAdminRouter := RouterGroup.Group("/casbin").Use(middleware.ParseJwt(), middleware.Casbin())
	{
		casbinAdminRouter.GET("getAllPolicy", api.GetAllPolicy)                    //获取全部权限
		casbinAdminRouter.POST("getPolicyByRoleIds", api.GetPolicyByRoleIds)       //获取用户权限ByRoleIds
		casbinAdminRouter.POST("updateCasbinPolicy", api.UpdateCasbinPolicy)       //更新casbin权限
		casbinAdminRouter.POST("updateCasbinPolicyNew", api.UpdateCasbinPolicyNew) //更新casbin权限
	}
	//websocket
	websocketRouter := RouterGroup.Group("/websocket").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		websocketRouter.GET("msg", api.WebSocketMsg)
	}
	//upload

	uploadRouter := RouterGroup.Group("/upload").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		uploadRouter.GET("newPictureUrl", api.NewPictureUrl)
		uploadRouter.POST("getPictureList", api.GetPictureList)
	}
	//报表
	reportRouter := RouterGroup.Group("/report").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		reportRouter.GET("getDB", api.GetDB)
		reportRouter.POST("getTables", api.GetTables)
		reportRouter.POST("getColumn", api.GetColumnNew)
		reportRouter.POST("reportSubmit", api.ReportSubmit)

	}
	//文章
	articleRouter := RouterGroup.Group("/article").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		articleRouter.POST("newArticle", api.NewArticle)
		articleRouter.POST("deleteArticle", api.DeleteArticle)
		articleRouter.POST("updaterticle", api.UpdateArticle)
		articleRouter.POST("getArticle", api.GetArticle)
	}
	//折扣
	couponRouter := RouterGroup.Group("/coupon").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		couponRouter.POST("newCoupon", api.NewCoupon)
		couponRouter.POST("deleteCoupon", api.DeleteCoupon)
		couponRouter.POST("updateCoupon", api.UpdateCoupon)
		couponRouter.POST("getCoupon", api.GetCoupon)
	}
	//isp
	ispRouter := RouterGroup.Group("/isp").Use(middleware.RateLimitIP(), middleware.ParseJwt(), middleware.Casbin(), middleware.RateLimitVisit())
	{
		ispRouter.POST("sendCode", api.SendCode)
		ispRouter.POST("ispLogin", api.ISPLogin)
		//ispRouter.POST("queryPackage", api.QueryPackage) //
		ispRouter.POST("getMonitorByUserID", api.GetMonitorByUserID)

	}
	ispRouterNoVerify := RouterGroup.Group("/isp").Use(middleware.RateLimitIP())
	{
		ispRouterNoVerify.GET("queryPackage", api.QueryPackage) //
	}

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(global.Config.SystemParams.HTTPPort),
		Handler: Router,
	}
	srvTls := &http.Server{
		Addr:    ":" + strconv.Itoa(global.Config.SystemParams.HTTPSPort),
		Handler: Router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logrus.Fatalf("listen: %s\n", err)
		}
	}()
	go func() {
		// 服务连接
		if err := srvTls.ListenAndServeTLS("./air.cer", "./air.key"); err != nil && err != http.ErrServerClosed {
			global.Logrus.Error("tls listen: %s\n", err)
		}
	}()

	// 等待中断信号关闭服务器(设置 5 秒的超时时间)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	global.Logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Logrus.Fatalf("Server Shutdown:", err)
	}
	if err := srvTls.Shutdown(ctx); err != nil {
		global.Logrus.Fatalf("Server Shutdown:", err)
	}
	global.Logrus.Info("Server exiting")

}
