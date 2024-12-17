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


- $ ./tester.sh
- Without binding: 1010.82mill / with binding: 243.06mill
Without binding: 1231.89mill / with binding: 243.39mill
Without binding: 1078.84mill / with binding: 243.49mill
Without binding: 804.31mill / with binding: 711.75mill
Without binding: 1317.37mill / with binding: 530.38mill
Without binding: 1707.83mill / with binding: 243.35mill
Without binding: 800.58mill / with binding: 544.58mill
Without binding: 1005.27mill / with binding: 244.59mill
Without binding: 721.46mill / with binding: 243.40mill
Without binding: 1048.11mill / with binding: 355.92mill
