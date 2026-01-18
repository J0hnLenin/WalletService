.PHONY: gen-api gen-moks
gen-api:
	powershell -ExecutionPolicy Bypass -File ./scripts/gen_api.ps1
gen-moks:
	powershell -ExecutionPolicy Bypass -File ./scripts/gen-moks.ps1

