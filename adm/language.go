package main

var langs = map[string]map[string]string{
    "cn": {
        "user login name": "用户登录名",
        "user nickname":   "用户昵称",
        "user password":   "用户密码",

        "choose a driver":          "选择数据库驱动",
        "sql address":              "连接地址",
        "sql port":                 "端口",
        "sql username":             "用户名",
        "sql schema":               "Schema",
        "sql database name":        "数据库名",
        "sql file":                 "数据库文件地址",
        "sql password":             "密码",
        "choose table to generate": "选择要生成的表格",

        "wrong config file path": "错误的配置文件路径",
        "user record exists":     "用户记录已存在",
        "empty tables":           "表格不能为空",

        "tables to generate, use comma to split": "要生成权限的表格，用逗号分隔",

        "no tables, you should build a table of your own business first.": "表格不能为空，请先创建您的业务表",
        "no table is selected": "没有选择表格",

        "set package name":     "设置包名",
        "set connection name":  "设置连接",
        "set file output path": "设置文件输出路径",

        "generate permission records for tables": "是否生成表格权限",

        "Query":                 "查询",
        "Show Edit Form Page":   "编辑页显示",
        "Show Create Form Page": "新建记录页显示",
        "Edit":                  "编辑",
        "Create":                "新建",
        "Delete":                "删除",
        "Export":                "导出",

        "Use arrows to move, type to filter, enter to select": "使用方向键去移动，空格键选择，输入进行筛选",
        "select all": "选择全部",
        "Use arrows to move, space to select, type to filter": "使用方向键去移动，空格键选择，输入进行筛选",
        "Add admin user success~~":                            "增加用户成功~~",
        "Add table permissions success~~":                     "增加表格权限成功~~",
        "Generate data table models success~~":                "生成数据模型文件成功~~",
        "see the docs: ":                                      "查看文档：",
        "visit forum: ":                                       "访问论坛：",
        "generating: ":                                        "生成中：",

        "choose a theme":                          "选择主题",
        "choose language":                         "选择语言",
        "choose framework":                        "选择框架",
        "choose a orm":                            "选择一个ORM",
        "none":                                    "不使用",
        "Generate project template success~~":     "生成项目模板成功~~",
        "1 Import and initialize database:":       "1 安装初始化数据库：",
        "2 Execute the following command to run:": "2 执行以下命令运行：",
        "3 Visit and login:":                      "3 访问并登陆：",
        "4 See more in README.md":                 "4 在README.md中查看更多",
        "account: admin  password: admin":         "账号：admin，密码：admin",
        "Login: ":                                 "登陆：",
        "Generate CRUD models: ":                  "生成CRUD模型：",

        "GoAdmin CLI error: CLI has not supported MINGW64 for now, please use cmd terminal instead.": "GoAdmin CLI" +
            "错误：目前不支持 MINGW64，请使用 CMD 终端。",
        "Know more: http://discuss.go-admin.com/t/goadmin-cli-adm-does-not-support-git-bash-mingw64-for-now/77": "了解更多：" +
            "http://discuss.go-admin.com/t/goadmin-cli-adm-git-bash-mingw64/22",

        "port":        "端口",
        "url prefix":  "路由前缀",
        "module path": "模块路径",

        "yes": "是",
        "no":  "否",

        "cn": "简体中文",
        "en": "英文",
        "jp": "日文",
        "tc": "繁体中文",

        "web.database settings":     "数据库设置",
        "web.installation settings": "应用安装设置",
        "web.application settings":  "应用设置",
        "web.optional settings":     "可选设置",

        "web.database type":     "数据库类型",
        "web.database host":     "数据库主机",
        "web.database user":     "数据库用户",
        "web.database password": "数据库用户密码",
        "web.database name":     "数据库名称",
        "web.database file":     "数据库文件路径",
        "web.database port":     "数据库端口",
        "web.database schema":   "数据库Schema",

        "web.theme":             "主题",
        "web.language":          "语言",
        "web.web framework":     "Web框架",
        "web.module name":       "模块名",
        "web.http port":         "HTTP 端口号",
        "web.url prefix":        "路由前缀",
        "web.website title":     "网站标题",
        "web.login page logo":   "登录页面logo",
        "web.sidebar logo":      "侧边栏logo",
        "web.sidebar mini logo": "侧边栏缩小logo",
        "web.use orm":           "使用 ORM",
        "web.no use":            "不使用",
        "web.input":             "输入",

        "web.simplified chinese":  "简体中文",
        "web.traditional chinese": "繁体中文",
        "web.english":             "英文",
        "web.japanese":            "日文",

        "web.where the framework sql data install to":            "GoAdmin框架数据表安装位置。",
        "web.the file path of sqlite3 database":                  "SQLite3 数据库文件路径。",
        "web.please use absolute path when you start as service": "作为服务启动时，请使用绝对路径。",
        "web.module name is the path of go module":               "模块名指在go module中的路径",
        "web.port number which application will listen on":       "应用监听的端口号。",
        "web.url prefix of the running application":              "运行访问的路由前缀。",

        "web.official website":                 "官方网站",
        "web.current version":                  "当前版本",
        "web.goadmin web installation program": "GoAdmin Web安装程序",
        "web.installation page":                "安装页面",
        "web.install now":                      "立即安装",

        "web.result":          "安装结果",
        "web.ok":              "好的",
        "web.wrong parameter": "错误的参数",
        "web.install success": "安装成功~~",
    },
    "en": {
        "cn": "Chinese",
        "en": "English",
        "jp": "Japanese",
        "tc": "Traditional Chinese",

        "web.database settings":     "Database Settings",
        "web.installation settings": "Installation Settings",
        "web.application settings":  "Application General Settings",
        "web.optional Settings":     "Optional Settings",

        "web.database type":     "Database Type",
        "web.database host":     "Host",
        "web.database user":     "User",
        "web.database password": "Password",
        "web.database name":     "Database Name",
        "web.database file":     "Path",
        "web.database port":     "port",
        "web.database schema":   "Schema",

        "web.theme":             "Theme",
        "web.language":          "Language",
        "web.web framework":     "Web Framework",
        "web.module name":       "Module Name",
        "web.http port":         "HTTP Port",
        "web.url prefix":        "Url Prefix",
        "web.website title":     "Website Title",
        "web.login page logo":   "Login Page Logo",
        "web.sidebar logo":      "SideBar Logo",
        "web.sidebar mini logo": "SideBar Mini Logo",
        "web.use orm":           "Use ORM",
        "web.no use":            "No use",
        "web.input":             "Input",

        "web.simplified chinese":  "Simplified Chinese",
        "web.traditional chinese": "Traditional Chinese",
        "web.english":             "English",
        "web.japanese":            "Japanese",

        "web.where the framework sql data install to":            "Where the framework sql data will be installed to。",
        "web.the file path of sqlite3 database. ":                "The file path of SQLite3 database. ",
        "web.please use absolute path when you start as service": "Please use absolute path when you start as service.",
        "web.module name is the path of go module":               "Module name is the path of go module.",
        "web.port number which application will listen on":       "Port number which application will listen on.",
        "web.url prefix of the running application":              "Url prefix of the running application.",

        "web.official website":                 "Official Website",
        "web.current version":                  "Current Version",
        "web.goadmin web installation program": "GoAdmin Web Installation Program",
        "web.installation page":                "Installation Page",
        "web.install now":                      "Install GoAdmin",

        "web.result":          "Installation Result",
        "web.ok":              "Ok",
        "web.wrong parameter": "Wrong parameter",
        "web.install success": "Install Success~~",
    },
}

var defaultLang = "en"

func setDefaultLangSet(set string) {
    if set != "" && (set == "cn" || set == "en") {
        defaultLang = set
    }
}

func local(lang string) func(string) string {
    if _, ok := langs[defaultLang]; ok {
        return func(msg string) string {
            return langs[lang][msg]
        }
    }
    return nil
}

func getWord(msg string) string {
    if word, ok := langs[defaultLang][msg]; ok {
        return word
    }
    return msg
}
