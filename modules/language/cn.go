// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package language

import "strings"

var cn = LangSet{
	"managers":  "管理员管理",
	"name":      "用户名",
	"nickname":  "昵称",
	"role":      "角色",
	"header":    "头部",
	"createdat": "创建时间",
	"updatedat": "更新时间",
	"path":      "路径",
	"submit":    "提交",
	"filter":    "筛选",
	"can like":  "支持模糊搜索",

	"new":             "新建",
	"export":          "导出",
	"action":          "操作",
	"toggle dropdown": "下拉",
	"delete":          "删除",
	"refresh":         "刷新",
	"back":            "返回",
	"reset":           "重置",
	"save":            "保存",
	"edit":            "编辑",
	"expand":          "展开",
	"collapse":        "折叠",
	"online":          "在线",
	"setting":         "设置",
	"sign out":        "登出",
	"home":            "首页",
	"all":             "全部",
	"more":            "更多",
	"browse":          "打开",
	"remove":          "移除",

	"permission manage": "权限管理",
	"menus manage":      "菜单管理",
	"roles manage":      "角色管理",
	"operation log":     "操作日志",
	"method":            "方法",
	"input":             "输入",
	"operation":         "操作",
	"menu name":         "菜单名",
	"reload succeeded":  "加载成功",
	"search":            "搜索",

	"permission denied": "没有权限",
	"error":             "错误",
	"success":           "成功",
	"fail":              "失败",
	"current page":      "当前页",

	"goadmin is now running. \nrunning in \"debug\" mode. switch to \"release\" mode in production.\n\n": "GoAdmin 启动成功。\n目前处于 \"debug\" 模式。请在生产环境中切换为 \"release\" 模式。\n\n",

	"wrong goadmin version, theme %s required goadmin version are %s":    "错误的 GoAdmin 版本，当前主题 %s 需要 GoAdmin 版本为 %s",
	"wrong theme version, goadmin %s required version of theme %s is %s": "错误的主题版本, GoAdmin %s 需要主题 %s 的版本为 %s",

	"adapter is nil, import the default adapter or use addadapter method add the adapter": "适配器为空，请先 import 对应的适配器或使用 AddAdapter 方法引入",

	"are you sure to delete": "你确定要删除吗？",
	"yes":                    "确定",
	"confirm":                "确认",
	"got it":                 "知道了",
	"cancel":                 "取消",
	"refresh succeeded":      "刷新成功",
	"delete succeed":         "删除成功",
	"edit fail":              "编辑失败",
	"create fail":            "新增失败",
	"delete fail":            "删除失败",
	"confirm password":       "确认密码",
	"all method if empty":    "为空默认为所有方法",

	"detail": "详情",

	"avatar":     "头像",
	"password":   "密码",
	"username":   "用户名",
	"slug":       "标志",
	"permission": "权限",
	"userid":     "用户ID",
	"content":    "内容",
	"parent":     "父级",
	"icon":       "图标",
	"uri":        "路径",
	"close":      "关闭",

	"login":      "登录",
	"login fail": "登录失败",

	"admin":     "管理",
	"user":      "用户",
	"users":     "用户",
	"roles":     "角色",
	"menu":      "菜单",
	"dashboard": "仪表盘",

	"continue editing":  "继续编辑",
	"continue creating": "继续新增",

	"username and password can not be empty":        "用户名密码不能为空",
	"operation not allow":                           "不允许的操作",
	"password does not match":                       "密码不一致",
	"should be unique":                              "需要保证唯一",
	"slug exists":                                   "标志已经存在了",
	"no corresponding options?":                     "没找到对应选项？",
	"create here.":                                  "在这里新建一个。",
	"use for login":                                 "用于登录",
	"use to display":                                "用来展示",
	"a path a line, without global prefix":          "一行一个路径，换行输入新路径，路径不包含全局路由前缀",
	"slug or http_path or name should not be empty": "标志或路径或权限名不能为空",
	"no roles":                                      "无角色",
	"no permission":                                 "没有权限",
	"fixed the sidebar":                             "固定侧边栏",
	"enter fullscreen":                              "进入全屏",
	"exit fullscreen":                               "退出全屏",
	"wrong captcha":                                 "错误的验证码",
	"modify success":                                "修改成功",

	"not found":      "找不到记录",
	"internal error": "系统内部错误",
	"unauthorized":   "未认证",

	"login overdue, please login again": "登录信息过期，请重新登录",
	"login info":                        "登录信息",

	"initialize configuration":        "初始化配置",
	"initialize navigation buttons":   "初始化导航栏按钮",
	"initialize plugins":              "初始化插件",
	"initialize database connections": "初始化数据库连接",
	"initialize success":              "初始化成功🍺🍺",

	"plugins":          "插件",
	"plugin store":     "插件商店",
	"get more plugins": "获取更多插件",
	"uninstalled":      "未安装",
	"plugin setting":   "插件设置",

	"showing <b>%s</b> to <b>%s</b> of <b>%s</b> entries": "显示第 <b>%s</b> 到第 <b>%s</b> 条记录，总共 <b>%s</b> 条记录",

	"second":  "秒",
	"seconds": "秒",
	"minute":  "分",
	"minutes": "分",
	"hour":    "小时",
	"hours":   "小时",
	"day":     "天",
	"days":    "天",
	"week":    "周",
	"weeks":   "周",
	"month":   "月",
	"months":  "月",
	"year":    "年",
	"years":   "年",

	"config.domain":          "网站域名",
	"config.language":        "网站语言",
	"config.url prefix":      "URL前缀",
	"config.theme":           "主题",
	"config.title":           "标题",
	"config.index url":       "首页URL",
	"config.login url":       "登录URL",
	"config.env":             "开发环境",
	"config.color scheme":    "颜色主题",
	"config.cdn url":         "cdn资源URL",
	"config.login title":     "登录标题",
	"config.auth user table": "登录用户表",
	"config.extra":           "额外配置",
	"config.store":           "文件存储设置",
	"config.databases":       "数据库设置",
	"config.general":         "通用",
	"config.log":             "日志",
	"config.site setting":    "网站设置",
	"config.custom":          "定制",
	"config.debug":           "Debug模式",
	"config.site off":        "关闭网站",
	"config.true":            "是",
	"config.false":           "否",

	"config.test":  "测试环境",
	"config.prod":  "生产环境",
	"config.local": "本地环境",

	"config.logo":                        "Logo",
	"config.mini logo":                   "Mini Logo",
	"config.session life time":           "Session时长",
	"config.bootstrap file path":         "插件文件路径",
	"config.go mod file path":            "go.mod文件路径",
	"config.custom head html":            "自定义Head HTML",
	"config.custom foot html":            "自定义Foot HTML",
	"config.custom 404 html":             "自定义404页面",
	"config.custom 403 html":             "自定义403页面",
	"config.custom 500 html":             "自定义500页面",
	"config.hide config center entrance": "隐藏配置中心入口",
	"config.hide app info entrance":      "隐藏应用信息入口",
	"config.hide tool entrance":          "隐藏工具入口",
	"config.hide plugin entrance":        "隐藏插件列表入口",
	"config.footer info":                 "自定义底部信息",
	"config.login logo":                  "登录Logo",
	"config.no limit login ip":           "取消限制多IP登录",
	"config.operation log off":           "关闭操作日志",
	"config.allow delete operation log":  "允许删除操作日志",
	"config.animation type":              "动画类型",
	"config.animation duration":          "动画间隔（秒）",
	"config.animation delay":             "动画延迟（秒）",
	"config.file upload engine":          "文件上传引擎",

	"config.logger rotate":             "日志切割设置",
	"config.logger rotate max size":    "存储最大文件大小（m）",
	"config.logger rotate max backups": "存储最多文件数",
	"config.logger rotate max age":     "最长存储时间（天）",
	"config.logger rotate compress":    "压缩",

	"config.info log path":         "信息日志存储路径",
	"config.error log path":        "错误日志存储路径",
	"config.access log path":       "访问日志存储路径",
	"config.info log off":          "关闭信息日志",
	"config.error log off":         "关闭错误日志",
	"config.access log off":        "关闭访问日志",
	"config.access assets log off": "关闭静态资源访问日志",
	"config.sql log on":            "打开SQL日志",
	"config.log level":             "日志级别",

	"config.logger rotate encoder":                "日志encoder设置",
	"config.logger rotate encoder time key":       "Time Key",
	"config.logger rotate encoder level key":      "Level Key",
	"config.logger rotate encoder name key":       "Name Key",
	"config.logger rotate encoder caller key":     "Caller Key",
	"config.logger rotate encoder message key":    "Message Key",
	"config.logger rotate encoder stacktrace key": "Stacktrace Key",
	"config.logger rotate encoder level":          "Level字段编码",
	"config.logger rotate encoder time":           "Time字段编码",
	"config.logger rotate encoder duration":       "Duration字段编码",
	"config.logger rotate encoder caller":         "Caller字段编码",
	"config.logger rotate encoder encoding":       "输出格式",

	"config.capital":        "大写",
	"config.capitalcolor":   "大写带颜色",
	"config.lowercase":      "小写",
	"config.lowercasecolor": "小写带颜色",

	"config.seconds":     "秒",
	"config.nanosecond":  "纳秒",
	"config.microsecond": "微秒",
	"config.millisecond": "毫秒",

	"config.full path":  "完整路径",
	"config.short path": "简短路径",

	"config.do not modify when you have not set up all assets": "不要修改，当你还没有设置好所有资源文件的时候",
	"config.it will work when theme is adminlte":               "当主题为adminlte时生效",
	"config.must bigger than 900 seconds":                      "必须大于900秒",

	"config.language." + CN:                  "中文",
	"config.language." + EN:                  "英文",
	"config.language." + JP:                  "日文",
	"config.language." + strings.ToLower(TC): "繁体中文",
	"config.language." + PTBR:                "Brazilian Portuguese",

	"config.modify site config":         "修改网站配置",
	"config.modify site config success": "修改网站配置成功",
	"config.modify site config fail":    "修改网站配置失败",

	"system.system info":     "应用系统信息",
	"system.application":     "应用信息",
	"system.application run": "应用运行信息",
	"system.system":          "系统信息",

	"system.process_id":                           "进程ID",
	"system.golang_version":                       "Golang版本",
	"system.server_uptime":                        "服务运行时间",
	"system.current_goroutine":                    "当前 Goroutines 数量",
	"system.current_memory_usage":                 "当前内存使用量",
	"system.total_memory_allocated":               "所有被分配的内存",
	"system.memory_obtained":                      "内存占用量",
	"system.pointer_lookup_times":                 "指针查找次数",
	"system.memory_allocate_times":                "内存分配次数",
	"system.memory_free_times":                    "内存释放次数",
	"system.current_heap_usage":                   "当前 Heap 内存使用量",
	"system.heap_memory_obtained":                 "Heap 内存占用量",
	"system.heap_memory_idle":                     "Heap 内存空闲量",
	"system.heap_memory_in_use":                   "正在使用的 Heap 内存",
	"system.heap_memory_released":                 "被释放的 Heap 内存",
	"system.heap_objects":                         "Heap 对象数量",
	"system.bootstrap_stack_usage":                "启动 Stack 使用量",
	"system.stack_memory_obtained":                "被分配的 Stack 内存",
	"system.mspan_structures_usage":               "MSpan 结构内存使用量",
	"system.mspan_structures_obtained":            "被分配的 MSpan 结构内存",
	"system.mcache_structures_usage":              "MCache 结构内存使用量",
	"system.mcache_structures_obtained":           "被分配的 MCache 结构内存",
	"system.profiling_bucket_hash_table_obtained": "被分配的剖析哈希表内存",
	"system.gc_metadata_obtained":                 "被分配的 GC 元数据内存",
	"system.other_system_allocation_obtained":     "其它被分配的系统内存",
	"system.next_gc_recycle":                      "下次 GC 内存回收量",
	"system.last_gc_time":                         "距离上次 GC 时间",
	"system.total_gc_time":                        "GC 执行时间总量",
	"system.total_gc_pause":                       "GC 暂停时间总量",
	"system.last_gc_pause":                        "上次 GC 暂停时间",
	"system.gc_times":                             "GC 执行次数",

	"system.cpu_logical_core": "cpu逻辑核数",
	"system.cpu_core":         "cpu物理核数",
	"system.os_platform":      "系统平台",
	"system.os_family":        "系统家族",
	"system.os_version":       "系统版本",
	"system.load1":            "1分钟内负载",
	"system.load5":            "5分钟内负载",
	"system.load15":           "15分钟内负载",
	"system.mem_total":        "总内存",
	"system.mem_available":    "可用内存",
	"system.mem_used":         "使用内存",

	"system.app_name":         "应用名",
	"system.go_admin_version": "应用版本",
	"system.theme_name":       "主题",
	"system.theme_version":    "主题版本",

	"tool.tool":                 "工具",
	"tool.table":                "表格",
	"tool.connection":           "连接",
	"tool.output path is empty": "输出路径为空",
	"tool.package":              "包名",
	"tool.output":               "输出路径",
	"tool.field":                "字段",
	"tool.title":                "标题",
	"tool.field name":           "字段名",
	"tool.db type":              "数据类型",
	"tool.form type":            "表单类型",
	"tool.generate table model": "生成CRUD模型",
	"tool.primarykey":           "主键",
	"tool.field filterable":     "可筛选",
	"tool.field sortable":       "可排序",
	"tool.yes":                  "是",
	"tool.no":                   "否",
	"tool.hide":                 "隐藏",
	"tool.show":                 "显示",
	"tool.generate success":     "生成成功",
	"tool.display":              "显示",
	"tool.use absolute path":    "使用绝对路径",
	"tool.basic info":           "基本信息",
	"tool.table info":           "表格信息",
	"tool.form info":            "表单信息",
	"tool.field editable":       "允许编辑",
	"tool.field can add":        "允许新增",
	"tool.info field editable":  "可编辑",
	"tool.field default":        "默认值",
	"tool.filter area":          "筛选框",
	"tool.new button":           "新建按钮",
	"tool.export button":        "导出按钮",
	"tool.edit button":          "编辑按钮",
	"tool.delete button":        "删除按钮",
	"tool.extra import package": "导入包",
	"tool.detail button":        "详情按钮",
	"tool.filter button":        "筛选按钮",
	"tool.row selector":         "列选择按钮",
	"tool.pagination":           "分页",
	"tool.query info":           "查询信息",
	"tool.filter form layout":   "筛选表单布局",
	"tool.generate":             "生成",
	"tool.generated tables":     "生成过的表格",
	"tool.description":          "描述",
	"tool.label":                "标签",
	"tool.image":                "图片",
	"tool.bool":                 "布尔",
	"tool.link":                 "链接",
	"tool.fileSize":             "文件大小",
	"tool.date":                 "日期",
	"tool.icon":                 "Icon",
	"tool.dot":                  "标点",
	"tool.progressBar":          "进度条",
	"tool.loading":              "Loading",
	"tool.downLoadable":         "可下载",
	"tool.copyable":             "可复制",
	"tool.carousel":             "图片轮播",
	"tool.qrcode":               "二维码",
	"tool.field hide":           "隐藏",
	"tool.field display":        "显示",
	"tool.table permission":     "生成表格权限",
	"tool.extra code":           "额外代码",

	"tool.detail display":             "显示",
	"tool.detail info":                "详情页信息",
	"tool.follow list page":           "跟随列表页",
	"tool.inherit from list page":     "继承列表页",
	"tool.independent from list page": "独立",

	"tool.continue edit checkbox": "继续编辑按钮",
	"tool.continue new checkbox":  "继续新增按钮",
	"tool.reset button":           "重设按钮",
	"tool.back button":            "返回按钮",

	"tool.field display normal":     "显示",
	"tool.field diplay hide":        "隐藏",
	"tool.field diplay edit hide":   "编辑隐藏",
	"tool.field diplay create hide": "新建隐藏",

	"tool.generate table model success": "生成成功",
	"tool.generate table model fail":    "生成失败",

	"generator.query":                 "查询",
	"generator.show edit form page":   "编辑页显示",
	"generator.show create form page": "新建记录页显示",
	"generator.edit":                  "编辑",
	"generator.create":                "新建",
	"generator.delete":                "删除",
	"generator.export":                "导出",

	"plugin.plugin":                         "插件",
	"plugin.plugin detail":                  "插件详情",
	"plugin.introduction":                   "介绍",
	"plugin.website":                        "网站",
	"plugin.version":                        "版本",
	"plugin.created at":                     "创建日期",
	"plugin.updated at":                     "更新日期",
	"plugin.provided by %s":                 "由 %s 提供",
	"plugin.upgrade":                        "升级",
	"plugin.install":                        "安装",
	"plugin.info":                           "详细信息",
	"plugin.download":                       "下载",
	"plugin.buy":                            "购买",
	"plugin.downloading":                    "下载中",
	"plugin.login":                          "登录",
	"plugin.login to goadmin member system": "登录到GoAdmin会员系统",
	"plugin.account":                        "账户名",
	"plugin.password":                       "密码",
	"plugin.learn more":                     "了解更多",

	"plugin.no account? click %s here %s to register.":    "没有账号？点击%s这里%s注册。",
	"plugin.download fail, wrong name":                    "下载失败，错误的参数",
	"plugin.change to debug mode first":                   "先切换到debug模式",
	"plugin.download fail, plugin not exist":              "下载失败，插件不存在",
	"plugin.download fail":                                "下载失败",
	"plugin.golang develop environment does not exist":    "golang开发环境不存在",
	"plugin.download success, restart to install":         "下载成功，重启程序进行安装",
	"plugin.restart to install":                           "重启程序进行安装",
	"plugin.can not connect to the goadmin remote server": "连接到GoAdmin远程服务器失败，请检查您的网络连接。",

	"admin.basic admin": "基础Admin",
	"admin.a built-in plugins of goadmin which help you to build a crud manager platform quickly.": "一个内置GoAdmin插件，帮助您快速搭建curd简易管理后台。",
	"admin.official": "GoAdmin官方",
}
