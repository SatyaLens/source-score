install-tools:
	echo Installing tools from tools/tools.go && \
	cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go insatll %