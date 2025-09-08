MAKEFLAGS	+=	--quiet

WHT	= \033[0;37m
BLK	= \033[0;30m
RED	= \033[0;31m
YEL	= \033[0;33m
BLU	= \033[0;34m
GRN	= \033[0;32m

NAME	=	rlox

DIR_S			=	src
DIR_RESOURCES	=	resources

TARGET		 = --release
DIR_TARGET	 = release

ifeq ($(OPT),true)
	FLAGS =	--cfg optimize
endif

RM		=	rm -fdr

$(NAME):
			export RUSTFLAGS=$(FLAGS) && \
			cargo build $(TARGET)
			mkdir -p ./bin
			cp target/$(DIR_TARGET)/$(NAME) ./bin/$(NAME)

install:	build
# todo: add it to paths so it can be run globally

build:		$(NAME)
			printf "$(WHT)[$(GRN)$(NAME) PROGRAM COMPILED$(WHT)]\n"

debug:		FLAGS:="--cfg debug $(FLAGS)"
debug:		TARGET=
debug:		DIR_TARGET="debug"
debug:		$(NAME)

speedtest:	FLAGS:="--cfg speedtest $(FLAGS)"
speedtest:  build
			printf "$(WHT)[$(BLU)RUNNING SPEED TEST$(WHT)]\n"
			hyperfine "./bin/$(NAME) $(DIR_RESOURCES)/speedtest.lox" --warmup=1

all:		build install

clean:
			$(RM) bin target
			printf "$(WHT)[$(YEL)REMOVED $(NAME) BINARIES$(WHT)]\n"

.PHONY:		all build install speedtest debug clean
