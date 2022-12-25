#!/bin/sh

TARGET_FILE=$1

MIDDLEWARE_PATH="\t\tmiddleware(c)\n\t\tif c.IsAborted() {\n\t\t\treturn\n\t\t}"

sed -i "s/.*middleware(c).*/$MIDDLEWARE_PATH/" "$TARGET_FILE"

rm -f "$TARGET_FILE.bak"