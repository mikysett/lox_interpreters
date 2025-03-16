package main

type Config struct {
	ForbidUnusedVariable        bool
	ForbidUninitializedVariable bool
	AllowImplicitStringCast     bool
	AllowStaticMethods          bool
	AllowAnonymousFunctions     bool
	AllowGettersInClasses       bool
}

var GlobalConfig = ConfigWithExtras

var ConfigWithExtras = Config{
	ForbidUnusedVariable:        true,
	ForbidUninitializedVariable: true,
	AllowImplicitStringCast:     true,
	AllowStaticMethods:          true,
	AllowAnonymousFunctions:     true,
	AllowGettersInClasses:       true,
}

var BasicConfig = Config{}
