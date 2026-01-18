.PHONY: gen-api gen-moks
gen-api:
	powershell -ExecutionPolicy Bypass -File ./scripts/gen-api.ps1
gen-moks:
	powershell -ExecutionPolicy Bypass -File ./scripts/gen-moks.ps1
run-tests:
	powershell -ExecutionPolicy Bypass -File ./scripts/run-tests.ps1

