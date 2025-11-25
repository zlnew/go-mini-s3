public-list:
	hurl --variables-file ./.env ./collections/list.hurl

public-stream:
	hurl --variables-file ./.env ./collections/stream.hurl

public-download:
	hurl --variables-file ./.env ./collections/download.hurl

public-upload:
	hurl --variables-file ./.env ./collections/upload.hurl --file-root /home/zlnew/www/learn/go-nano-cloud/.hurl/files

public-delete:
	hurl --variables-file ./.env ./collections/delete.hurl

# =========================================================

private-list:
	hurl --variables-file ./.env ./collections/list.hurl --variable bucket=private

private-stream:
	hurl --variables-file ./.env ./collections/stream.hurl --variable bucket=private

private-download:
	hurl --variables-file ./.env ./collections/download.hurl --variable bucket=private

private-upload:
	hurl --variables-file ./.env ./collections/upload.hurl --variable bucket=private --file-root /home/zlnew/www/learn/go-nano-cloud/.hurl/files

private-delete:
	hurl --variables-file ./.env ./collections/delete.hurl --variable bucket=private

