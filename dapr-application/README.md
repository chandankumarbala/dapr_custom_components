mkdir dapr

# Clone dapr
git clone https://github.com/dapr/dapr.git dapr/dapr

# Clone component-contrib
git clone https://github.com/dapr/components-contrib.git dapr/components-contrib

https://github.com/dapr/components-contrib/blob/main/docs/developing-component.md

go mod edit -replace github.com/dapr/components-contrib=../components-contrib
make modtidy-all
make DEBUG=1 build
/dapr/dist/windows_amd64/debug/daprd.exe /c/Users/USER/.dapr/bin 


dapr run --enable-api-logging  --app-id relay-calls  --dapr-http-port 3500 --app-port 3000 --resources-path ./components --config ./components/config.yaml --log-level debug  -- node index.js