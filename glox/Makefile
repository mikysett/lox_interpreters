MAKEFLAGS	+=	--quiet

WHT	= \033[0;37m
BLK	= \033[0;30m
RED	= \033[0;31m
YEL	= \033[0;33m
BLU	= \033[0;34m
GRN	= \033[0;32m

NAME	=	glox

DIR_S	=	src

RM		=	rm -fdr

$(NAME):
			cd $(DIR_S) && go build -o ../bin/$(NAME)

install:
			cd $(DIR_S) && go install
			printf "$(WHT)[$(GRN)$(NAME) PROGRAM INSTALLED$(WHT)]\n"

build:		$(NAME)
			printf "$(WHT)[$(GRN)$(NAME) PROGRAM COMPILED$(WHT)]\n"

all:		test build install

clean:
			$(RM) bin
			printf "$(WHT)[$(YEL)REMOVED $(NAME) BINARIES$(WHT)]\n"

test: build
	printf "$(WHT)[$(RED)TODO$(WHT)]\n"

.PHONY:		all build install clean test
