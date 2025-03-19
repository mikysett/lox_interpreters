package main

type Config struct {
	ForbidUnusedVariable        bool
	ForbidUninitializedVariable bool
	AllowImplicitStringCast     bool
	AllowStaticMethods          bool
	AllowAnonymousFunctions     bool
	AllowGettersInClasses       bool
	AllowContinueKeyword        bool
	AllowTernaryOperator        bool
	AllowModuloOperator         bool
}

var GlobalConfig = ConfigWithExtras

var ConfigWithExtras = Config{
	ForbidUnusedVariable:        true,
	ForbidUninitializedVariable: true,
	AllowImplicitStringCast:     true,
	AllowStaticMethods:          true,
	AllowAnonymousFunctions:     true,
	AllowGettersInClasses:       true,
	AllowContinueKeyword:        true,
	AllowTernaryOperator:        true,
	AllowModuloOperator:         true,
}

var BasicConfig = Config{}
