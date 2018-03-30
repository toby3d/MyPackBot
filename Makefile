# Build package by default with all general tools and things
all:
	make localization
	make build

# Build minimal package only
build:
	go build

# Build debug version with more logs
debug:
	go build -tags=debug

# Format the source code
fmt:
	go fmt

check:
	gometalinter --fast

# Build localization files with separated untranslated strings
translation:
	goi18n merge -format yaml \
	-sourceLanguage en \
	-outdir ./i18n/ \
	./i18n/src/*/*

# Build localization files and merge untranslated strings
localization:
	make translation
	goi18n -format yaml \
	-sourceLanguage en \
	-outdir ./i18n/ \
	./i18n/*.all.yaml ./i18n/*.untranslated.yaml