#!/usr/bin/env bash
gox --output "out/{{.Dir}}_{{.OS}}_{{.Arch}}" && ghr -u Webstrates $1 out
