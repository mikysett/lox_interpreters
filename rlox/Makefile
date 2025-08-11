MAKEFLAGS	+=	--quiet

WHT	= \033[0;37m
BLK	= \033[0;30m
RED	= \033[0;31m
YEL	= \033[0;33m
BLU	= \033[0;34m
GRN	= \033[0;32m

NAME	=	rlox

DIR_S	=	src

RM		=	rm -fdr

$(NAME):
			cargo build --release
			mkdir -p ./bin
			cp target/release/$(NAME) ./bin/$(NAME)

install:	build
# todo: add it to paths so it can be run globally

build:		$(NAME)
			printf "$(WHT)[$(GRN)$(NAME) PROGRAM COMPILED$(WHT)]\n"

debug:
	        export RUSTFLAGS="--cfg debug" && cargo run

all:		build install

clean:
			$(RM) bin target
			printf "$(WHT)[$(YEL)REMOVED $(NAME) BINARIES$(WHT)]\n"

.PHONY:		all build install clean
