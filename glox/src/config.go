package main

type Config struct {
	ForbidUnusedVariable        bool
	ForbidUninitializedVariable bool
	AllowImplicitStringCast     bool
	AllowStaticMethods          bool
}

var GlobalConfig = ConfigWithExtras

var ConfigWithExtras = Config{
	ForbidUnusedVariable:        true,
	ForbidUninitializedVariable: true,
	AllowImplicitStringCast:     true,
	AllowStaticMethods:          true,
}

var BasicConfig = Config{}
