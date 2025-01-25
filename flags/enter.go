package flags

import (
	"flag"
	"os"

	"github.com/LiangNing7/BlogX/flags/flag_user"
)

type Options struct {
	File    string
	DB      bool
	Version bool
	Type    string
	Sub     string
	ES      bool
}

var FlagOptions = new(Options)

func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FlagOptions.Version, "v", false, "版本")
	flag.StringVar(&FlagOptions.Type, "t", "", "类型")
	flag.StringVar(&FlagOptions.Sub, "s", "", "子类")
	flag.BoolVar(&FlagOptions.ES, "es", false, "创建索引")
	flag.Parse()
}

func Run() {
	if FlagOptions.DB {
		// 执行数据库迁移
		FlagDB()
		os.Exit(0)
	}
	if FlagOptions.ES {
		EsIndex()
		os.Exit(0)
	}

	switch FlagOptions.Type {
	case "user":
		u := flag_user.FlagUser{}
		switch FlagOptions.Sub {
		case "create":
			u.Create()
			os.Exit(0)
		}
	}
}
