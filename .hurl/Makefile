ifneq (,$(wildcard .env))
include .env
export $(shell sed 's/=.*//' .env)
endif

# Test mode toggle
ifeq ($(TEST),true)
HURL_MODE = --test
else
HURL_MODE =
endif

public-list:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/list.hurl
public-stream:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/stream.hurl --output $(FILE_ROOT)/${FILEPATH}
public-download:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/download.hurl --output $(FILE_ROOT)/${FILEPATH}
public-upload:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/upload.hurl --file-root $(FILE_ROOT)
public-delete:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/delete.hurl

# =========================================================

private-list:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/list.hurl --variable BUCKET=private
private-stream:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/stream.hurl --variable BUCKET=private --output $(FILE_ROOT)/${FILEPATH}
private-download:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/download.hurl --variable BUCKET=private --output $(FILE_ROOT)/${FILEPATH}
private-upload:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/upload.hurl --variable BUCKET=private --file-root $(FILE_ROOT)
private-delete:
	hurl $(HURL_MODE) --variables-file ./.env ./collections/delete.hurl --variable BUCKET=private

# Run groups
public-all: public-list public-stream public-download public-upload public-delete
private-all: private-list private-stream private-download private-upload private-delete
all: public-all private-all

